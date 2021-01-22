package dao

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_chat/common/conn"
	"github.com/mangenotwork/mange_chat/common/constname"
	"github.com/mangenotwork/mange_chat/common/utils"
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
	key := fmt.Sprintf(constname.OnebyOneMsgKey, room)
	_, err = rc.Do("RPUSH", key, msg)
	if err != nil {
		//log.Println("GET error", err.Error())
		return
	}
	//log.Println("保存成功: ", res)
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
	key := fmt.Sprintf(constname.OnebyOneMsgKey, room)
	res, err := redis.Strings(rc.Do("LRANGE", key, 0, -1))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	//fmt.Println("历史消息 = ", res)
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
	key := fmt.Sprintf(constname.UnreadMsgKey, me)

	// 返回哈希表 key 中给定域 field 的值。
	log.Println("保存未读消息条数 ")
	res, err := redis.String(rc.Do("HGET", key, you))
	if err != nil {
		log.Println("GET error", err.Error())
	}
	log.Println(res)

	if res == "" {
		//如果不存在
		//log.Println("不存在")
		_, err := rc.Do("HSETNX", key, you, "1")
		if err != nil {
			log.Println("GET error", err.Error())
			return
		}
		//log.Println(res1)
	} else {
		//如果存在
		//log.Println("存在 = ", res)
		count := utils.Str2Int(res) + 1
		countStr := utils.Int2Str(count)
		_, err := rc.Do("HSET", key, you, countStr)
		if err != nil {
			log.Println("GET error", err.Error())
		}
		//log.Println(res)
	}
}

// 清零未读消息条数
func (dao *DaoMsg) EmptyUnreadMsg(me, you string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := fmt.Sprintf(constname.UnreadMsgKey, me)

	log.Println("执行redis : ", "HDEL", key)
	res, err := rc.Do("HDEL", key, you)
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	log.Println(res)
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

	key := fmt.Sprintf(constname.UnreadMsgKey, me)

	datas, err = redis.StringMap(rc.Do("HGETALL", key))
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	//log.Println("获取未读消息 = ", datas)
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

	key := fmt.Sprintf(constname.UserInfoKey, me)

	args := redis.Args{}.Add(key)
	for k, v := range data {
		args = args.Add(k)
		args = args.Add(v)
	}

	//log.Println("执行redis : ", "HMSET", args)
	_, err = rc.Do("HMSET", args...)
	if err != nil {
		log.Println("GET error", err.Error())
	}
	//log.Println(res)
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

	key := fmt.Sprintf(constname.UserInfoKey, me)

	//log.Println("执行redis : ", "HGETALL", key)
	data, err = redis.StringMap(rc.Do("HGETALL", key))
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	//log.Println(data)

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
	key := fmt.Sprintf(constname.UserInfoKey, name)

	return EXISTSKey(rc, key)
}

//EXISTSKey 检查给定 key 是否存在。
// true 存在
// false 不存在
func EXISTSKey(rc redis.Conn, keyname string) bool {
	//log.Println("[Execute redis command]: ", "EXISTS", keyname)
	datas, err := redis.String(rc.Do("DUMP", keyname))
	if err != nil || datas == "0" {
		fmt.Println("GET error", err.Error())
		return false
	}
	return true
}

// 进入用户列表
// 列表是一个zset, 后面的score(大到小排序) 是登录时间戳
// 如果用户下线 后面的score 为0(也就是为0就是没在线)
func (dao *DaoMsg) UserToOnline(name string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := constname.AllUserKey

	args := redis.Args{}.Add(key)
	args = args.Add(time.Now().Unix())
	args = args.Add(name)

	//log.Println("执行redis : ", "ZADD", args)
	_, err = rc.Do("ZADD", args...)
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	//log.Println(res)
}

// 用户下线
func (dao *DaoMsg) UserOutOnline(name string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := constname.AllUserKey

	args := redis.Args{}.Add(key)
	args = args.Add(0)
	args = args.Add(name)

	//log.Println("执行redis : ", "ZADD", args)
	_, err = rc.Do("ZADD", args...)
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	//log.Println(res)
}

// 在redis获取用户列表
func (dao *DaoMsg) GetAllUser() (res map[string]string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := constname.AllUserKey

	start := 0
	stop := -1

	//log.Println("执行redis : ", "ZREVRANGE", key, start, stop, "WITHSCORES")
	res, err = redis.StringMap(rc.Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	//log.Println(res)
	return
}

//   ---------------------------  群聊  ---------------------------

// 新建群聊
// 用zet 保存群聊
func (dao *DaoMsg) NewRoom(name string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := constname.RoomList

	args := redis.Args{}.Add(key)
	args = args.Add(time.Now().Unix())
	args = args.Add(name)

	//log.Println("执行redis : ", "ZADD", args)
	_, err = rc.Do("ZADD", args...)
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
}

// 获取所有群聊房间
func (dao *DaoMsg) GetRoomList() (res map[string]string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := constname.RoomList

	start := 0
	stop := -1

	//log.Println("执行redis : ", "ZREVRANGE", key, start, stop, "WITHSCORES")
	res, err = redis.StringMap(rc.Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		log.Println("GET error", err.Error())
		return
	}
	//log.Println(res)
	return
}

// 保存群聊消息
func (dao *DaoMsg) SaveRoomMsg(room, msg string) {
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := fmt.Sprintf(constname.RoomMsgKey, room)
	_, err = rc.Do("RPUSH", key, msg)
	if err != nil {
		//log.Println("GET error", err.Error())
		return
	}
}

// 获取群聊消息
func (dao *DaoMsg) GetRoomMsg(room string) (datas []string) {
	datas = make([]string, 0)
	rc, err := conn.RedisConn()
	if err != nil {
		log.Println("未获取到redis连接")
		return
	}
	defer rc.Close()
	key := fmt.Sprintf(constname.RoomMsgKey, room)
	res, err := redis.Strings(rc.Do("LRANGE", key, 0, -1))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	//fmt.Println("历史消息 = ", res)
	datas = res
	return
}
