package common

type HttpStatusCode int

// Http
const (
	SUCCESS_REGISTER = 1000 + iota
	SUCCESS_LOGIN
	SUCCESS_ADD_FRIEND
	SUCCESS_CREATE_GROUP
	SUCCESS_GET_GROUP_USERS
)
const (
	ERROR_INVALID_PARAMS = 2000 + iota
	ERROR_USER_NOT_EXIST
	ERROR_USER_EXIST
	ERROR_NEED_LOGIN
	ERROR_GENERATE_JWT
	ERROR_INVALID_TOKEN
	ERROR_ADD_FRIEND
	ERROR_CREATE_GROUP
	ERROR_GET_GROUP_USERS
	ERROR_USER_LOGINED
	ERROR_MYSQL
	ERROR_REDIS
)

var HttpMsg map[HttpStatusCode]string = map[HttpStatusCode]string{
	SUCCESS_REGISTER:        "注册成功",
	SUCCESS_LOGIN:           "登录成功",
	SUCCESS_ADD_FRIEND:      "添加好友成功",
	SUCCESS_CREATE_GROUP:    "创建群聊成功",
	SUCCESS_GET_GROUP_USERS: "获取群成员成功",
	ERROR_INVALID_PARAMS:    "无效参数",
	ERROR_USER_EXIST:        "用户存在",
	ERROR_USER_NOT_EXIST:    "用户不存在",
	ERROR_NEED_LOGIN:        "需要登录",
	ERROR_GENERATE_JWT:      "生成jwt失败",
	ERROR_INVALID_TOKEN:     "无效的token",
	ERROR_ADD_FRIEND:        "添加好友失败",
	ERROR_CREATE_GROUP:      "创建群聊失败",
	ERROR_GET_GROUP_USERS:   "获取群成员失败",
	ERROR_USER_LOGINED:      "用户已登录",
	ERROR_MYSQL:             "MySQL数据库内部错误",
	ERROR_REDIS:             "Redis数据库内部错误",
}

func (h HttpStatusCode) GetMsg() string {
	return HttpMsg[h]
}
