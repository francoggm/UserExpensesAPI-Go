package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/francoggm/go_expenses_api/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func deleteOldLogs(daysRange int) {
	startDay := time.Now().AddDate(0, 0, -daysRange)

	files, _ := os.ReadDir("./logs")

	for _, e := range files {
		fileName := e.Name()

		if strings.Contains(fileName, ".log") {
			fileName = strings.Split(fileName, ".")[0]

			if strings.Contains(fileName, "_") {
				fileName = strings.Split(fileName, "_")[0]
			}

			fileDate, err := time.Parse("2006-01-02", fileName)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if startDay.After(fileDate) {
				os.Remove("./logs/" + e.Name())
			}
		}
	}
}

func getLogName() string {
	logName := time.Now().Format("2006-01-02") + ".log"
	logCount := 1

	for {
		if _, err := os.Stat("./logs/" + logName); err == nil {
			logName = time.Now().Format("2006-01-02") + "_" + strconv.Itoa(logCount) + ".log"
			logCount++
		} else {
			break
		}
	}

	return logName
}

func NewLogger() (*zap.SugaredLogger, error) {
	cfg := configs.GetConfigs()

	deleteOldLogs(cfg.LogRemoveDays)

	logName := getLogName()
	path, _ := filepath.Abs("./logs/" + logName)

	logBuilder := zap.NewProductionConfig()
	logBuilder.OutputPaths = []string{path, "stderr"}
	logBuilder.EncoderConfig.CallerKey = zapcore.OmitKey

	logger, err := logBuilder.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
