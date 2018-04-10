package log

import (
	"os"
	"syscall"
	//"fmt"
	//lcf "github.com/Robpol86/logrus-custom-formatter"
	"img-crawler/src/conf"

	"github.com/sirupsen/logrus"
)

type Level logrus.Level
type Fields logrus.Fields

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var logger *logrus.Logger
var fp *os.File

func init() {
	logger = logrus.New()

	var err error
	fp, err = os.OpenFile(conf.Config.Log_path,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Fatalf("error opening log file: %v", err)
	}

	defer fp.Close()

	/*
		lcf.WindowsEnableNativeANSI(true)
		template := fmt.Sprintf("[%%[shortLevelName]s %%[ascTime]s%s] %%-45[message]s%%[fields]s\n", getHostname())
		logger.Formatter = lcf.NewFormatter(template, nil)
	*/

	//logrus.SetOutput(fp)
	syscall.Dup2(int(fp.Fd()), 1) /* -- stdout */
	syscall.Dup2(int(fp.Fd()), 2) /* -- stderr */

	logger.Level = logrus.InfoLevel
}

func getHostname() (hostname string) {
	hostname = os.Getenv("HOSTNAME")
	if hostname != "" {
		hostname = " " + hostname
	}
	return
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

func WithFields(fields Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields(fields))
}

func WithError(err error) *logrus.Entry {
	return logger.WithError(err)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Panicln(args ...interface{}) {
	logger.Panicln(args...)
}

func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}

func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

func SetLevel(level Level) {
	logger.Level = logrus.Level(level)
}
