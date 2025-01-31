package common

import (
	"log"

	"go.uber.org/zap"
)

// Log global logger
var Log *zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	Log = logger.Sugar()
	defer Log.Sync()
}
