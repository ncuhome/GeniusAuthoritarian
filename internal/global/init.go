package global

func init() {
	initConfig()
	initLog()
	initJwtEd25519()
	checkConfig()
}
