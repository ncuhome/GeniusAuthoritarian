package callback

const (
	cErrForm uint8 = iota + 1
	cErrSiteNotAllow
	cErrDBOperation
	cErrRemoteOperationFailed
	cErrUnauthorized
)

var (
	ErrForm = &Msg{
		Code:       cErrForm,
		Msg:        "参数错误，请把输入截图发给前端同学",
		HttpStatus: 400,
	}
	ErrSiteNotAllow = &Msg{
		Code:       cErrSiteNotAllow,
		Msg:        "该站点不在授权范围，请联系管理员添加",
		HttpStatus: 403,
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
)
