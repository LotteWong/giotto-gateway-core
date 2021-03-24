package redis

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/garyburd/redigo/redis"
)

// RedisConfPipeline Execute batch commands
func RedisConfPipeline(pipe ...func(conn redis.Conn)) error {
	conn, err := lib.RedisConnFactory("default")
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, item := range pipe {
		item(conn)
	}
	conn.Flush()
	return nil
}

// RedisConfDo Execute single command
func RedisConfDo(cmd string, args ...interface{}) (interface{}, error) {
	conn, err := lib.RedisConnFactory("default")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.Do(cmd, args...)
}
