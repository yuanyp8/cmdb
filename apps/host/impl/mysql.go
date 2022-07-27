package impl

import (
	"database/sql"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/yuanyp8/cmdb/apps"
	"github.com/yuanyp8/cmdb/apps/host"
	"github.com/yuanyp8/cmdb/conf"
)

// 接口实现的静态检查
// 这样写, 会造成 conf.C()并准备好, 造成conf.C().MySQL.GetDB()该方法panic
// var impl = NewHostServiceImpl()

var impl = &HostServiceImpl{}

type HostServiceImpl struct {
	l  logger.Logger
	db *sql.DB
}

func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		l:  zap.L().Named("Host"),
		db: conf.C().MySQL.MustGetDB(),
	}
}

// Name 服务的名称
func (i *HostServiceImpl) Name() string {
	return host.AppName
}

// Config 其作用是将实例注册进来
func (i *HostServiceImpl) Config() {
	i.l = zap.L().Named("Host")
	i.db = conf.C().MySQL.MustGetDB()
}

func init() {
	apps.RegistryImpl(impl)
}
