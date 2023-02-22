package logger

import log "github.com/sirupsen/logrus"

func init() {
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "timestamp",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "message",
		},
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)
}

func Info(format string, args ...interface{}) {
	log.WithField("app", "git-workflows").Infof(format, args...)
}

func Debug(format string, args ...interface{}) {
	log.WithField("app", "git-workflows").Debugf(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.WithField("app", "git-workflows").Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	log.WithField("app", "git-workflows").Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.WithField("app", "git-workflows").Fatalf(format, args...)
}

func EnableDebug() {
	log.SetLevel(log.DebugLevel)
}
