package main

import (
	"flag"
	"os"
	"strings"

	// 导入 kratos 的 dtm 驱动
	_ "github.com/dtm-labs/driver-kratos"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc/resolver"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"
	etcdAPI "go.etcd.io/etcd/client/v3"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// main.go 创建 kratos 应用生命周期管理
func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	var etcdServer = "10.211.55.29:2379"

	//采用Kratos的方法创建registry, https://go-kratos.dev/docs/component/registry/
	client, err := etcdAPI.New(etcdAPI.Config{
		Endpoints: strings.Split(etcdServer, ","),
	})
	if err != nil {
		panic(err)
	}
	registry := etcd.New(client)
	defer client.Close()

	//注册全局的resolver,  现在business server可以使用 discovery:///dtmservice 来访问dtm
	resolver.Register(discovery.NewBuilder(registry, discovery.WithInsecure(true)))
	// with registrar
	//kratos.Registrar(registry)

	//var dtmServer = "discovery:///dtmservice"

	app, cleanup, err := wireApp(flagconf, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
