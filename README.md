# mange_chat
 web socket  IM 练手项目

# 技术选型
- socket : github.com/gorilla/websocket
- web : github.com/gin-gonic/gin
- UI : https://github.com/lihongxun945/jquery-weui
- 存储 redis : github.com/garyburd/redigo/redis
- 存储 mongo :

# 文档:
- socket : https://pkg.go.dev/github.com/gorilla/websocket
- gin : https://www.kancloud.cn/shuangdeyu/gin_book/949412
- UI : https://jqweui.cn/components


# 记录:
1. v0.1 20210118 : 匿名聊天，聊天室群聊
2. v0.2 20210119 : 一对一聊天室，首页大厅
3. v0.3 20210120 : redis存储，消息记录，未读消息数量
4. v0.4 20210121 : 首页UI设计
5. v0.5 20210121 : 聊天页面UI设计
6. v0.6 20210122 : 登录，退出, 头像选择，用户存储到redis
7. v0.7 20210122 : 用户列表(redis),未读消息提示
8. v0.8 20210122 : 群聊房间列表(redis)， 群聊消息保存到redis 

# TODO:
v0.9 发送图片
