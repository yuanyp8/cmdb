package conf

// LogFormat 日志输出格式
type LogFormat string

const (
	// TextFormat 文本格式
	TextFormat = LogFormat("text")
	// JSONFormat json格式
	JSONFormat = LogFormat("json")
)

// LogTo 日志记录到哪儿
type LogTo string

const (
	// ToFile 保存到文件
	ToFile = LogTo("file")
	// ToStdout 打印到标准输出
	ToStdout = LogTo("stdout")
)

// LogLevel 定义日志等级
type LogLevel string

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warning"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
)
