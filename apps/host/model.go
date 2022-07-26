package host

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"time"
)

type Vendor int

const (
	// PRIVATE_IDC 枚举的默认值
	PRIVATE_IDC Vendor = iota
	// ALIYUN 阿里云
	ALIYUN
	// TXYUN 腾讯云
	TXYUN
)

// Host 模型的定义
type Host struct {
	// 资源公共属性部分
	*Resource
	// 资源独有属性部分
	*Describe
}

func NewHost() *Host {
	return &Host{
		Resource: &Resource{},
		Describe: &Describe{},
	}
}

type Resource struct {
	Id          string `json:"id"  validate:"required"`     // 全局唯一Id
	Vendor      Vendor `json:"vendor"`                      // 厂商
	Region      string `json:"region"  validate:"required"` // 地域
	CreateAt    int64  `json:"create_at"`                   // 创建时间
	ExpireAt    int64  `json:"expire_at"`                   // 过期时间
	Type        string `json:"type"  validate:"required"`   // 规格
	Name        string `json:"name"  validate:"required"`   // 名称
	Description string `json:"description"`                 // 描述
	Status      string `json:"status"`                      // 服务商中的状态
	// Tags        map[string]string `json:"tags"`                        // 标签
	UpdateAt  int64  `json:"update_at"`  // 更新时间
	SyncAt    int64  `json:"sync_at"`    // 同步时间
	Account   string `json:"account"`    // 资源的所属账号
	PublicIP  string `json:"public_ip"`  // 公网IP
	PrivateIP string `json:"private_ip"` // 内网IP
}

type Describe struct {
	CPU          int    `json:"cpu" validate:"required"`    // 核数
	Memory       int    `json:"memory" validate:"required"` // 内存
	GPUAmount    int    `json:"gpu_amount"`                 // GPU数量
	GPUSpec      string `json:"gpu_spec"`                   // GPU类型
	OSType       string `json:"os_type"`                    // 操作系统类型，分为Windows和Linux
	OSName       string `json:"os_name"`                    // 操作系统名称
	SerialNumber string `json:"serial_number"`              // 序列号
}

type HostSet struct {
	Total int     `json:"total"`
	Items []*Host `json:"items"`
}

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

// Add 将向主机列表增加主机的操作封装成一个函数，使得代码变得优雅
func (s *HostSet) Add(item *Host) {
	s.Items = append(s.Items, item)
}

// Put 对象全量更新
func (h *Host) Put(obj *Host) error {
	if obj.Id != h.Id {
		return fmt.Errorf("id not equal")
	}

	*h.Resource = *obj.Resource
	*h.Describe = *obj.Describe
	return nil
}

// Patch 对象的局部更新
func (h *Host) Patch(obj *Host) error {
	// if obj.Name != "" {
	// 	h.Name = obj.Name
	// }
	// if obj.CPU != 0 {
	// 	h.CPU = obj.CPU
	// }
	// 比如 obj.A  obj.B  只想修改obj.B该属性
	return mergo.Merge(h, obj)
}

var validate = validator.New()

// Validate 字段合法性校验
func (h *Host) Validate() error {
	return validate.Struct(h)
}

func (h *Host) InjectDefault() {
	if h.CreateAt == 0 {
		h.CreateAt = time.Now().UnixMilli()
	}
}

// QueryHostRequest 入参结果体
// 由于Query是在前端发出的请求，假设我们查到了1万条数据,不能全部都返回给前端，这里要做一个分页操作
type QueryHostRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	Keywords   string `json:"kws"`
}

func NewQueryHostRequest() *QueryHostRequest {
	return &QueryHostRequest{
		PageSize:   20,
		PageNumber: 1,
	}
}

func (req *QueryHostRequest) GetPageSize() uint {
	return uint(req.PageSize)
}

func (req *QueryHostRequest) OffSet() int64 {
	return int64((req.PageNumber - 1) * req.PageSize)
}

// DescribeHostRequest 查询主机详情的时候，我们可以通过主键进行关联查询，后面我们会详细说明
type DescribeHostRequest struct {
	Id string
}

func NewDescribeHostRequestWithId(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		Id: id,
	}
}

// UPDATE_MODE 这里的Update分为全量更新或者部分更新
type UPDATE_MODE string

const (
	// UPDATE_MODE_PUT 全量更新
	UPDATE_MODE_PUT UPDATE_MODE = "put"
	// UPDATE_MODE_PATCH 局部更新
	UPDATE_MODE_PATCH UPDATE_MODE = "patch"
)

type UpdateHostRequest struct {
	UpdateMode UPDATE_MODE `json:"update_mode"`
	*Host
}

func NewPutUpdateHostRequest(id string) *UpdateHostRequest {
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		UpdateMode: UPDATE_MODE_PUT,
		Host:       h,
	}
}

func NewPatchUpdateHostRequest(id string) *UpdateHostRequest {
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		UpdateMode: UPDATE_MODE_PATCH,
		Host:       h,
	}
}

type DeleteHostRequest struct {
	Id string
}

func NewDeleteHostRequest(id string) *DeleteHostRequest {
	return &DeleteHostRequest{id}
}
