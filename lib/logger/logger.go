package logger

import (
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Log *zap.SugaredLogger
}

var log Logger

const (
	TraceLevel   zapcore.Level = zapcore.InfoLevel - 1 // Define a custom level
	VerboseLevel zapcore.Level = zapcore.InfoLevel - 2 // Define a custom level
)

func init() {
	if log.Log != nil {
		return
	}

	// Initialize the logger with a default verbosity level
	log.Log = createLogger(0)
}

func New(verbosity int) *zap.SugaredLogger {
	log.Log = createLogger(verbosity)

	log.Log.Debug("Logger initialized at level", log.Log.Level())

	return log.Log
}

// func NewDevZapLogger(lvl zapcore.LevelEnabler) *zap.Logger {
// 	encCfg := zap.NewDevelopmentEncoderConfig()
// 	encCfg.EncodeLevel = capitalColorLevelEncoder
// 	zapCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), zapcore.AddSync(os.Stderr), lvl)
// 	return zap.New(zapCore)
// }

func capitalColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if level == TraceLevel {
		enc.AppendString(color.CyanString("TRACE"))
		return
	}

	if level == VerboseLevel {
		enc.AppendString(color.GreenString("VERBOSE"))
		return
	}

	zapcore.CapitalColorLevelEncoder(level, enc)
}

func createLogger(verbosity int) *zap.SugaredLogger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = capitalColorLevelEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(getLevel(verbosity)),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	// zap.ReplaceGlobals(NewDevZapLogger(zap.NewAtomicLevelAt(verboseLevel)))
	// zap.ReplaceGlobals(NewDevZapLogger(zap.NewAtomicLevelAt(traceLevel)))

	return zap.Must(config.Build()).Sugar()
}

// Helper Functions

func Debug(msg string, args ...any) {
	if log.Log != nil {
		log.Log.Debugf(msg, args...)
	}
}

func Trace(msg string, args ...any) {
	if log.Log != nil {
		log.Log.Logf(TraceLevel, msg, args...)
	}
}

func Verbose(msg string, args ...any) {
	if log.Log != nil {
		log.Log.Logf(VerboseLevel, msg, args...)
	}
}

func Info(msg string, args ...any) {
	if log.Log != nil {
		log.Log.Infof(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if log.Log != nil {
		log.Log.Warnf(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if log.Log != nil {
		log.Log.Errorf(msg, args...)
	}
}

func Fatal(msg string, args ...any) {
	if log.Log != nil {
		log.Log.Fatalf(msg, args...)
	}
}

func Panic(err error) {
	if log.Log != nil {
		log.Log.DPanic(err)
	}
}

// Private functions

func getLevel(verbosity int) zapcore.Level {
	if verbosity >= 4 {
		return zapcore.DebugLevel
	}

	if verbosity >= 3 {
		return TraceLevel
	}

	if verbosity >= 2 {
		return VerboseLevel
	}

	if verbosity >= 1 {
		return zapcore.InfoLevel
	}

	return zapcore.ErrorLevel
}
