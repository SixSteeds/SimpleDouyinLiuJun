package common

/*
	此处考虑后期统一错误返回码
*/
// 成功返回
const OK int32 = 200

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码
const SERVER_COMMON_ERROR int32 = 201
const REUQEST_PARAM_ERROR int32 = 202
const TOKEN_EXPIRE_ERROR int32 = 203
const TOKEN_GENERATE_ERROR int32 = 204
const DB_ERROR int32 = 205
const DB_UPDATE_AFFECTED_ZERO_ERROR int32 = 206

//用户模块
