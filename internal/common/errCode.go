package common

/*
	统一错误返回码
*/

// 成功返回

const OK int32 = 0 //tzx修改，文档上成功是0

// 通用错误

const ServerCommonError int32 = 201
const RequestParamError int32 = 202
const TokenExpireError int32 = 203
const TokenGenerateError int32 = 204
const DbError int32 = 205
const DbUpdateAffectedZeroError int32 = 206
const RedisError int32 = 207
const DataUseUp int32 = 208
const TokenParseError int32 = 209
const UsernameRepetition = 220
const AuthorizationError = 221
