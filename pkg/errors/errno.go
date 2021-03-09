package errors

const (
	ErrCustom int32 = 999999
)

// TODO：完善错误处理
var (
	ErrSign                = newError(100001, "签名错误")
	ErrInvalidParam        = newError(100002, "参数错误")
	ErrInvalidToken        = newError(100003, "无效的令牌")
	ErrTokenExpired        = newError(100004, "令牌过期")
	ErrTokenRevoked        = newError(100005, "令牌已失效")
	ErrUnAuthorized        = newError(100006, "未授权")
	ErrNoLogin             = newError(100007, "未登录")
	ErrPassword            = newError(100008, "用户名或密码错误")
	ErrUserNotExists       = newError(100009, "用户不存在")
	ErrUserAlreadyExists   = newError(100010, "用户已存在")
	ErrAccountNotAvailable = newError(100011, "账户不可用")
	ErrSendCodeTooFrequent = newError(100012, "验证码发送过于频繁")
	ErrSendCodeLimit       = newError(100012, "当日发送验证码次数过多")
	ErrRequireSendCode     = newError(100013, "请先发送短信验证码")
	ErrVerifyCodeExpired   = newError(100014, "验证码过期")
	ErrVerifyCodeLimit     = newError(100015, "验证码失效")
	ErrVerifyCode          = newError(100016, "验证码错误")
	ErrInvalidFromUserTag  = newError(100017, "技能已被删除") // 已方技能标签已删除
	ErrInvalidToUserTag    = newError(100018, "该用户已失效") // 对方技能标签已删除
)
