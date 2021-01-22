package dao

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_chat/common/conn"
	"github.com/mangenotwork/mange_chat/common/utils"
	"github.com/mangenotwork/mange_chat/obj"
)

// 聊天消息
type DaoMsg struct{}

// 保存消息到redis
func (dao *DaoMsg) Save(room, msg string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := fmt.Sprintf(obj.OnebyOneMsgKey, room)
	res, err := rc.Do("RPUSH", key, msg)
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	log.Println("保存成功: ", res)
}

// 从redis 获取该房间的消息,也就是历史消息
func (dao *DaoMsg) Get(room string) (datas []string) {
	datas = make([]string, 0)
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := fmt.Sprintf(obj.OnebyOneMsgKey, room)
	res, err := redis.Strings(rc.Do("LRANGE", key, 0, -1))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println("历史消息 = ", res)
	datas = res
	return
}

// 保存未读消息条数
func (dao *DaoMsg) SaveUnreadMsg(me, you string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := fmt.Sprintf(obj.UnreadMsgKey, me)

	// 返回哈希表 key 中给定域 field 的值。
	log.Println("保存未读消息条数 ")
	res, err := redis.String(rc.Do("HGET", key, you))
	if err != nil {
		log.Println("GET error", err.Error())
	}
	log.Println(res)

	if res == "" {
		//如果不存在
		log.Println("不存在")
		res1, err := rc.Do("HSETNX", key, you, "1")
		if err != nil {
			log.Println("GET error", err.Error())
			return
		}
		log.Println(res1)
	} else {
		//如果存在
		log.Println("存在 = ", res)
		count := utils.Str2Int(res) + 1
		countStr := utils.Int2Str(count)
		res, err := rc.Do("HSET", key, you, countStr)
		if err != nil {
			log.Println("GET error", err.Error())
		}
		log.Println(res)
	}

}

// 获取未读消息
func (dao *DaoMsg) GetUnreadMsg(me string) (datas map[string]string) {
	datas = make(map[string]string, 0)
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()

	key := fmt.Sprintf(obj.UnreadMsgKey, me)

	datas, err = redis.StringMap(rc.Do("HGETALL", key))
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	log.Println("获取未读消息 = ", datas)
	return
}

// 保存用户信息到redis
func (dao *DaoMsg) SetUser(me string, data map[string]string) {

	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()

	key := fmt.Sprintf(obj.UserInfoKey, me)

	args := redis.Args{}.Add(key)
	for k, v := range data {
		args = args.Add(k)
		args = args.Add(v)
	}

	log.Println("执行redis : ", "HMSET", args)
	res, err := rc.Do("HMSET", args...)
	if err != nil {
		log.Println("GET error", err.Error())
	}
	log.Println(res)
}

// 从redis获取用户信息
func (dao *DaoMsg) GetUserInfo(me string) (data map[string]string) {

	data = make(map[string]string, 0)

	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()

	key := fmt.Sprintf(obj.UserInfoKey, me)

	fmt.Println("执行redis : ", "HGETALL", key)
	data, err = redis.StringMap(rc.Do("HGETALL", key))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(data)

	return
}

// 从redis 获取名称是否重复
func (dao *DaoMsg) UserNameRepeat(name string) bool {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return true
	}
	defer rc.Close()
	key := fmt.Sprintf(obj.UserInfoKey, name)

	return EXISTSKey(rc, key)
}

//EXISTSKey 检查给定 key 是否存在。
// true 存在
// false 不存在
func EXISTSKey(rc redis.Conn, keyname string) bool {
	fmt.Println("[Execute redis command]: ", "EXISTS", keyname)
	datas, err := redis.String(rc.Do("DUMP", keyname))
	if err != nil || datas == "0" {
		fmt.Println("GET error", err.Error())
		return false
	}
	return true
}
