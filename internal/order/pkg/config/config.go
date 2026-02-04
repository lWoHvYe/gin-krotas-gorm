package config

import (
	"helloworld-go/internal/conf"
	"log"
	"time"

	cfg "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/config"
	etcdAPI "go.etcd.io/etcd/client/v3"
)

func NewConfig(flagconf string) (*conf.Bootstrap, error) {
	// create an etcd client
	client, err := etcdAPI.New(etcdAPI.Config{
		Endpoints:   []string{"10.211.55.29:2379"},
		DialTimeout: time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	// configure the source, "path" is required
	source, err := cfg.New(client, cfg.WithPath(flagconf), cfg.WithPrefix(true))
	if err != nil {
		log.Fatalln(err)
	}

	// create a config instance with source
	c := config.New(config.WithSource(source))
	defer c.Close()

	// load sources before get
	if err := c.Load(); err != nil {
		log.Fatalln(err)
	}

	// acquire config value
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}

	/*if err := c.Watch("service.name", func(key string, value config.Value) {
		fmt.Printf("config changed: %s = %v\n", key, value)
		// Write your callback logic here
	}); err != nil {
		log.Error(err)
	}*/

	return &bc, nil
}

func NewDiscovery(conf *conf.Bootstrap) (*etcd.Registry, error) {
	// new etcd client
	client, err := etcdAPI.New(etcdAPI.Config{
		Endpoints: conf.Registry.Etcd.Endpoints,
	})
	if err != nil {
		panic(err)
	}
	// new dis with etcd client
	dis := etcd.New(client)
	//defer client.Close()

	return dis, err
}
