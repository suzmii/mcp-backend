package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogrusWriter struct{}

func (l *LogrusWriter) Alert(v any) {
	logrus.WithField("type", "alert").Error(v)
}

func (l *LogrusWriter) Close() error {
	return nil
}

func (l *LogrusWriter) Debug(v any, fields ...logx.LogField) {
	logrus.WithFields(convertFields(fields)).Debug(v)
}

func (l *LogrusWriter) Error(v any, fields ...logx.LogField) {
	logrus.WithFields(convertFields(fields)).Error(v)
}

func (l *LogrusWriter) Info(v any, fields ...logx.LogField) {
	logrus.WithFields(convertFields(fields)).Info(v)
}

func (l *LogrusWriter) Severe(v any) {
	logrus.WithField("type", "severe").Fatal(v)
}

func (l *LogrusWriter) Slow(v any, fields ...logx.LogField) {
	logrus.WithFields(convertFields(fields)).Warn("[SLOW] " + fmt.Sprint(v))
}

func (l *LogrusWriter) Stack(v any) {
	// 通常用于 panic trace
	logrus.WithField("stack", v).Error("stack trace")
}

func (l *LogrusWriter) Stat(v any, fields ...logx.LogField) {
	// 统计型日志一般为 info
	logrus.WithFields(convertFields(fields)).Info("[STAT] " + fmt.Sprint(v))
}

func convertFields(fields []logx.LogField) logrus.Fields {
	result := logrus.Fields{}
	for _, f := range fields {
		result[f.Key] = f.Value
	}
	return result
}

func InitLogx() {
	logx.SetWriter(&LogrusWriter{})
}
