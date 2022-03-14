package discovery

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

type Callback func(instances []*registry.ServiceInstance)

func NewDNSDiscovery(log *log.Helper, callback Callback) registry.Discovery {
	return &DNSDiscovery{
		services: map[string]registry.Watcher{},
		log:      log,
		callback: callback,
	}
}

// DNSDiscovery
// 服务名使用: "discovery:///cdp-usergroup-svc.cdp:9000"
type DNSDiscovery struct {
	services map[string]registry.Watcher
	log      *log.Helper
	callback Callback
}

func (s *DNSDiscovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return s.services[serviceName].(*DNSWatcher).toArray(), nil
}

func (s *DNSDiscovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	tmp := strings.Split(serviceName, ":")
	serviceName = tmp[0]
	port := "9000"
	if len(tmp) > 1 {
		port = tmp[1]
	}
	if s.services == nil {
		s.services = make(map[string]registry.Watcher)
	}
	if v, ok := s.services[serviceName]; ok {
		return v, nil
	}
	dw := &DNSWatcher{
		name:     serviceName,
		port:     port,
		changed:  make(chan struct{}, 1),
		log:      s.log,
		callback: s.callback,
	}
	go dw.watch()
	s.services[serviceName] = dw
	return s.services[serviceName], nil
}

type DNSWatcher struct {
	port     string
	name     string
	changed  chan struct{}
	latest   map[string]*registry.ServiceInstance
	log      *log.Helper
	callback Callback
}

func (m *DNSWatcher) Next() ([]*registry.ServiceInstance, error) {
	if m.latest == nil {
		time.Sleep(time.Second)
	} else {
		<-m.changed
	}
	if m.callback != nil {
		go m.callback(m.toArray())
	}
	return m.toArray(), nil
}

func (m *DNSWatcher) toArray() []*registry.ServiceInstance {
	ins := make([]*registry.ServiceInstance, len(m.latest))
	idx := 0
	for _, v := range m.latest {
		ins[idx] = v
		idx++
	}
	return ins
}

func (m *DNSWatcher) watch() {
	for {
		// 每5重新解析dns
		m.lookup()
		time.Sleep(time.Second * 5)
	}
}
func (m *DNSWatcher) lookup() {
	_, srvs, err := net.LookupSRV("grpclb", "tcp", m.name)
	if err != nil {
		// m.log.Debugf("resolve grpclb failed, hostname=%s,message=%s", m.name, err.Error())
		srvs = make([]*net.SRV, 0)
	}

	// 没有找到grpclb SRV记录再解析A记录
	if len(srvs) == 0 {
		addrs, err := net.LookupHost(m.name)
		if err != nil {
			m.log.Debugf("resolve host failed, hostname=%s,message=%s", m.name, err.Error())
		}
		for _, v := range addrs {
			if v == "::1" {
				continue
			}
			srvs = append(srvs, &net.SRV{
				Target: v,
				Port:   0,
			})
		}
	}

	hasChanged := false
	latest := make(map[string]*registry.ServiceInstance, 0)
	for _, v := range srvs {
		addrs, err := net.LookupHost(v.Target)
		if err != nil {
			m.log.Debugf("resolve host failed, hostname=%s,message=%s", v.Target, err.Error())
			continue
		}
		for _, addr := range addrs {
			id := fmt.Sprintf("%s:%s", addr, m.port)
			m.log.Debugf("resolve host found, serviceName=%s, instance=%s", m.name, id)
			if v, ok := m.latest[id]; ok {
				latest[id] = v
			} else {
				hasChanged = true
				latest[id] = &registry.ServiceInstance{
					ID:        id,
					Name:      m.name,
					Version:   "v1",
					Metadata:  map[string]string{},
					Endpoints: []string{fmt.Sprintf("grpc://%s:%s?isSecure=false", addr, m.port)},
				}
			}
		}
	}
	if !hasChanged {
		for n, _ := range m.latest {
			if _, ok := latest[n]; ok {
				continue
			}
			hasChanged = true
			break
		}
	}
	if hasChanged {
		m.latest = latest
		m.changed <- struct{}{}
	}
}

func (m *DNSWatcher) Stop() error {
	return nil
}
