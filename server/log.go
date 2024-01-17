package server

import log "github.com/sirupsen/logrus"

func init() {
	formatter := log.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		TimestampFormat:           "2006-01-02 15:04:05",
		FullTimestamp:             true,
	}
	log.SetFormatter(&formatter)
}
