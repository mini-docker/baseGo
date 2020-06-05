package main

import (
	"fmt"
	"time"

	"fecho/registry"
)

func main() {

	service := &registry.Service{
		Name:      "aaaa",
		Version:   "1",
		Metadata:  nil,
		Endpoints: nil,
		Nodes: []*registry.Node{
			{
				Id:       "1",
				Address:  "localhost",
				Port:     7071,
				Metadata: nil,
			},
		},
	}

	r := registry.NewEtcdV3Registry(
		registry.Timeout(5*time.Second),
		registry.Addrs("127.0.0.1:2379"),
	)

	err := r.Register(service, registry.RegisterTTL(5*time.Second))

	registry.DefaultRegistry = r

	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(1 * time.Second)
		services, err := r.ListServices()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, v := range services {
			fmt.Println("aaaa --> ", v.Name)
		}
	}
}
