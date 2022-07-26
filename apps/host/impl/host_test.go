package impl_test

import (
	"context"
	"fmt"
	"github.com/infraboard/mcube/logger/zap"
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

	ins.Id = "ins-1"
	ins.Name = "test1"
	ins.Region = "cn-hangzhou"
	ins.Description = "接口测试"
	ins.Type = "sm"
	ins.CPU = 2
	ins.Memory = 4096

	ins, err := service.CreateHost(context.Background(), ins)
	if should.NoError(err) {
		fmt.Println(ins)
	}
}

func TestHostServiceImpl_QueryHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewQueryHostRequest()
	req.Keywords = "接口测试"
	set, err := service.QueryHost(context.Background(), req)
	if should.NoError(err) {
		for i := range set.Items {
			fmt.Println(set.Items[i].Id)
		}
	}
}

func TestHostServiceImpl_DescribeHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewDescribeHostRequestWithId("ins-2")
	h, err := service.DescribeHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(h.Id, h.Name)
	}
}

func TestHostServiceImpl_UpdateHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewPutUpdateHostRequest("ins-2")
	req.Name = "更新测试02"
	req.Region = "rg 02"
	req.Type = "small"
	req.CPU = 1
	req.Memory = 2048
	req.Description = "测试更新"
	h, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(h.Id, h.Id, h.Description)
	}
}

func TestHostServiceImpl_DeleteHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewDeleteHostRequest("ins-2")
	h, err := service.DeleteHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(h.Id, h.Id, h.Description)
	}
}

func init() {
	// 初始化全局Logger，不默认打印是为了减少性能损耗
	zap.DevelopmentSetup()

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
