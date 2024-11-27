package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

}
