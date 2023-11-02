package feishu

import "encoding/json"

type EventHeader struct {
	EventID    string `json:"event_id"`
	Token      string `json:"token"`
	CreateTime string `json:"create_time"`
	EventType  string `json:"event_type"`
	TenantKey  string `json:"tenant_key"`
	AppID      string `json:"app_id"`
}

type Event struct {
	Schema string          `json:"schema"`
	Header EventHeader     `json:"header"`
	Event  json.RawMessage `json:"event"`

	// 以下配置订阅时使用

	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

type UserDeletedEvent struct {
	Object struct {
		OpenId          string `json:"open_id"`
		UnionId         string `json:"union_id"`
		UserId          string `json:"user_id"`
		Name            string `json:"name"`
		EnName          string `json:"en_name"`
		Nickname        string `json:"nickname"`
		Email           string `json:"email"`
		EnterpriseEmail string `json:"enterprise_email"`
		JobTitle        string `json:"job_title"`
		Mobile          string `json:"mobile"`
		Gender          int    `json:"gender"`
		Avatar          struct {
			Avatar72     string `json:"avatar_72"`
			Avatar240    string `json:"avatar_240"`
			Avatar640    string `json:"avatar_640"`
			AvatarOrigin string `json:"avatar_origin"`
		} `json:"avatar"`
		Status struct {
			IsFrozen    bool `json:"is_frozen"`
			IsResigned  bool `json:"is_resigned"`
			IsActivated bool `json:"is_activated"`
			IsExited    bool `json:"is_exited"`
			IsUnjoin    bool `json:"is_unjoin"`
		} `json:"status"`
		DepartmentIds []string `json:"department_ids"`
		LeaderUserId  string   `json:"leader_user_id"`
		City          string   `json:"city"`
		Country       string   `json:"country"`
		WorkStation   string   `json:"work_station"`
		JoinTime      int      `json:"join_time"`
		EmployeeNo    string   `json:"employee_no"`
		EmployeeType  int      `json:"employee_type"`
		Orders        []struct {
			DepartmentId    string `json:"department_id"`
			UserOrder       int    `json:"user_order"`
			DepartmentOrder int    `json:"department_order"`
			IsPrimaryDept   bool   `json:"is_primary_dept"`
		} `json:"orders"`
		CustomAttrs []struct {
			Type  string `json:"type"`
			Id    string `json:"id"`
			Value struct {
				Text        string `json:"text"`
				Url         string `json:"url"`
				PcUrl       string `json:"pc_url"`
				OptionId    string `json:"option_id"`
				OptionValue string `json:"option_value"`
				Name        string `json:"name"`
				PictureUrl  string `json:"picture_url"`
				GenericUser struct {
					Id   string `json:"id"`
					Type int    `json:"type"`
				} `json:"generic_user"`
			} `json:"value"`
		} `json:"custom_attrs"`
		JobLevelId              string   `json:"job_level_id"`
		JobFamilyId             string   `json:"job_family_id"`
		DottedLineLeaderUserIds []string `json:"dotted_line_leader_user_ids"`
	} `json:"object"`
	OldObject struct {
		DepartmentIds []string `json:"department_ids"`
		OpenId        string   `json:"open_id"`
	} `json:"old_object"`
}
