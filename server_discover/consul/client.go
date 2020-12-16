package consul

import (
	consul "github.com/hashicorp/consul/api"
)

// Client Consul API设置三个方法
type Client interface {
	// 服务注册
	Register(r *consul.AgentServiceRegistration) error

	// 服务下线
	Deregister(r *consul.AgentServiceRegistration) error

	// Service
	Service(service, tag string, passingOnly bool, queryOpts *consul.QueryOptions) ([]*consul.ServiceEntry, *consul.QueryMeta, error)
}

type client struct {
	consul *consul.Client
}

// NewClient 包装一个consul客户端返回
func NewClient(c *consul.Client) Client {
	return &client{consul: c}
}
func (c *client) Register(r *consul.AgentServiceRegistration) error {
	return c.consul.Agent().ServiceRegister(r)
}

func (c *client) Deregister(r *consul.AgentServiceRegistration) error {
	return c.consul.Agent().ServiceDeregister(r.ID)
}

func (c *client) Service(service, tag string, passingOnly bool, queryOpts *consul.QueryOptions) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	return c.consul.Health().Service(service, tag, passingOnly, queryOpts)
}
