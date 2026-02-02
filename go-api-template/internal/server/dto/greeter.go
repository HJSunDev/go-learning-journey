// Package dto 定义 HTTP 请求/响应的数据传输对象
// DTO 用于接收 HTTP 请求并进行验证，然后转换为 Service 层使用的 Proto 类型
package dto

import (
	v1 "go-api-template/api/helloworld/v1"
)

// SayHelloRequest 是 POST /api/v1/greeter/say-hello 的请求体
// binding tag 定义验证规则，由 Gin 内置的 Validator 库执行验证
type SayHelloRequest struct {
	// Name 要问候的用户名称
	// required: 必填
	// min=1: 最小长度 1
	// max=100: 最大长度 100
	Name string `json:"name" binding:"required,min=1,max=100" example:"World"`
}

// ToProto 将 DTO 转换为 Proto 类型
// HTTP Handler 用 DTO 接收并验证请求后，调用此方法转换，再传给 Service 层
func (r *SayHelloRequest) ToProto() *v1.SayHelloRequest {
	return &v1.SayHelloRequest{
		Name: r.Name,
	}
}
