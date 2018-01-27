package common

import (
	"fmt"
	"log/syslog"

	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

// Deprecated.
func GetLocalSyslogLogger(logPriority syslog.Priority) (*log.Logger, error) {
	newLog := log.New()
	hook, err := lSyslog.NewSyslogHook("", "", logPriority, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get syslog hook: %s", err)
	}
	newLog.Hooks.Add(hook)

	return newLog, nil
}

func NewLocalSyslogLogger() (*log.Logger, error) {
	newLog := log.New()
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_DEBUG, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get syslog hook: %s", err)
	}
	newLog.Hooks.Add(hook)

	customFormatter := new(log.JSONFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	newLog.Formatter = customFormatter

	return newLog, nil
}

func ConvertLogLevel2Str(logLevel uint32) string {
	var logStr string
	switch logLevel {
	case 0:
		logStr = "panic"
	case 1:
		logStr = "fatal"
	case 2:
		logStr = "error"
	case 3:
		logStr = "warn"
	case 4:
		logStr = "info"
	default:
		logStr = "debug"
	}

	return logStr
}
