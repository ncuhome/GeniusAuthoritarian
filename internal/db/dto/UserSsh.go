package dto

type SshDeploy struct {
	UID       uint `gorm:"column:uid"`
	PublicSsh string
}
