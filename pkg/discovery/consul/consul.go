// @author: mongo
package consul

import (
	"fmt"
	"github.com/xulichen/halfway/pkg/discovery"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

// ConsulResource defines resource
type Resource struct {
	// consulClient consul 的 local agent
	consulClient *api.Client
	//
	serverConf      *discovery.ServerConfig
	claimedServices []string
	serviceConf     *discovery.ServiceConfig
}

// @todo 讲依赖声明加到初始化的地方
// New returns the point of ConsulResource
func New(cfg *discovery.ServerConfig) (*Resource, error) {
	consulConfig := &api.Config{
		Address: fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		Token:   cfg.Token,
	}
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("NewClient consul error\t%v \t address : %s", err, cfg.Address))
	}
	return &Resource{consulClient: client, serverConf: cfg}, nil
}

func (cr *Resource) WithServiceConfig(cfg *discovery.ServiceConfig) {
	cr.serviceConf = cfg
}

// RegisterService 注册服务
func (cr *Resource) RegisterService() error {
	//register consul
	agent := cr.consulClient.Agent()
	interval, deregister := 3*time.Second, 1*time.Second
	// 区分 rpc 和 http 服务
	check := new(api.AgentServiceCheck)
	if cr.serviceConf.IsRPC {
		check = &api.AgentServiceCheck{ // 健康检查
			Interval:                       interval.String(),              // 健康检查间隔
			GRPC:                           cr.serviceConf.HealthyCheckURL, // 健康检查地址
			DeregisterCriticalServiceAfter: deregister.String(),            // 注销时间，相当于过期时间
		}
	} else {
		check = &api.AgentServiceCheck{ // 健康检查
			Interval:                       interval.String(),              // 健康检查间隔
			HTTP:                           cr.serviceConf.HealthyCheckURL, // 健康检查地址
			DeregisterCriticalServiceAfter: deregister.String(),            // 注销时间，相当于过期时间
		}
	}
	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s", cr.serviceConf.Name, cr.serviceConf.IP), // 服务节点的名称
		Name:    cr.serviceConf.Name,                                          // 服务名称
		Tags:    cr.serviceConf.Tags,                                          // tag，可以为空
		Port:    cr.serviceConf.Port,                                          // 服务端口
		Address: cr.serviceConf.IP,                                            // 服务 IP
		Check:   check,
		Meta:    cr.serviceConf.Meta,
	}
	if err := agent.ServiceRegister(reg); err != nil {
		return errors.New(fmt.Sprintf("Service Register error : %v", err))
	}
	return nil
}

// DeregisterService sign off from consul by defined serviceID
func (cr *Resource) DeregisterService() error {
	agent := cr.consulClient.Agent()
	if agent == nil {
		return errors.New("fail get consul client agent()")
	}
	return cr.consulClient.Agent().ServiceDeregister(fmt.Sprintf("%s-%s", cr.serviceConf.Name, cr.serviceConf.IP))
}

// ClaimServices 声明依赖的服务
func (cr *Resource) ClaimServices(dependServices []string) bool {
	for _, service := range dependServices {
		InitResolver(cr.serverConf.Address, cr.serverConf.Port, cr.serverConf.Token, service)
	}
	return true
}

// @todo 设置成单例子？看 gRPC client 代码实现是否处理过了。
// DialService sign off from consul by defined serviceID
func (cr *Resource) Dial(service string) (*grpc.ClientConn, error) {
	//@todo 判断是否声明过依赖
	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s:%d/%s", "consul",
			cr.serverConf.Address,
			cr.serverConf.Port,
			service),
		//不能block => blockingPicker打开，在调用轮询时picker_wrapper => picker时若block则不进行robin操作直接返回失败
		//grpc.WithBlock(),
		grpc.WithInsecure(),
		//指定初始化round_robin => balancer (后续可以自行定制balancer和 register、resolver 同样的方式)
		grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
