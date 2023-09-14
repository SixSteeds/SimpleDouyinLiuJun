package common

var message map[int32]string

func init() {
	message = make(map[int32]string)
	message[OK] = "SUCCESS"
	message[ServerCommonError] = "服务器开小差啦,稍后再来试一试"
	message[RequestParamError] = "参数错误"
	message[TokenExpireError] = "token失效，请重新登陆"
	message[TokenGenerateError] = "生成token失败"
	message[DbError] = "数据库繁忙,请稍后再试"
	message[DbUpdateAffectedZeroError] = "更新数据影响行数为0"
	message[RedisError] = "redis执行失败，请重新尝试"
	message[DataUseUp] = "数据库资源已经全部展示给你啦"
	message[TokenParseError] = "TOKEN解析失败"

	message[UsernameRepetition] = "用户名重复"
	message[AuthorizationError] = "用户名或密码错误"

}

func MapErrMsg(errcode int32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

//func IsCodeErr(errcode int32) bool {
//	if _, ok := message[errcode]; ok {
//		return true
//	} else {
//		return false
//	}
//}
