package keypair

import "fmt"

type Type string
type Format string

func (t Type) String() string {
	return string(t)
}
func (t Type) Format(f Format) FormattedType {
	return FormattedType{
		Type:   t,
		Format: f,
	}
}

const (
	TypePemPublic   Type = "PUBLIC KEY"
	TypePemPrivate  Type = "PRIVATE KEY"
	TypeCertificate Type = "CERTIFICATE"

	FormatRSA     Format = "RSA"
	FormatDSA     Format = "DSA"
	FormatECDSA   Format = "EC"
	FormatOpenSSH Format = "OPENSSH"
)

type FormattedType struct {
	Type
	Format
}

func (f FormattedType) String() string {
	return fmt.Sprintf("%s %s", f.Format, f.Type)
}

type PemType interface {
	String() string
}
