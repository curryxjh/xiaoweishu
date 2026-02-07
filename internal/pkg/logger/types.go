package logger

// 兼容性最好
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func LoggerExample() {
	var l Logger
	phone := "159****7246"
	l.Info("user not registerd, phone: %s", phone)
}

type Field struct {
	Key   string
	Value any
}

// 认同参数需要有名字就是用LoggerV1
type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	With(args ...Field) LoggerV1
}

func LoggerV1Example() {
	var l LoggerV1
	phone := "159****7246"
	l.Info("user not registerd", Field{Key: "phone", Value: phone})
}

// 有完善的代码评审流程就是用，否则不建议，因为不确定用户是否按照规定使用
type LoggerV2 interface {
	// Debug @param args必须是偶数，并且按照 key-value 来组织
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...Field) LoggerV2
}
