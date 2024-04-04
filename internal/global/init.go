package global

func init() {
	initConfig()
	initLog()
	checkConfig()

	initJwtEd25519()
	initCaIssuer()
}
