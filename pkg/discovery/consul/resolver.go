package consul

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const defaultSize = 50

var defaultSeed = time.Now().UnixNano()

// InitResolver returns
func InitResolver(host string, port int, token string, serviceName string) {
	log.Println("calling consul init")
	//resolver.Register(CacheBuilder())
	resolver.Register(NewBuilder(host, port, token, serviceName))
}

type consulBuilder struct {
	host        string
	port        int
	token       string
	serviceName string
}

type consulResolver struct {
	address              string
	token                string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	Ch                   chan int
	subsetSize           int
}

// NewBuilder ...
func NewBuilder(h string, p int, t string, sn string) resolver.Builder {
	return &consulBuilder{host: h, port: p, token: t, serviceName: sn}
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	subsetSize := viper.GetInt("subsetSize")
	if subsetSize == 0 {
		subsetSize = defaultSize
	}
	cr := &consulResolver{
		address:              fmt.Sprintf("%s:%d", cb.host, cb.port),
		token:                cb.token,
		name:                 cb.serviceName,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		Ch:                   make(chan int, 0),
		subsetSize:           subsetSize,
	}
	go cr.watcher()
	return cr, nil

}

func (cr *consulResolver) watcher() {
	log.Printf("calling [%s] consul watcher", cr.name)
	config := api.DefaultConfig()
	config.Address = cr.address
	config.Token = cr.token
	client, err := api.NewClient(config)
	if err != nil {
		log.Printf("error create consul client: %v", err)
		return
	}
	t := time.NewTicker(2000 * time.Millisecond)
	defer func() {
		log.Println("watcher defer")
	}()
	i := 0
	for {
		select {
		case <-t.C:
			//fmt.Println("定时")
		case <-cr.Ch:
			//fmt.Println("ch call")
		}
		//api添加了 lastIndex   consul api中并不兼容附带lastIndex的查询
		services, _, err := client.Health().Service(cr.name, "", true, &api.QueryOptions{})
		if err != nil {
			log.Printf("error retrieving instances from Consul: %v", err)
			if i%5 == 0 {
				i = 0
				// _ = alert.SendDingMsg(fmt.Sprintf("consul health err: %v", err), false)
			}
			i++
		} else {
			i = 0
		}
		newAddrs := make([]resolver.Address, 0)
		for _, service := range services {
			addr := net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))
			newAddrs = append(newAddrs, resolver.Address{
				Addr: addr,
				//type：不能是grpclib，grpclib在处理链接时会删除最后一个链接地址，不用设置即可 详见=> balancer_conn_wrappers => updateClientConnState
				ServerName: service.Service.Service,
			})
		}
		// 如何地址长度大于subsetSize 取地址集合的子集
		if len(newAddrs) > cr.subsetSize {
			rand.Seed(defaultSeed)
			rand.Shuffle(len(newAddrs), func(i, j int) { newAddrs[i], newAddrs[j] = newAddrs[j], newAddrs[i] })
			newAddrs = newAddrs[:cr.subsetSize]
		}
		cr.cc.UpdateState(resolver.State{Addresses: newAddrs})
	}
}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	cr.Ch <- 1
}

func (cr *consulResolver) Close() {
}
