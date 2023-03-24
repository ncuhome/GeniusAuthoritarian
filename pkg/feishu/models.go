package feishu

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type TenantAccessTokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int64  `json:"expire"`
}

type UserInfoResp struct {
	User struct {
		Avatar struct {
			Avatar240    string `json:"avatar_240"`
			Avatar640    string `json:"avatar_640"`
			Avatar72     string `json:"avatar_72"`
			AvatarOrigin string `json:"avatar_origin"`
		} `json:"avatar"`
		City            string        `json:"city"`
		Country         string        `json:"country"`
		CustomAttrs     []interface{} `json:"custom_attrs"`
		DepartmentIds   []string      `json:"department_ids"`
		Description     string        `json:"description"`
		EmployeeNo      string        `json:"employee_no"`
		EmployeeType    int           `json:"employee_type"`
		EnName          string        `json:"en_name"`
		EnterpriseEmail string        `json:"enterprise_email"`
		Gender          int           `json:"gender"`
		IsTenantManager bool          `json:"is_tenant_manager"`
		JobTitle        string        `json:"job_title"`
		JoinTime        int           `json:"join_time"`
		MobileVisible   bool          `json:"mobile_visible"`
		Name            string        `json:"name"`
		OpenId          string        `json:"open_id"`
		Orders          []struct {
			DepartmentId    string `json:"department_id"`
			DepartmentOrder int    `json:"department_order"`
			UserOrder       int    `json:"user_order"`
		} `json:"orders"`
		Status struct {
			IsActivated bool `json:"is_activated"`
			IsExited    bool `json:"is_exited"`
			IsFrozen    bool `json:"is_frozen"`
			IsResigned  bool `json:"is_resigned"`
			IsUnjoin    bool `json:"is_unjoin"`
		} `json:"status"`
		UnionId     string `json:"union_id"`
		WorkStation string `json:"work_station"`
	} `json:"user"`
}
