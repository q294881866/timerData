package task

import (
	"fmt"
	"time"
)

type Task interface {
	Exec()
}

type Producer struct {
}

func (p *Producer) Exec() {
	fmt.Println("The machines will running!")
	for i := 1000; i < 6000; i++ {
		datas := make([]string, 60)
		filer := time.NewTicker(time.Minute * 1)
		MiTask(filer, datas, i)
	}
}

// 工作am8:00-pm9:00
// 改为任何时间都工作 代码注释掉了
func WorkTime() {
//	now := time.Now()
//	hour := now.Hour()
//	//第二天早上八点的距离
//	sleep := time.Date(2017, 1, 2, 8, 0, 0, 0, time.Local).UnixNano() - time.Date(2017, 1, 1, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.Local).UnixNano()
//	if hour > 21 {
//		fmt.Println(sleep)
//		time.Sleep(time.Duration(sleep))
//	}
}
