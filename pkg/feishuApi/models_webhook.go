package feishuApi

import (
	"encoding/json"
)

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
	Object    User `json:"object"`
	OldObject struct {
		DepartmentIds []string `json:"department_ids"`
		OpenId        string   `json:"open_id"`
	} `json:"old_object"`
}

type UserUpdatedEvent struct {
	Object User `json:"object"`
	// 只有变更字段有值
	OldObject User `json:"old_object"`
}
