package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// 存储应用程序自身相关配置
type app struct {
	// app name
	Name string `toml:"name" env:"APP_NAME"`
	// 应用程序在http protocol下运行的相关配置
	HttpHost string `toml:"http_host" env:"APP_HTTP_HOST"`
	HttpPort string `toml:"http_port" env:"APP_HTTP_PORT"`
	// tls认证的key文件
	EncryptKey string `toml:"encrypt_key" env:"APP_ENCRYPT_KEY"`
}

// 构造函数
func newDefaultAPP() *app {
	return &app{
		Name:       "cmdb",
		HttpHost:   "127.0.0.1",
		HttpPort:   "8060",
		EncryptKey: "defualt app encrypt key",
	}
}

// HTTPAddr 用于ListenAndServe时拼接
func (a *app) HTTPAddr() string {
	return fmt.Sprintf("%s:%s", a.HttpHost, a.HttpPort)
}

// 应用日志相关配置
type log struct {
	Level   LogLevel  `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}

func newDefaultLog() *log {
	return &log{
		Level:  DebugLevel,
		Format: TextFormat,
		To:     ToStdout,
	}
}

type mySQL struct {
	Host     string `toml:"host" env:"MYSQL_HOST"`
	Port     string `toml:"port" env:"MYSQL_PORT"`
	UserName string `toml:"username" env:"MYSQL_USERNAME"`
	Password string `toml:"password" env:"MYSQL_PASSWORD"`
	Database string `toml:"database" env:"MYSQL_DATABASE"`
	/*使用的是MySQL连接池，所以需要进行一些初始化配置*/
	MaxOpenConn int `toml:"max_open_conn" env:"MYSQL_MAX_OPEN_CONN"`
	MaxIdleConn int `toml:"max_idle_conn" env:"MYSQL_MAX_IDLE_CONN"`
	MaxLifeTime int `toml:"max_life_time" env:"MYSQL_MAX_LIFE_TIME"`
	MaxIdleTime int `toml:"max_idle_time" env:"MYSQL_MAX_idle_TIME"`

	lock sync.Mutex
}

var (
	db *sql.DB
)

func newDefaultMySQL() *mySQL {
	return &mySQL{
		Database:    "cmdb",
		Host:        "127.0.0.1",
		Port:        "3306",
		MaxOpenConn: 200,
		MaxIdleConn: 50,
		MaxLifeTime: 1800,
		MaxIdleTime: 600,
	}
}

// initPool 初始化连接池
func (m *mySQL) initPool() error {
	var err error

	dsnTemplate := "%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true"
	dsn := fmt.Sprintf(dsnTemplate, m.UserName, m.Password, m.Host, m.Port, m.Database)

	db2, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	// configuration MySQL pool
	db2.SetMaxOpenConns(m.MaxOpenConn)
	db2.SetMaxIdleConns(m.MaxIdleConn)
	db2.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	db2.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db2.PingContext(ctx); err != nil {
		return fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	db = db2
	return err
}

func (m *mySQL) GetDB() (*sql.DB, error) {
	// 加载全局实例，保证连接池只初始化一次
	m.lock.Lock()
	defer m.lock.Unlock()

	if db == nil {
		if err := m.initPool(); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// MustGetDB 确保数据库连接不会异常时才调用，如果调用异常则直接panic
func (m *mySQL) MustGetDB() *sql.DB {
	db, err := m.GetDB()
	if err != nil {
		panic(err)
	}
	return db
}

// Config 应用配置
// 通过封装为一个对象, 来与外部配置进行对接
type Config struct {
	App   *app   `toml:"app"`
	Log   *log   `toml:"log"`
	MySQL *mySQL `toml:"mysql"`
}

func newConfig() *Config {
	return &Config{
		App:   newDefaultAPP(),
		Log:   newDefaultLog(),
		MySQL: newDefaultMySQL(),
	}
}
