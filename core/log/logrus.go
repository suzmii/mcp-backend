package log

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const pkgPath = "core/log"

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"

	TimeColor = Blue    // Blue
	FuncColor = Magenta // Magenta

	// BackRed     = "\033[41m"
	// BackGreen   = "\033[42m"
	// BackYellow  = "\033[43m"
	// BackBlue    = "\033[44m"
	// BackMagenta = "\033[45m"
	// BackCyan    = "\033[46m"
)

func init() {
	logrus.SetFormatter(&formatter{})
}

type formatter struct{}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b strings.Builder

	// 日志级别
	b.WriteString(levelColor(entry.Level))
	b.WriteString(shotLevel(entry.Level))
	b.WriteString(Reset)

	// 时间
	b.WriteString(TimeColor)
	b.WriteByte('[')
	b.WriteString(time.Now().Format("2006-01-02 15:04:05.000"))
	b.WriteString("]")
	b.WriteString(Reset)
	b.WriteByte(' ')

	// caller file and func
	funcName, file, line := findCaller()
	b.WriteString(FuncColor)
	b.WriteString(file)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(line))
	b.WriteByte(':')
	b.WriteString(funcName)
	b.WriteString("():\n")
	b.WriteString(Reset)

	// 日志内容
	b.WriteString(levelColor(entry.Level))
	b.WriteString(entry.Message)
	b.WriteString(Reset)
	b.WriteByte('\n')

	// 如果是严重错误，打印堆栈

	if entry.Level <= logrus.ErrorLevel {
		stackBuf := make([]byte, 4096)
		n := runtime.Stack(stackBuf, false)
		stackLines := strings.Split(string(stackBuf[:n]), "\n")

		b.WriteString(levelColor(entry.Level))
		for i := 0; i < len(stackLines)-1; i += 2 { // 一次走两行：第一行是 "函数名"，第二行是 "文件名:行号"
			funcName := strings.TrimSpace(stackLines[i])
			fileLine := strings.TrimSpace(stackLines[i+1])

			if isInternal(funcName) {
				continue
			}

			b.WriteString(funcName)
			b.WriteByte('\n')
			b.WriteString(fileLine)
			b.WriteByte('\n')
		}
		b.WriteString(Reset)
	}

	return []byte(b.String()), nil
}
func findCaller() (funcName, file string, line int) {
	for skip := 4; skip < 15; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			continue
		}
		if !isInternal(file) {
			funcName = runtime.FuncForPC(pc).Name()
			return funcName, file, line
		}
	}
	return "unknown", "unknown", 0
}

func isInternal(file string) bool {
	return strings.Contains(file, "runtime.") || strings.Contains(file, "sirupsen/logrus") || strings.Contains(file, pkgPath)
}

func levelColor(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return Magenta
	case logrus.DebugLevel:
		return Green
	case logrus.InfoLevel:
		return Cyan
	case logrus.WarnLevel:
		return Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return Red
	default:
		return Reset
	}
}

// func levelBackColor(level logrus.Level) string {
// 	switch level {
// 	case logrus.TraceLevel:
// 		return BackMagenta
// 	case logrus.DebugLevel:
// 		return BackGreen
// 	case logrus.InfoLevel:
// 		return BackCyan
// 	case logrus.WarnLevel:
// 		return BackYellow
// 	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
// 		return BackRed
// 	default:
// 		return Reset
// 	}
// }

func shotLevel(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return "[TRACE]"
	case logrus.DebugLevel:
		return "[DEBUG]"
	case logrus.InfoLevel:
		return "[INFO] "
	case logrus.WarnLevel:
		return "[WARN] "
	case logrus.ErrorLevel:
		return "[ERROR]"
	case logrus.FatalLevel:
		return "[FATAL]"
	case logrus.PanicLevel:
		return "[PANIC]"
	default:
		return "[DEFAULT]"
	}
}
