package log

import (
	"fmt"
	"github.com/hhzhhzhhz/logs-go"
	"os"
)

var logger logs_go.Logf

func Logger() logs_go.Logf {
	return logger
}

func Init() {
	var err error
	cfg := logs_go.NewLogfConfig()

	cfg.WriteDisk.GenerateRule = fmt.Sprintf("%s%s", os.Getenv("LOG_DIR"), "%Y-%d-%m/%Y%m%d%H.log")
	cfg.WriteDisk.MaxAge = 30
	cfg.WriteDisk.RotationTime = 30
	cfg.Stdout = true
	logger, err = cfg.BuildLogf()
	if err != nil {
		panic(err)
	}
}

func Close() {
	logger.Close()
}


