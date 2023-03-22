package callback

var (
	ErrForm = &Msg{
		Code:       1,
		Msg:        "参数错误，请把输入截图发给前端同学",
		HttpStatus: 400,
	}
)
