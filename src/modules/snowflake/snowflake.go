/**
 * 用这个生成id
 */
package snowflake

var Worker *IdWorker

func InitWorker(workerId int64) error {
	var err error
	Worker, err = NewIdWorker(workerId)
	if err != nil {
		return err
	}

	return nil
}
