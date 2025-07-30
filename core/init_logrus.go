package core

import (
	"blogx_server/global"
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// 颜色常量，用于在控制台输出中为不同级别的日志着色。
// 它们是标准的ANSI转义序列颜色代码。
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

// LogFormatter 是一个自定义的日志格式化器。
// 通过实现logrus.Formatter接口，我们可以完全控制日志的输出格式。
type LogFormatter struct{}

// Format 是logrus.Formatter接口的核心方法。
// 当logrus需要格式化一条日志记录时，就会调用这个方法。
// entry 参数包含了日志的所有信息（级别、消息、时间、调用者信息等）。
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 1. 根据日志级别选择不同的颜色。
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	// 2. 创建或复用一个缓冲区，用于拼接最终的日志字符串。
	// 这种复用机制可以减少内存分配，提升性能。
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 3. 格式化时间戳为"年-月-日 时:分:秒"的格式。
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// 4. 检查日志条目是否包含了调用者信息（文件名、行号、函数名）。
	if entry.HasCaller() {
		// 提取函数名。
		funcVal := entry.Caller.Function
		// 提取文件名和行号，并格式化为 "文件名:行号" 的形式。
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		// 按照自定义的格式将所有信息写入缓冲区。
		// 格式为：[时间] [颜色][级别][颜色重置] 文件:行号 函数名 日志消息
		// \x1b[?m 是用于控制终端颜色的ANSI转义序列。
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		// 如果没有调用者信息，则使用简化的格式。
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s\n", timestamp, levelColor, entry.Level, entry.Message)
	}
	// 5. 返回格式化好的日志字节流。
	return b.Bytes(), nil
}

// FileDateHook 是一个自定义的logrus钩子（Hook），用于实现日志的按日分割归档。
// 它实现了logrus.Hook接口。
type FileDateHook struct {
	file     *os.File // 当前正在写入的日志文件句柄
	logPath  string   // 日志文件的根目录
	fileDate string   // 当前日志文件对应的日期，用于判断是否需要创建新文件
	appName  string   // 应用程序名称，用作日志文件名的一部分
}

// Levels 告诉logrus，这个Hook关心哪些日志级别。
// 返回logrus.AllLevels意味着任何级别的日志都会触发这个Hook。
func (hook *FileDateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 是logrus.Hook接口的核心方法。
// 每当有符合Levels()定义的日志产生时，logrus就会调用这个方法。
// 注意：这里的hook是值传递，意味着在Fire方法内部对hook字段的修改（如hook.fileDate）不会影响到原始的hook实例。
// 这会导致每次跨天后，每条新日志都会触发一次文件切换，因为fileDate永远是旧的。
// 正确的做法是使用指针接收器 `func (hook *FileDateHook) Fire...`，但这会修改代码，此处仅作说明。
func (hook *FileDateHook) Fire(entry *logrus.Entry) error {
	// 1. 获取当前日志的日期（"年-月-日"）和格式化后的完整日志行。
	timer := entry.Time.Format("2006-01-02")
	line, _ := entry.String() // entry.String() 会调用上面我们定义的Format方法

	// 2. 检查当前日志的日期是否与Hook中记录的日期一致。
	if hook.fileDate == timer {
		// 如果日期一致，直接将日志写入当前文件。
		hook.file.Write([]byte(line))
		return nil
	}

	// 3. 如果日期不一致（即跨天了），则执行日志轮转操作。
	// 首先关闭旧的日志文件。
	hook.file.Close()
	// 创建一个新的以当天日期命名的子目录。
	os.MkdirAll(fmt.Sprintf("%s/%s", hook.logPath, timer), os.ModePerm)
	// 拼接出新的日志文件名。
	filename := fmt.Sprintf("%s/%s/%s.log", hook.logPath, timer, hook.appName)

	// 打开（或创建）新的日志文件以进行写入。
	hook.file, _ = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	// 更新Hook中记录的日期为当前日期。
	hook.fileDate = timer
	// 将当前这条日志写入新的日志文件。
	hook.file.Write([]byte(line))
	return nil
}

// InitFile 负责初始化文件日志记录。
// 它创建初始的日志目录和文件，并设置FileDateHook。
func InitFile(logPath, appName string) {
	// 1. 获取当前日期，用于创建初始的日志目录和文件名。
	fileDate := time.Now().Format("2006-01-02")
	// 创建当天的日志目录，例如 "logs/2023-10-27"。
	// os.MkdirAll 如果目录已存在，不会报错。
	err := os.MkdirAll(fmt.Sprintf("%s/%s", logPath, fileDate), os.ModePerm)
	if err != nil {
		logrus.Error(err)
		return
	}

	// 2. 拼接出当天的日志文件名。
	filename := fmt.Sprintf("%s/%s/%s.log", logPath, fileDate, appName)
	// 打开或创建当天的日志文件。
	// O_APPEND 表示以追加模式写入，O_CREATE 表示如果文件不存在则创建。
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		logrus.Error(err)
		return
	}
	// 3. 创建并初始化FileDateHook实例。
	fileHook := FileDateHook{file, logPath, fileDate, appName}
	// 4. 将创建好的Hook添加到logrus中。
	// 这里传递的是指针 &fileHook，这非常关键，确保了logrus持有了对我们fileHook实例的引用。
	// 因此，当Fire方法被调用时，它是在同一个实例上操作，使得文件句柄和日期的状态可以在多次调用间保持。
	logrus.AddHook(&fileHook)
}

// InitLogrus 是日志系统的总入口和总配置函数。
func InitLogrus() { //新建一个实例
	// 设置日志输出到标准输出（控制台）。
	logrus.SetOutput(os.Stdout)
	// 开启调用者信息（文件名、行号、函数名）的报告。
	logrus.SetReportCaller(true)
	// 设置使用我们自定义的LogFormatter进行日志格式化。
	logrus.SetFormatter(&LogFormatter{})
	// 设置日志记录的最低级别，低于此级别的日志将不会被记录。
	logrus.SetLevel(logrus.DebugLevel)
	// 从全局配置中获取日志目录和应用名。
	l := global.Config.Log
	// 调用InitFile来启动文件日志记录和按日轮转功能。
	InitFile(l.Dir, l.App)
}
