package registry_module

import (
	"fmt"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	Start(
		WithRegistryAddr("127.0.0.1:2379"),
		WithAddr("localhost"),
		WithContainerId("localhost"),
		WithPort(7070),
		WithSites("aaaa"),
		WithTTL(5),
		WithVersion("testing"),
	)

	// wait for registry
	time.Sleep(1 * time.Second)

	services := GetServices("aaa", "a")

	if len(services) != 1 {
		t.Errorf("services list failed ")
	}

	s := services[0]
	if s.SiteId != "aaa" && s.SiteIndexId != "a" && s.HostPort != "localhost:7070" {
		t.Errorf("service info error %v", s)
	}

	all := ListServices()
	if len(all) != 1 {
		t.Errorf("services list failed ")
	}

	ss, ok := all["aaaa"]
	if !ok {
		t.Errorf("failed to list services ")
	}

	s = ss[0]
	if s.SiteId != "aaa" && s.SiteIndexId != "a" && s.HostPort != "localhost:7070" {
		t.Errorf("service info error %v", s)
	}
}

func TestAll2(t *testing.T) {
	Start(
		WithRegistryAddr("127.0.0.1:2379"),
		WithAddr("localhost"),
		WithContainerId("localhost"),
		WithPort(7070),
		WithSites("aaaa"),
		WithTTL(5),
		WithVersion("testing"),
	)

	go Start(
		WithRegistryAddr("127.0.0.1:2379"),
		WithAddr("localhost"),
		WithContainerId("localhost2"),
		WithPort(7071),
		WithSites("aaab"),
		WithTTL(5),
		WithVersion("testing"),
	)

	// wait for registry
	time.Sleep(1 * time.Second)

	services := GetServices("aaa", "a")

	for _, v := range services {
		fmt.Println(v.HostPort)
	}
	if len(services) != 1 {
		t.Errorf("services list failed ")
	}

	all := ListServices()
	if len(all) != 2 {
		t.Errorf("services list failed ")
	}
}
