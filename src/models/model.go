package models

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//  _ "github.com/jinzhu/gorm/dialects/postgres"
	//  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Drivers map[string]*gorm.DB

type GormService struct {
	Enable bool
	Debug  bool
	Driver string
	Host   string
	DB     string
	User   string
	Passwd string

	Path string // for sqlite,tidb

	MaxIdle int // 连接池的空闲数大小
	MaxOpen int // 最大打开连接数
	LogPath string
}

func Model(name ...string) *gorm.DB {
	k := "default"
	if len(name) > 0 {
		k = name[0]
	}
	if db, ok := Drivers[k]; ok {
		return db
	}

	panic(fmt.Errorf("model 不存在 %s", k))

	return nil
}

func HasModel(name string) bool {
	_, ok := Drivers[name]
	return ok
}

func InitModels(confs map[string]GormService) error {
	Drivers = make(map[string]*gorm.DB)
	for k, v := range confs {
		if !v.Enable {
			continue
		}
		db, err := newGorm(v)
		if err != nil {
			return err
		}

		Drivers[k] = db
	}

	return nil
}

func newGorm(conf GormService) (*gorm.DB, error) {
	dsn := ""
	switch conf.Driver {
	case "mysql":
		//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.User, conf.Passwd, conf.Host, conf.DB)
	// case "postgres":
	// case "sqlite3":
	default:
		return nil, fmt.Errorf("未知的 gorm 驱动：%s", conf.Driver)

	}

	db, err := gorm.Open(conf.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("newGorm 连接数据库失败 %v", err)
	}

	// if set this to true, `User`'s default table name will be `user`, table name setted with `TableName` won't be affected
	db.SingularTable(true)

	if conf.MaxIdle > 0 {
		db.DB().SetMaxIdleConns(conf.MaxIdle)
	}

	if conf.MaxOpen > 0 {
		db.DB().SetMaxOpenConns(conf.MaxOpen)
	}

	if conf.Debug {
		db.LogMode(true)
	}

	logpath := "./log/gorm.log"
	if len(conf.LogPath) > 0 {
		logpath, _ = filepath.Abs(conf.LogPath)
	}

	os.MkdirAll(path.Dir(logpath), os.ModePerm)
	// 日志
	f, err := os.Create(logpath)
	if err != nil {

		return nil, err
	}

	db.SetLogger(log.New(f, "\n", 0))

	return db, nil

}
