package dto

type SshDeploy struct {
	UID       uint `gorm:"column:uid"`
	PublicSsh string
}

type SshSecrets struct {
	Username string     `json:"username"`
	Pem      SshKeyPair `json:"pem"`
	Ssh      SshKeyPair `json:"ssh"`
}

type SshKeyPair struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}
