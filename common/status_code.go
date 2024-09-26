package common

type HttpStatusCode int

// Http
const (
	SUCCESS_REGISTER = 1000 + iota
	SUCCESS_LOGIN
)
const (
	ERROR_INVALID_PARAMS = 2000 + iota
	ERROR_USER_NOT_EXIST
	ERROR_USER_EXIST
	ERROR_NEED_LOGIN
	ERROR_GENERATE_JWT
	ERROR_INVALID_TOKEN
	ERROR_MYSQL
)

var HttpMsg map[HttpStatusCode]string = map[HttpStatusCode]string{
	SUCCESS_REGISTER:     "注册成功",
	SUCCESS_LOGIN:        "登录成功",
	ERROR_INVALID_PARAMS: "无效参数",
	ERROR_USER_EXIST:     "用户存在",
	ERROR_USER_NOT_EXIST: "用户不存在",
	ERROR_NEED_LOGIN:     "需要登录",
	ERROR_GENERATE_JWT:   "生成jwt失败",
	ERROR_INVALID_TOKEN:  "无效的token",
	ERROR_MYSQL:          "数据库内部错误",
}

func (h HttpStatusCode) GetMsg() string {
	return HttpMsg[h]
}
