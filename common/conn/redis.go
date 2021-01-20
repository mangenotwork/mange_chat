package conn

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

var (
	RedisIP       string = "192.168.0.197"
	RedisPort     int    = 6379
	RedisPassword string = "love0021$%"
	DBID          int    = 4
)

// RConn 普通连接
// ip : redis服务地址
// port:  Redis 服务端口
// password  Redis 服务密码
// 返回redis连接
func RConn(ip string, port int, password string) (redis.Conn, error) {
	host := fmt.Sprintf("%s:%d", ip, port)
	conn, err := redis.Dial("tcp", host)
	if nil != err {
		log.Println("dial to redis addr err: ", err)
		return nil, err
	}

	log.Println(conn)

	//针对无密码的连接
	if password == "" {
		return conn, nil
	}

	if _, authErr := conn.Do("AUTH", password); authErr != nil {
		log.Println("redis auth password error: ", authErr)
		return nil, fmt.Errorf("redis auth password error: %s", authErr)
	}

	return conn, nil
}

//指定db的连接
func RedisConn() (redis.Conn, error) {
	rc, err := RConn(RedisIP, RedisPort, RedisPassword)
	if err != nil {
		log.Println("redis conn error: ", err)
	}
	_, err = rc.Do("select", fmt.Sprintf("%d", DBID))
	if err != nil {
		log.Println("redis select db error: ", err)
	}
	return rc, err
}
