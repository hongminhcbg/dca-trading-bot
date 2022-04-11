package logs

import (
	"dca-bot/conf"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func Init(cfg *conf.Config) {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	if cfg.IsTestnet {
		// Only log the warning severity or above.
		log.SetLevel(log.InfoLevel)
		return
	}

	log.SetLevel(log.ErrorLevel)
}

func Info(msg string, kvs ...interface{}) {
	n := len(kvs)
	fields := make(log.Fields)
	for i := 0; i < n/2; i++ {
		fields[fmt.Sprintf("%v", kvs[2*i])] = kvs[2*i+1]
	}

	log.WithFields(fields).Info(msg)
}

func Error(err error, msg string, kvs ...interface{}) {
	n := len(kvs)
	fields := make(log.Fields)
	for i := 0; i < n/2; i++ {
		fields[fmt.Sprintf("%v", kvs[2*i])] = kvs[2*i+1]
	}
	fields["error"] = err

	log.WithFields(fields).Error(msg)
}
