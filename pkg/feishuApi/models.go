package feishuApi

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

type User struct {
	UserId string `json:"user_id"`
	Avatar struct {
		Avatar240    string `json:"avatar_240"`
		Avatar640    string `json:"avatar_640"`
		Avatar72     string `json:"avatar_72"`
		AvatarOrigin string `json:"avatar_origin"`
	} `json:"avatar"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CustomAttrs []struct {
		Type  string `json:"type"`
		Id    string `json:"id"`
		Value struct {
			Text        string `json:"text"`
			Url         string `json:"url"`
			PcUrl       string `json:"pc_url"`
			OptionValue string `json:"option_value"`
			Name        string `json:"name"`
			PictureUrl  string `json:"picture_url"`
			GenericUser struct {
				Id   string `json:"id"`
				Type int    `json:"type"`
			} `json:"generic_user"`
		} `json:"value"`
	} `json:"custom_attrs"`
	DepartmentIds   []string `json:"department_ids"`
	Description     string   `json:"description"`
	EmployeeNo      string   `json:"employee_no"`
	EmployeeType    int      `json:"employee_type"`
	EnName          string   `json:"en_name"`
	EnterpriseEmail string   `json:"enterprise_email"`
	Gender          int      `json:"gender"`
	IsTenantManager bool     `json:"is_tenant_manager"`
	JobTitle        string   `json:"job_title"`
	JoinTime        int      `json:"join_time"`
	Email           string   `json:"email"`
	Mobile          string   `json:"mobile"`
	MobileVisible   bool     `json:"mobile_visible"`
	Name            string   `json:"name"`
	OpenId          string   `json:"open_id"`
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
}

type UserInfoResp struct {
	User User `json:"user"`
}

type Department struct {
	Name     string `json:"name"`
	I18NName struct {
		ZhCn string `json:"zh_cn"`
		JaJp string `json:"ja_jp"`
		EnUs string `json:"en_us"`
	} `json:"i18n_name"`
	ParentDepartmentId string   `json:"parent_department_id"`
	DepartmentId       string   `json:"department_id"`
	OpenDepartmentId   string   `json:"open_department_id"`
	LeaderUserId       string   `json:"leader_user_id"`
	ChatId             string   `json:"chat_id"`
	Order              string   `json:"order"`
	UnitIds            []string `json:"unit_ids"`
	MemberCount        int      `json:"member_count"`
	Status             struct {
		IsDeleted bool `json:"is_deleted"`
	} `json:"status"`
	CreateGroupChat bool `json:"create_group_chat"`
	Leaders         []struct {
		LeaderType int    `json:"leaderType"`
		LeaderID   string `json:"leaderID"`
	} `json:"leaders"`
	DepartmentHrbps []string `json:"department_hrbps"`
}

type ListDepartmentResp struct {
	HasMore   bool         `json:"has_more"`
	PageToken string       `json:"page_token"`
	Items     []Department `json:"items"`
}

type ListUserResp struct {
	HasMore   bool   `json:"has_more"`
	PageToken string `json:"page_token"`
	Items     []User `json:"items"`
}
