package wrr

import (
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

// Name is the name of round_robin balancer.
const Name = "round_robin"

var (
	logger = grpclog.Component("roundrobin")
)

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &rrPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct{}

func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("roundrobinPicker: newPicker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	rr := &rrPicker{colorConnMap: make(map[string][]balancer.SubConn)}

	var scs []balancer.SubConn
	for sc, addr := range info.ReadySCs {
		if color, ok := addr.Address.Attributes.Value("color").(string); ok {
			rr.colorConnMap[color] = append(rr.colorConnMap[color], sc)
		}

		scs = append(scs, sc)
	}
	rr.subConns = scs
	return rr
}

type rrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns     []balancer.SubConn
	colorConnMap map[string][]balancer.SubConn
	mu           sync.Mutex
	next         int
}

func (p *rrPicker) Pick(pi balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if color, ok := pi.Ctx.Value("color").(string); ok {
		if scs, ok := p.colorConnMap[color]; ok {
			sc := scs[p.next]
			p.next = (p.next + 1) % len(p.subConns)
			return balancer.PickResult{SubConn: sc}, nil
		}
	}
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	return balancer.PickResult{SubConn: sc}, nil
}
