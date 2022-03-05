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

func NewDNSDiscovery(log *log.Helper) registry.Discovery {
	return &DNSDiscovery{
		services: map[string]registry.Watcher{},
		log:      log,
	}
}

// DNSDiscovery
// 服务名使用: "discovery:///cdp-usergroup-svc.cdp:9000"
type DNSDiscovery struct {
	services map[string]registry.Watcher
	log      *log.Helper
}

func (s *DNSDiscovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return s.services[serviceName].(*DNSWatcher).latest, nil
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
		name:      serviceName,
		port:      port,
		instances: map[string]*registry.ServiceInstance{},
		changed:   make(chan []*registry.ServiceInstance, 1),
		log:       s.log,
	}
	go dw.watch()
	s.services[serviceName] = dw
	return s.services[serviceName], nil
}

type DNSWatcher struct {
	port      string
	name      string
	instances map[string]*registry.ServiceInstance
	changed   chan []*registry.ServiceInstance
	latest    []*registry.ServiceInstance
	log       *log.Helper
}

func (m *DNSWatcher) Next() ([]*registry.ServiceInstance, error) {
	return <-m.changed, nil
}
func (m *DNSWatcher) watch() {
	for {
		m.watch1()
		time.Sleep(time.Second * 5)
	}
}
func (m *DNSWatcher) watch1() {
	a, err := net.LookupHost(m.name)
	if err != nil {
		m.log.Warnf("resolve host failed, hostname=%s,message=%s", m.name, err.Error())
	}
	_, srvs, err := net.LookupSRV("grpclb", "tcp", m.name)
	if err != nil {
		m.log.Warnf("resolve grpclb failed, hostname=%s,message=%s", m.name, err.Error())
		srvs = make([]*net.SRV, 0)
	}

	if len(srvs) == 0 {
		for _, v := range a {
			srvs = append(srvs, &net.SRV{
				Target: v,
				Port:   0,
			})
		}
	}

	newIns := make(map[string]struct{})
	hasChange := false

	for _, v := range srvs {
		a, err := net.LookupHost(v.Target)
		if err != nil {
			m.log.Warnf("resolve host failed, hostname=%s,message=%s", v.Target, err.Error())
			continue
		}
		for _, addr := range a {
			id := fmt.Sprintf("%s:%s", addr, m.port)
			newIns[id] = struct{}{}
			if _, ok := m.instances[id]; ok {
				continue
			}

			address := fmt.Sprintf("grpc://%s:%s?isSecure=false", addr, m.port)
			m.instances[id] = &registry.ServiceInstance{
				ID:        id,
				Name:      m.name,
				Version:   "v1",
				Metadata:  map[string]string{},
				Endpoints: []string{address},
			}
			hasChange = true
			break
		}
	}
	list := make([]*registry.ServiceInstance, 0)
	for n, ins := range m.instances {
		if _, ok := newIns[n]; !ok {
			delete(m.instances, n)
			hasChange = true
			continue
		}
		list = append(list, ins)
	}
	m.latest = list
	if hasChange {
		m.changed <- list
	}
}

func (m *DNSWatcher) Stop() error {
	return nil
}
