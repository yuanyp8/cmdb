package impl_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuanyp8/cmdb/apps/host"
	"github.com/yuanyp8/cmdb/apps/host/impl"
	"github.com/yuanyp8/cmdb/conf"
	"testing"
)

var service host.Service

func TestHostServiceImpl_CreateHost(t *testing.T) {
	should := assert.New(t)
	ins := host.NewHost()

	ins.Id = "ins-01"
	ins.Name = "test"
	ins.Region = "cn-hangzhou"
	ins.Type = "sm1"
	ins.CPU = 1
	ins.Memory = 2048

	ins, err := service.CreateHost(context.Background(), ins)
	if should.NoError(err) {
		fmt.Println(ins)
	}
}

func init() {
	// 加载配置文件
	cfg, err := conf.LoadConfigFromToml("../../../conf/testdata/config-good.toml")
	if err != nil {
		panic(err)
	}
	if err := cfg.InitGlobal(); err != nil {
		panic(err)
	}

	// 初始化全局Logger
	service = impl.NewHostServiceImpl()
}
