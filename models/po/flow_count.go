package po

import "time"

type FlowCount struct {
	ServiceName string        `description:"服务名"`
	Interval    time.Duration `description:"间隔期"`
	Unix        int64         `description:"时间戳"`
	Qps         int64         `description:"单位时间内计数"`
	TickerCount int64         `description:"间隔时间内计数"`
	TotalCount  int64         `description:"总计时间内计数"`
}
