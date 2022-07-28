package http

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"
	"github.com/yuanyp8/cmdb/apps/host"
)

func (h *Handler) createHost(c *gin.Context) {
	ins := host.NewHost()

	// 将用户传过来的参数映射到结构体中
	if err := c.Bind(ins); err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 进行接口调用
	ins, err := h.svc.CreateHost(c.Request.Context(), ins)
	// 注册失败了
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	// 注册成功了
	response.Success(c.Writer, ins)
}

func (h *Handler) queryHost(c *gin.Context) {
	// 从http请求的query string 中获取参数
	req := host.NewQueryHostFromHTTP(c.Request)
	// 进行接口调用, 返回 肯定有成功或者失败
	set, err := h.svc.QueryHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, set)
}

func (h *Handler) describeHost(c *gin.Context) {
	req := host.NewDescribeHostRequestWithId(c.Param("id"))

	set, err := h.svc.DescribeHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}

func (h *Handler) putHost(c *gin.Context) {
	req := host.NewPutUpdateHostRequest(c.Param("id"))
	// 解析Body里面的数据
	if err := c.Bind(req.Host); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	req.Id = c.Param("id")

	// 进行接口调用
	host, err := h.svc.UpdateHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, host)
}

func (h *Handler) patchHost(c *gin.Context) {
	req := host.NewPatchUpdateHostRequest(c.Param("id"))

	if err := c.Bind(req.Host); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	req.Id = c.Param("id")

	host, err := h.svc.UpdateHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, host)
}

func (h *Handler) deleteHost(c *gin.Context) {
	req := host.NewDeleteHostRequest(c.Param("id"))

	ins, err := h.svc.DeleteHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, ins)
}
