package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yuanyp8/cmdb/apps"
	"github.com/yuanyp8/cmdb/apps/host"
)

// Handler 该实体类把内部的Service接口通过http协议暴露出去
type Handler struct {
	svc host.Service
}

func (h *Handler) Config() {
	// 从IOC里面获取HostService的实例对象
	h.svc = apps.GetImpl(host.AppName).(host.Service)
}

func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/hosts", h.createHost)
	r.GET("/hosts", h.queryHost)
	r.GET("/hosts/:id", h.describeHost)
	r.PUT("/hosts/:id", h.putHost)
	r.PATCH("/hosts/:id", h.patchHost)
	r.DELETE("/hosts/:id", h.deleteHost)
}

func (h *Handler) Name() string {
	return host.AppName
}

var handler = &Handler{}

// 完成Http Handler注册
func init() {
	apps.RegistryGin(handler)
}
