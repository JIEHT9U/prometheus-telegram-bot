package logger

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type Event struct {
	id    int
	value string
}

func (e Event) toString() string {
	return ""
}

var (
	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid arg value: %s => %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
)

type Logger struct {
	entryErr  *logrus.Entry
	entryInfo *logrus.Entry
}

// func (l *Logger) InvalidArg(name string) {
// 	l.entry.Errorf(invalidArgMessage.toString(), name)

// }
// func (l *Logger) InvalidArgValue(name string, value interface{}) {
// 	l.entry.WithField("arg."+name, value).Errorf(invalidArgValueMessage.toString(), name, value)
// }
// func (l *Logger) MissingArg(name string) {
// 	l.entry.Errorf(missingArgMessage.toString(), name)
// }

/* func (l *Logger) Info(args ...interface{}) {
	l.entryInfo.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.entryInfo.Infof(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.entryErr.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entryErr.Errorf(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.entryInfo.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entryInfo.Debugf(format, args...)
}
*/

func (l *Logger) ErrEntry() *logrus.Entry {
	return l.entryErr
}
func (l *Logger) InfoEntry() *logrus.Entry {
	return l.entryInfo
}

func (l *Logger) InfoCT(templateName string, chatID int64, args ...interface{}) {
	l.entryInfo.WithField("chat_id", chatID).WithField("template_name", templateName).Info(args...)
}
func (l *Logger) DebugfCT(templateName string, chatID int64, format string, args ...interface{}) {
	l.entryInfo.WithField("chat_id", chatID).WithField("template_name", templateName).Debugf(format, args...)
}

func (l *Logger) ErrorCT(templateName string, chatID int64, args ...interface{}) {
	l.entryInfo.WithField("chat_id", chatID).WithField("template_name", templateName).Error(args...)
}

func (l *Logger) ErrorfCT(templateName string, chatID int64, format string, args ...interface{}) {
	l.entryInfo.WithField("chat_id", chatID).WithField("template_name", templateName).Errorf(format, args...)
}

func (l *Logger) ReqError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("content-type", "application/javascript")
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{Error: err.Error()})
}

func New(debug, json bool) (*Logger, error) {
	var timestampFormat = "02-01-2006 15:04:05"
	var formatter logrus.Formatter

	if json {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
		}
	} else {
		formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: timestampFormat,
			DisableSorting:  true,
		}
	}

	hostname, error := os.Hostname()

	if error != nil {
		return &Logger{}, error
	}

	loggerInfo := logrus.New()
	loggerInfo.Formatter = formatter
	loggerInfo.Out = os.Stdout

	loggerError := logrus.New()
	loggerError.Formatter = formatter
	loggerError.Out = os.Stderr

	if debug {
		loggerError.Level = logrus.DebugLevel
		loggerInfo.Level = logrus.DebugLevel

	}

	return &Logger{
		entryInfo: loggerInfo.WithFields(logrus.Fields{
			"hostname": hostname,
		}),
		entryErr: loggerError.WithFields(logrus.Fields{
			"hostname": hostname,
		}),
	}, nil
}
