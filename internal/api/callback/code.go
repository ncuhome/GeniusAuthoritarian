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
	cErrAlreadyExist
)

var (
	ErrForm = &Msg{
		Code:       cErrForm,
		Msg:        "参数错误，请把输入截图发给前端同学",
		HttpStatus: 400,
	}
	ErrSiteNotAllow = &Msg{
		Code:       cErrSiteNotAllow,
		Msg:        "目标站点不在授权范围，请联系管理员添加",
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
		Msg:        "身份校验失败，权限不足",
		HttpStatus: 401,
	}
	ErrFindUnit = &Msg{
		Code:       cErrFindUnit,
		Msg:        "找不到身份组，请联系管理员处理",
		HttpStatus: 401,
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
	ErrAlreadyExist = &Msg{
		Code:       cErrAlreadyExist,
		Msg:        "目标已存在",
		HttpStatus: 400,
	}
)
