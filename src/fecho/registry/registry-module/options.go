package registry_module

import "strings"

type Options struct {
	registryAddr []string
	containerId  string // 容器id .
	addr         string
	port         int
	ttl          int
	sites        string
	version      string
	tcpAddr      string
	wsAddr       string
	httpAddr     string
	name         string
	group        string
}

type Option func(*Options)

func WithName(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

func WithTcpAddr(tcpAddr string) Option {
	return func(o *Options) {
		o.tcpAddr = tcpAddr
	}
}

func WithWsAddr(wsAddr string) Option {
	return func(o *Options) {
		o.wsAddr = wsAddr
	}
}

func WithHttpAddr(httpAddr string) Option {
	return func(o *Options) {
		o.httpAddr = httpAddr
	}
}

func WithContainerId(containerId string) Option {
	return func(o *Options) {
		o.containerId = containerId
	}
}

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.addr = addr
	}
}

func WithPort(port int) Option {
	return func(o *Options) {
		o.port = port
	}
}

func WithSites(sites ...string) Option {
	return func(o *Options) {
		o.sites = strings.Join(sites, ",")
	}
}

func WithTTL(ttl int) Option {
	return func(o *Options) {
		o.ttl = ttl
	}
}

func WithVersion(version string) Option {
	return func(o *Options) {
		o.version = version
	}
}

func WithRegistryAddr(addr ...string) Option {
	return func(o *Options) {
		o.registryAddr = addr
	}
}
