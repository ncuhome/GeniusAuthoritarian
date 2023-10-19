package callback

const (
	cErrForm uint8 = iota + 1
	cErrSiteNotAllow
	cErrAppCodeNotFound
	cErrDBOperation
	cErrRemoteOperationFailed
	cErrUnauthorized
	cErrFindUnit
	cErrUnexpected
	cErrSignatureExpired
	cErrOperationIllegal
	cErrInsufficientPermissions
	cErrAppNameExist
	cErrGroupNotFound
	cErrAppNotFound
	cErrMfaCode
	cErrMfaAlreadyExist
	cErrMfaAddExpired
	cErrMfaNotExist
	cErrRequestFrequency
	cErrNetContextChanged
	cErrMfaRequired
	cErrSmsCoolDown
	cErrSendSmsFailed
	cErrIdentityCodeNotCorrect
	cErrSshNotFound
	cErrUserIdentity
	cErrLoginSessionExpired
)

var (
	ErrForm = &Msg{
		Code:       cErrForm,
		Msg:        "参数错误，请把输入截图发给前端同学",
		HttpStatus: 400,
	}
	ErrSiteNotAllow = &Msg{
		Code:       cErrSiteNotAllow,
		Msg:        "回调站点不在白名单内，请联系管理员添加",
		HttpStatus: 403,
	}
	ErrAppCodeNotFound = &Msg{
		Code:       cErrAppCodeNotFound,
		Msg:        "没有找到授权码",
		HttpStatus: 404,
	}
	ErrDBOperation = &Msg{
		Code:       cErrDBOperation,
		Msg:        "数据库操作失败，请反馈后端同学修复",
		HttpStatus: 500,
	}
	ErrRemoteOperationFailed = &Msg{
		Code:       cErrRemoteOperationFailed,
		Msg:        "远程调用异常，请稍候重试",
		HttpStatus: 500,
	}
	ErrUnauthorized = &Msg{
		Code:       cErrUnauthorized,
		Msg:        "身份校验失败，令牌失效或过期",
		HttpStatus: 401,
	}
	ErrFindUnit = &Msg{
		Code:       cErrFindUnit,
		Msg:        "找不到身份组，请联系管理员处理",
		HttpStatus: 403,
	}
	ErrUnexpected = &Msg{
		Code:       cErrUnexpected,
		Msg:        "发生预期外错误，请反馈后端同学",
		HttpStatus: 500,
	}
	ErrSignatureExpired = &Msg{
		Code:       cErrSignatureExpired,
		Msg:        "请求已过期，请重试",
		HttpStatus: 403,
	}
	ErrOperationIllegal = &Msg{
		Code:       cErrOperationIllegal,
		Msg:        "操作不合法",
		HttpStatus: 403,
	}
	ErrInsufficientPermissions = &Msg{
		Code:       cErrInsufficientPermissions,
		Msg:        "缺少对应身份组权限",
		HttpStatus: 403,
	}
	ErrAppNameExist = &Msg{
		Code:       cErrAppNameExist,
		Msg:        "名称已被占用",
		HttpStatus: 400,
	}
	ErrGroupNotFound = &Msg{
		Code:       cErrGroupNotFound,
		Msg:        "找不到匹配身份组",
		HttpStatus: 404,
	}
	ErrAppNotFound = &Msg{
		Code:       cErrAppNotFound,
		Msg:        "应用不存在",
		HttpStatus: 404,
	}
	ErrMfaCode = &Msg{
		Code:       cErrMfaCode,
		Msg:        "双因素校验码错误",
		HttpStatus: 403,
	}
	ErrMfaAlreadyExist = &Msg{
		Code:       cErrMfaAlreadyExist,
		Msg:        "双因素校验已开启",
		HttpStatus: 400,
	}
	ErrMfaAddExpired = &Msg{
		Code:       cErrMfaAddExpired,
		Msg:        "双因素认证绑定已过期，请刷新重试",
		HttpStatus: 404,
	}
	ErrMfaNotExist = &Msg{
		Code:       cErrMfaNotExist,
		Msg:        "双因素认证未开启",
		HttpStatus: 400,
	}
	ErrRequestFrequency = &Msg{
		Code:       cErrRequestFrequency,
		Msg:        "操作频繁，请稍后重试",
		HttpStatus: 403,
	}
	ErrNetContextChanged = &Msg{
		Code:       cErrNetContextChanged,
		Msg:        "网络环境异常变更，请重新登录",
		HttpStatus: 403,
	}
	ErrMfaRequired = &Msg{
		Code:       cErrMfaRequired,
		Msg:        "双因素校验码缺失",
		HttpStatus: 403,
	}
	ErrSmsCoolDown = &Msg{
		Code:       cErrSmsCoolDown,
		Msg:        "短信冷却中，请稍后重试",
		HttpStatus: 403,
	}
	ErrSendSmsFailed = &Msg{
		Code:       cErrSendSmsFailed,
		Msg:        "短信发送失败，请重试",
		HttpStatus: 500,
	}
	ErrIdentityCodeNotCorrect = &Msg{
		Code:       cErrIdentityCodeNotCorrect,
		Msg:        "身份校验码错误",
		HttpStatus: 403,
	}
	ErrSshNotFound = &Msg{
		Code:       cErrSshNotFound,
		Msg:        "SSH 账号未分配，请联系管理员",
		HttpStatus: 404,
	}
	ErrUserIdentity = &Msg{
		Code:       cErrUserIdentity,
		Msg:        "没有找到角色，请尝试使用其他登录方式或联系管理员",
		HttpStatus: 403,
	}
	ErrLoginSessionExpired = &Msg{
		Code:       cErrLoginSessionExpired,
		Msg:        "登录请求过期，请重新登录",
		HttpStatus: 403,
	}
)
