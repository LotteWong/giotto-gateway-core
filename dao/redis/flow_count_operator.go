package redis

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/e421083458/golang_common/lib"
	"github.com/garyburd/redigo/redis"
)

type FlowCountOperator struct{}

func NewFlowCountOperator() *FlowCountOperator {
	return &FlowCountOperator{}
}

func (o *FlowCountOperator) GetFlowCount(req *po.FlowCount) *po.FlowCount {
	res := &po.FlowCount{
		ServiceName: req.ServiceName,
		Interval:    req.Interval,
		Unix:        0,
	}

	go func() {
		//defer func() {
		//	if err := recover(); err != nil {
		//		log.Println(err)
		//	}
		//}()
		ticker := time.NewTicker(req.Interval)
		for {
			<-ticker.C

			tickerCount := atomic.LoadInt64(&res.TickerCount)
			atomic.StoreInt64(&res.TickerCount, 0)
			// o.Increase(res)

			currTime := time.Now()
			hourKey := o.GetRedisHourKey(currTime, req)
			dayKey := o.GetRedisDayKey(currTime, req)

			// increase in redis
			err := RedisConfPipeline(func(conn redis.Conn) {
				conn.Send("INCRBY", hourKey, tickerCount)
				conn.Send("EXPIRE", hourKey, 24*60*60*2)
				conn.Send("INCRBY", dayKey, tickerCount)
				conn.Send("EXPIRE", dayKey, 24*60*60*2)
			})
			if err != nil {
				continue
			}

			// increase in memory
			totalCount, err := o.GetRedisDayVal(currTime, req)
			if err != nil {
				continue
			}
			nowUnix := time.Now().Unix()
			if res.Unix == 0 {
				res.Unix = time.Now().Unix()
				continue
			}
			tickerCount = totalCount - res.TotalCount
			if nowUnix > res.Unix {
				res.TotalCount = totalCount
				res.Qps = tickerCount / (nowUnix - res.Unix)
				// res.TickerCount = tickerCount
				res.Unix = time.Now().Unix()
			}
		}
	}()

	return res
}

func (o *FlowCountOperator) Increase(req *po.FlowCount) {
	go func() {
		//defer func() {
		//	if err := recover(); err != nil {
		//		log.Println(err)
		//	}
		//}()
		atomic.AddInt64(&req.TickerCount, 1)
	}()
}

func (o *FlowCountOperator) GetRedisHourKey(t time.Time, req *po.FlowCount) string {
	key := constants.FlowHourCountKey
	date := t.In(lib.TimeLocation).Format("2006010215")
	app := req.ServiceName
	return fmt.Sprintf("%s_%s_%s", key, date, app)
}

func (o *FlowCountOperator) GetRedisHourVal(t time.Time, req *po.FlowCount) (int64, error) {
	return redis.Int64(RedisConfDo("GET", o.GetRedisHourKey(t, req)))
}

func (o *FlowCountOperator) GetRedisDayKey(t time.Time, req *po.FlowCount) string {
	key := constants.FlowDayCountKey
	date := t.In(lib.TimeLocation).Format("20060102")
	app := req.ServiceName
	return fmt.Sprintf("%s_%s_%s", key, date, app)
}

func (o *FlowCountOperator) GetRedisDayVal(t time.Time, req *po.FlowCount) (int64, error) {
	return redis.Int64(RedisConfDo("GET", o.GetRedisDayKey(t, req)))
}
