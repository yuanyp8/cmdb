// Package apps IOC层
package apps

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuanyp8/cmdb/apps/host"
)

var (
	HostService host.Service

	// 将多个实例进行管理，减少手工代码量
	// 在这里我们使用Interface 加断言来进行抽象
	implApps = map[string]ImplService{}
	ginApps  = map[string]GinService{}
)

type ImplService interface {
	Config()
	Name() string
}

// RegistryImpl 服务实例注册到svc map当中
func RegistryImpl(svc ImplService) {
	if _, ok := implApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}
	implApps[svc.Name()] = svc
	// 同时，我们也判断是否能注册到HostService
	if v, ok := svc.(host.Service); ok {
		HostService = v
	}
}

func InitApps() {
	for _, v := range implApps {
		v.Config()
	}
}

// GetImpl 如果指定了具体类型, 就导致没增加一种类型, 多一个Get方法
// func GetHostImpl(name string) host.Service
func GetImpl(name string) interface{} {
	for k, v := range implApps {
		if k == name {
			return v
		}
	}
	return nil
}

type GinService interface {
	Registry(r gin.IRouter)
	Config()
	Name() string
}

func RegistryGin(svc GinService) {
	if _, ok := ginApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}

	ginApps[svc.Name()] = svc
}

func InitGin(r gin.IRouter) {
	// 先初始化好所有对象
	for _, v := range ginApps {
		v.Config()
	}
	// 完成Http Handler的注册
	for _, v := range ginApps {
		v.Registry(r)
	}
}

// LoadedGinApps 已经加载完成的Gin App有哪些
func LoadedGinApps() (names []string) {
	for k := range ginApps {
		names = append(names, k)
	}
	return
}
