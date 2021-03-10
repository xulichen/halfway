package consts

const (
	ResponseCodeErrParameter        = 40023 // 参数错误
	ResponseCodeInternalServerError = 50000
)

const (
	ResponseRateLimitReachedStatusText = "您的手速太快了,服务器处理不过来。休息一下,稍后再来吧。" // 请求到达上限
	ResponseTextInternalServerError    = "服务器内部错误，请联系开发人员"
)
