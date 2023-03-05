package main

import (
	"byte_dance_5th/cmd/api/rpc"
	"byte_dance_5th/pkg/ymlconfig"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	h2config "github.com/hertz-contrib/http2/config"
	"github.com/hertz-contrib/http2/factory"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/hertz-contrib/gzip"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/registry/etcd"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

type HertzConfig struct {
	UseNetpoll bool  `json:"UseNetpoll" yaml:"UseNetpoll"`
	Http2      Http2 `json:"Http2" yaml:"Http2"`
	Tls        Tls   `json:"Tls" yaml:"Tls"`
}

type Http2 struct {
	Enable           bool     `json:"Enable" yaml:"Enable"`
	DisableKeepalive bool     `json:"DisableKeepalive" yaml:"DisableKeepalive"`
	ReadTimeout      Duration `json:"ReadTimeout" yaml:"ReadTimeout"`
}

type Tls struct {
	Enable bool `json:"Enable" yaml:"Enable"`
	Cfg    tls.Config
	Cert   string `json:"CertFile" yaml:"CertFile"`
	Key    string `json:"KeyFile" yaml:"KeyFile"`
	ALPN   bool   `json:"ALPN" yaml:"ALPN"`
}

type Duration struct {
	time.Duration
}

var (
	Config      = ymlconfig.ConfigInit("TIKTOK_API", "apiConfig")
	ServiceName = Config.Viper.GetString("Server.Name")
	ServiceAddr = fmt.Sprintf("%s:%d", Config.Viper.GetString("Server.Address"), Config.Viper.GetInt("Server.Port"))
	EtcdAddress = fmt.Sprintf("%s:%d", Config.Viper.GetString("Etcd.Address"), Config.Viper.GetInt("Etcd.Port"))
	hertzConfig HertzConfig
)

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// Init 初始化 API 配置
func Init() {
	rpc.InitRpc(&Config)
}

func InitHertzCfg() {
	hertzV, err := json.Marshal(Config.Viper.Sub("Hertz").AllSettings())
	if err != nil {
		hlog.Fatalf("Error marshalling Hertz config %s", err)
	}
	if err := json.Unmarshal(hertzV, &hertzConfig); err != nil {
		hlog.Fatalf("Error unmarshalling Hertz config %s", err)
	}
}

// InitHertz 初始化 Hertz
func InitHertz() *server.Hertz {
	InitHertzCfg()

	opts := []config.Option{server.WithHostPorts(ServiceAddr)}

	// 服务注册
	if Config.Viper.GetBool("Etcd.Enable") {
		r, err := etcd.NewEtcdRegistry([]string{EtcdAddress})
		if err != nil {
			hlog.Fatal(err)
		}
		opts = append(opts, server.WithRegistry(r, &registry.Info{
			ServiceName: ServiceName,
			Addr:        utils.NewNetAddr("tcp", ServiceAddr),
			Weight:      10,
			Tags:        nil,
		}))
	}

	// 链路追踪
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(ServiceName),
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	tracer, tracerCfg := hertztracing.NewServerTracer()
	opts = append(opts, tracer)

	// 网络库
	hertzNet := standard.NewTransporter
	if hertzConfig.UseNetpoll {
		hertzNet = netpoll.NewTransporter
	}
	opts = append(opts, server.WithTransport(hertzNet))

	// TLS & Http2
	tlsEnable := hertzConfig.Tls.Enable
	h2Enable := hertzConfig.Http2.Enable
	hertzConfig.Tls.Cfg = tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	if tlsEnable {
		cert, err := tls.LoadX509KeyPair(hertzConfig.Tls.Cert, hertzConfig.Tls.Key)
		if err != nil {
			hlog.Error(err)
		}
		hertzConfig.Tls.Cfg.Certificates = append(hertzConfig.Tls.Cfg.Certificates, cert)
		opts = append(opts, server.WithTLS(&hertzConfig.Tls.Cfg))

		if alpn := hertzConfig.Tls.ALPN; alpn {
			opts = append(opts, server.WithALPN(alpn))
		}
	} else if h2Enable {
		opts = append(opts, server.WithH2C(h2Enable))
	}

	// Hertz
	h := server.Default(opts...)
	h.Use(gzip.Gzip(gzip.DefaultCompression),
		hertztracing.ServerMiddleware(tracerCfg))

	// Protocol
	if h2Enable {
		h.AddProtocol("h2", factory.NewServerFactory(
			h2config.WithReadTimeout(hertzConfig.Http2.ReadTimeout.Duration),
			h2config.WithDisableKeepAlive(hertzConfig.Http2.DisableKeepalive)))
		if tlsEnable {
			hertzConfig.Tls.Cfg.NextProtos = append(hertzConfig.Tls.Cfg.NextProtos, "h2")
		}
	}

	return h
}
