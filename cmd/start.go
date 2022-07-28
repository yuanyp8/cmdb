package cmd

import (
	"fmt"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/spf13/cobra"
	"github.com/yuanyp8/cmdb/apps"
	_ "github.com/yuanyp8/cmdb/apps/all"
	"github.com/yuanyp8/cmdb/conf"
	"github.com/yuanyp8/cmdb/protocol"
	"os"
	"os/signal"
	"syscall"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start cmdb API",
	Long:  "start the  demo service's backend API",
	RunE:  start,
}

var confFile string

func start(cmd *cobra.Command, args []string) error {
	// 加载程序配置文件
	if err := conf.LoadConfig(confFile); err != nil {
		return err
	}

	// 初始化全局日志logger
	if err := loadGlobalLogger(); err != nil {
		return err
	}

	// 初始化HostServiceImpl
	apps.InitApps()

	svc := newManager()

	// 创建一个监听信号的管道
	ch := make(chan os.Signal, 1)
	defer close(ch)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)

	go svc.WaitStop(ch)

	return svc.Start()
}

func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)

	// 加载配置文件里面声明的log来配置全局logger
	configLog := conf.C().Log

	// 加载日志等级
	logLevel, err := zap.NewLevel(string(configLog.Level))
	if err != nil {
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = logLevel
		logInitMsg = fmt.Sprintf("log level: %s", logLevel)
	}

	// 使用默认配置初始化全局Logger
	zapConfig := zap.DefaultConfig()
	// 配置log level
	zapConfig.Level = level
	// 程序每启动一次，不必每次都生成一个日志文件
	zapConfig.Files.RotateOnStartup = false
	// 配置文件输出方式
	switch configLog.To {
	case conf.ToStdout:
		// 把日志打印到标准输出
		zapConfig.ToStderr = true
		// 并没在把日志输入输出到文件
		zapConfig.ToFiles = false
	case conf.ToFile:
		zapConfig.Files.Name = "api.log"
		zapConfig.Files.Path = configLog.PathDir
	}

	// 配置日志的输出格式
	switch configLog.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}

	// 把配置用于全局Logger
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

type manager struct {
	http *protocol.HttpService
	l    logger.Logger
}

func newManager() *manager {
	return &manager{
		http: protocol.NewHttpService(),
		l:    zap.L().Named("CLI"),
	}
}

func (m *manager) Start() error {
	return m.http.Start()
}

func (m *manager) WaitStop(ch <-chan os.Signal) {
	for v := range ch {
		switch v {
		// 后续我们还需要增加下对不同信号的不同处理逻辑
		case os.Kill:
			//TODO
		default:
			m.l.Infof("received signal: %s", v)
			m.http.Stop()
		}
	}
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "cmdb api 配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
