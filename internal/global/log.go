package global

import log "github.com/sirupsen/logrus"

func initLog() {
	log.SetLevel(log.DebugLevel)
}
