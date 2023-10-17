package logger

import (
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func deleteOldLogs(daysRange int) {

}

func getLogName() string {
	today := time.Now().Format("2006-01-02")
	return today
}

func NewLogger() (*zap.SugaredLogger, error) {
	//TODO: verify if already exists and create new
	//TODO: delete old logs
	go deleteOldLogs(7)

	logName := getLogName()
	path, _ := filepath.Abs("./logs/" + logName + ".log")

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{path, "stderr"}
	cfg.EncoderConfig.CallerKey = zapcore.OmitKey

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
