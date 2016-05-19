// This package provides unique id in distribute system
// the algorithm is inspired by Twitter's famous snowflake
// its link is: https://github.com/twitter/snowflake/releases/tag/snowflake-2010
//

// 0               41	           51			 64
// +---------------+----------------+------------+
// |timestamp(ms)  | worker node id | sequence	 |
// +---------------+----------------+------------+

// Copyright (C) 2016 by zheng-ji.info

package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	CEpoch         = 1457159242
	CWorkerIdBits  = 10 // Num of WorkerId Bits
	CSenquenceBits = 12 // Num of Sequence Bits

	CWorkerIdShift  = 12
	CTimeStampShift = 22
)

// IdWorker Struct
type IdWorker struct {
	workerId      int64
	lastTimeStamp int64
	sequence      int64
	sequenceMask  int64
	maxWorkerId   int64
	lock          *sync.Mutex
}

// NewIdWorker Func: Generate NewIdWorker with Given workerid
func NewIdWorker(workerid int64) (iw *IdWorker, err error) {
	iw = new(IdWorker)

	iw.maxWorkerId = getMaxWorkerId()

	if workerid > iw.maxWorkerId || workerid < 0 {
		return nil, errors.New("worker not fit")
	}
	iw.workerId = workerid
	iw.lastTimeStamp = -1
	iw.sequence = 0
	iw.sequenceMask = getSequenceMask()
	iw.lock = new(sync.Mutex)
	return iw, nil
}

func getMaxWorkerId() int64 {
	return -1 ^ -1<<CWorkerIdBits
}

func getSequenceMask() int64 {
	return -1 ^ -1<<CSenquenceBits
}

func (iw *IdWorker) timeGen() int64 {
	return time.Now().UnixNano()
}

func (iw *IdWorker) timeReGen(last int64) int64 {
	ts := time.Now().UnixNano()
	for {
		if ts < last {
			ts = iw.timeGen()
		} else {
			break
		}
	}
	return ts
}

// NewId Func: Generate next id
func (iw *IdWorker) NextId() (ts int64, err error) {
	iw.lock.Lock()
	defer iw.lock.Unlock()
	ts = iw.timeGen()
	if ts == iw.lastTimeStamp {
		iw.sequence = (iw.sequence + 1) & iw.sequenceMask
		if iw.sequence == 0 {
			ts = iw.timeReGen(ts)
		}
	} else {
		iw.sequence = 0
	}

	if ts < iw.lastTimeStamp {
		err = errors.New("Clock moved backwards, Refuse gen id")
		return 0, err
	}
	iw.lastTimeStamp = ts
	ts = ts - CEpoch<<CTimeStampShift | iw.workerId<<CWorkerIdShift | iw.sequence
	return ts, nil
}
