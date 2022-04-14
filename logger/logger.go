package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// error logger
var allLogger *zap.SugaredLogger
var warnLogger *zap.SugaredLogger
var errorLogger *zap.SugaredLogger
var g_strApp string
var g_strLevel string
var g_time time.Time
var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func Init(strApp, strLevel string) {
	g_strApp = strApp
	g_strLevel = strLevel
	g_time = time.Now()
	fileName := strApp + g_time.Format("2006-01-02") + ".log"
	level := getLoggerLevel(strLevel)
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   1 << 30, //1G
		LocalTime: true,
		Compress:  true,
	})
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(level))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	// 全部日志
	allLogger = logger.Sugar()

	// 错误日志
	fileName = strApp + g_time.Format("2006-01-02") + "_warn.log"
	level = getLoggerLevel(strLevel)
	syncWarnWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   1 << 30, //1G
		LocalTime: true,
		Compress:  true,
	})
	encoderWarn := zap.NewProductionEncoderConfig()
	encoderWarn.EncodeTime = zapcore.ISO8601TimeEncoder
	coreWarn := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderWarn), syncWarnWriter, zap.NewAtomicLevelAt(level))
	loggerWarn := zap.New(coreWarn, zap.AddCaller(), zap.AddCallerSkip(1))
	warnLogger = loggerWarn.Sugar()

	// 错误日志
	fileName = strApp + g_time.Format("2006-01-02") + "_error.log"
	level = getLoggerLevel(strLevel)
	syncErrorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   1 << 30, //1G
		LocalTime: true,
		Compress:  true,
	})
	encoderError := zap.NewProductionEncoderConfig()
	encoderError.EncodeTime = zapcore.ISO8601TimeEncoder
	coreError := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderError), syncErrorWriter, zap.NewAtomicLevelAt(level))
	loggerError := zap.New(coreError, zap.AddCaller(), zap.AddCallerSkip(1))
	errorLogger = loggerError.Sugar()
}

func Debug(args ...interface{}) {
	if g_time.Day() != time.Now().Day() {
		Init(g_strApp, g_strLevel)
	}
	allLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	if g_time.Day() != time.Now().Day() {
		Init(g_strApp, g_strLevel)
	}
	allLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	allLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	allLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	allLogger.Warn(args...)
	warnLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	allLogger.Warnf(template, args...)
	warnLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	allLogger.Error(args...)
	errorLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	allLogger.Errorf(template, args...)
	errorLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	allLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	allLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	allLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	allLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	allLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	allLogger.Fatalf(template, args...)
}
