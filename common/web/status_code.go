package web

type HttpStatusCode int

const (
	SUCCESS_REGISTER = 1000 + iota
	SUCCESS_LOGIN
)

const (
	ERROR_INVALID_PARAMS = 2000 + iota
)

var SuccessMsg map[HttpStatusCode]string = map[HttpStatusCode]string{
	SUCCESS_REGISTER: "注册成功",
	SUCCESS_LOGIN:    "登录成功",
}

var ErrorMsg map[HttpStatusCode]string = map[HttpStatusCode]string{
	ERROR_INVALID_PARAMS: "无效参数",
}

func (this HttpStatusCode) GetMsg() string {
	return ErrorMsg[this]
}
