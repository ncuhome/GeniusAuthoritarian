package redis

type Config struct {
	Addr     string
	Password string `config:"omitempty"`
}
