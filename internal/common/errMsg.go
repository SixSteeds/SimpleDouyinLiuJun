package common

var message map[int32]string

func init() {
	message = make(map[int32]string)
	message[OK] = "SUCCESS"
	message[SERVER_COMMON_ERROR] = "服务器开小差啦,稍后再来试一试"
	message[REUQEST_PARAM_ERROR] = "参数错误"
	message[TOKEN_EXPIRE_ERROR] = "token失效，请重新登陆"
	message[TOKEN_GENERATE_ERROR] = "生成token失败"
	message[DB_ERROR] = "数据库繁忙,请稍后再试"
	message[DB_UPDATE_AFFECTED_ZERO_ERROR] = "更新数据影响行数为0"
	message[REDIS_ERROR] = "redis执行失败，请重新尝试"
	message[DATA_USE_UP] = "数据库资源已经全部展示给你啦"

}

func MapErrMsg(errcode int32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errcode int32) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
