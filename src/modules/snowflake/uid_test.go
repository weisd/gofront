package snowflake

import (
	"fmt"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	fmt.Println("start generate")
	iw, _ := NewIdWorker(2)
	for i := 0; i < 100; i++ {
		id, err := iw.NextId()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(id)
		}
	}
	fmt.Println("end generate")
}
