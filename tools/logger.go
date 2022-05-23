package tools

import (
	"fmt"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

var (
	Logger = &logger{}

	LowercaseLevelEncoder      = zapcore.LowercaseLevelEncoder
	LowercaseColorLevelEncoder = zapcore.LowercaseColorLevelEncoder
	CapitalLevelEncoder        = zapcore.CapitalLevelEncoder
	CapitalColorLevelEncoder   = zapcore.CapitalColorLevelEncoder
)

type logger struct {
	*zap.Logger
	Config *LoggerConfig
}

type LoggerConfig struct {
	Director      string               `json:"director"`
	Level         zapcore.Level        `json:"level"`
	ShowLine      bool                 `json:"showLine"`
	StacktraceKey string               `json:"stacktraceKey"`
	LinkName      string               `json:"linkName"`
	LogInConsole  bool                 `json:"logInConsole"`
	Format        string               `json:"format"`
	EncodeLevel   zapcore.LevelEncoder `json:"encodeLevel"`
	Prefix        string               `json:"prefix"`
}

// InitLogger @Title 初始化日志工具
func InitLogger(config *LoggerConfig) *logger {
	Logger.Config = config
	if ok, _ := Logger.pathExists(config.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", config.Director)
		_ = os.Mkdir(config.Director, os.ModePerm)
	}

	if config.Level == zap.DebugLevel || config.Level == zap.ErrorLevel {
		zap.New(Logger.getEncoderCore(), zap.AddStacktrace(config.Level))
		Logger.Logger = zap.New(Logger.getEncoderCore(), zap.AddStacktrace(config.Level))
	} else {
		Logger.Logger = zap.New(Logger.getEncoderCore())
	}
	if config.ShowLine {
		Logger.Logger = Logger.Logger.WithOptions(zap.AddCaller())
	}
	return Logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func (l *logger) getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  l.Config.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    l.Config.EncodeLevel,
		EncodeTime:     l.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if config.EncodeLevel == nil {
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func (l *logger) getEncoder() zapcore.Encoder {
	if l.Config.Format == "json" {
		return zapcore.NewJSONEncoder(l.getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(l.getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func (l *logger) getEncoderCore() (core zapcore.Core) {
	writer, err := l.getWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(l.getEncoder(), writer, l.Config.Level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func (l *logger) CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(l.Config.Prefix + "2006/01/02 - 15:04:05.000"))
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: PathExists
//@description: 文件目录是否存在
//@param: path string
//@return: bool, error
func (l *logger) pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (l *logger) getWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(l.Config.Director, "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName(l.Config.LinkName),
		zaprotatelogs.WithMaxAge(time.Hour*24*7),
		zaprotatelogs.WithRotationTime(time.Hour*24),
	)
	if l.Config.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
