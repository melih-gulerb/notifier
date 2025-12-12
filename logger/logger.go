package logger

import (
	"context"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type LogFields struct {
	UserID        string                 `json:"userId,omitempty"`
	CorrelationID string                 `json:"correlationId,omitempty"`
	LogData       map[string]interface{} `json:"logData,omitempty"`
}

type Logger struct {
	nrApp *newrelic.Application
	log   *logrus.Logger
}

type LogEntry struct {
	logger  *Logger
	ctx     context.Context
	level   logrus.Level
	message string
	fields  LogFields
	err     error
}

func (l *LogFields) WithCorrelationID(correlationID string) *LogFields {
	l.CorrelationID = correlationID
	return l
}

func (l *LogFields) WithUserId(userId string) *LogFields {
	l.UserID = userId
	return l
}

func (l *LogFields) WithLogData(logData map[string]interface{}) *LogFields {
	l.LogData = logData
	return l
}

func NewLogger(nrApp *newrelic.Application) *Logger {
	nrlogrusFormatter := nrlogrus.NewFormatter(nrApp, &logrus.TextFormatter{})
	log := logrus.New()
	log.SetFormatter(nrlogrusFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	return &Logger{
		nrApp: nrApp,
		log:   log,
	}
}

func (l *Logger) Info(ctx context.Context) *LogEntry {
	return &LogEntry{
		logger: l,
		ctx:    ctx,
		level:  logrus.InfoLevel,
		fields: LogFields{LogData: make(map[string]interface{})},
	}
}

func (l *Logger) Error(ctx context.Context) *LogEntry {
	return &LogEntry{
		logger: l,
		ctx:    ctx,
		level:  logrus.ErrorLevel,
		fields: LogFields{LogData: make(map[string]interface{})},
	}
}

func (l *Logger) Warn(ctx context.Context) *LogEntry {
	return &LogEntry{
		logger: l,
		ctx:    ctx,
		level:  logrus.WarnLevel,
		fields: LogFields{LogData: make(map[string]interface{})},
	}
}

func (l *Logger) Debug(ctx context.Context) *LogEntry {
	return &LogEntry{
		logger: l,
		ctx:    ctx,
		level:  logrus.DebugLevel,
		fields: LogFields{LogData: make(map[string]interface{})},
	}
}

func (e *LogEntry) WithCorrelationId(correlationId string) *LogEntry {
	e.fields.CorrelationID = correlationId
	return e
}

func (e *LogEntry) WithUserId(userId string) *LogEntry {
	if userId != "" {
		e.fields.UserID = userId
	}

	return e
}

func (e *LogEntry) WithLogData(key string, value interface{}) *LogEntry {
	if e.fields.LogData == nil {
		e.fields.LogData = make(map[string]interface{})
	}
	e.fields.LogData[key] = value
	return e
}

func (e *LogEntry) WithError(err error) *LogEntry {
	e.err = err
	if err != nil && e.fields.LogData == nil {
		e.fields.LogData = make(map[string]interface{})
	}
	if err != nil {
		e.fields.LogData["error"] = err.Error()
	}
	return e
}

func (e *LogEntry) Log(message string) {
	e.message = message
	e.logger.logWithLevel(e.ctx, e.level, e.message, e.fields)
}

func (l *Logger) logWithLevel(ctx context.Context, level logrus.Level, message string, fields LogFields) {
	logEntry := l.log.WithFields(logrus.Fields{
		"userId":        fields.UserID,
		"correlationId": fields.CorrelationID,
		"logData":       fields.LogData,
	})

	switch level {
	case logrus.InfoLevel:
		logEntry.Info(message)
	case logrus.ErrorLevel:
		logEntry.Error(message)
	case logrus.WarnLevel:
		logEntry.Warn(message)
	case logrus.DebugLevel:
		logEntry.Debug(message)
	default:
		panic("unhandled default case")
	}

	l.sendToNewRelic(ctx, level, message, fields)
}

func (l *Logger) sendToNewRelic(ctx context.Context, level logrus.Level, message string, fields LogFields) {
	if l.nrApp == nil {
		return
	}

	txn := newrelic.FromContext(ctx)
	if txn == nil {
		txn = l.nrApp.StartTransaction("background-log")
		defer txn.End()
	}

	attributes := map[string]interface{}{
		"message":       message,
		"level":         level.String(),
		"userId":        fields.UserID,
		"correlationId": fields.CorrelationID,
		"project":       "shieldbearer",
	}

	if fields.LogData != nil {
		for k, v := range fields.LogData {
			attributes["logData."+k] = v
		}
	}

	l.nrApp.RecordLog(newrelic.LogData{
		Message:    message,
		Severity:   level.String(),
		Timestamp:  time.Now().UTC().Add(time.Hour * 3).UnixMilli(),
		Attributes: attributes,
	})
}

func (l *Logger) SetLevel(level string) {
	switch level {
	case "debug":
		l.log.SetLevel(logrus.DebugLevel)
	case "info":
		l.log.SetLevel(logrus.InfoLevel)
	case "warn":
		l.log.SetLevel(logrus.WarnLevel)
	case "error":
		l.log.SetLevel(logrus.ErrorLevel)
	default:
		l.log.SetLevel(logrus.InfoLevel)
	}
}

func InitNewRelic(appName, licenseKey string) (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		return nil, err
	}

	return app, nil
}
