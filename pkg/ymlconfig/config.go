package ymlconfig

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Viper *viper.Viper
}

var (
	configVar      string
	isRemoteConfig bool

	GlobalSource = pflag.String("global.source", "default(flag)", "identify the source of configuration")
	GlobalUnset  = pflag.String("global.unset", "", "this parameter do not appear in config file")
)

func init() {
	pflag.StringVar(&configVar, "config", "", "Config file path")
	pflag.BoolVar(&isRemoteConfig, "isRemoteConfig", false, "Whether to choose remote config")
}

func (c *Config) SetRemoteConfig(u *url.URL) {
	var provider string
	var endpoint string
	var path string

	schemes := strings.SplitN(u.Scheme, "+", 2)
	if len(schemes) < 1 {
		klog.Fatalf("invalid config scheme '%s'", u.Scheme)
	}

	provider = schemes[0]

	switch provider {
	case "etcd":
		if len(schemes) < 2 {
			klog.Fatalf("invalid config scheme '%s'", u.Scheme)
		}
		protocol := schemes[1]
		endpoint = fmt.Sprintf("%s://%s", protocol, u.Host)
		path = u.Path // u.Path = /path/to/key.yaml
	case "consul":
		endpoint = u.Host
		path = u.Path[1:] // u.Path = /key.json
	default:
		klog.Fatalf("unsupported provider '%s'", provider)
	}

	// 文件后缀
	ext := filepath.Ext(path)
	if ext == "" {
		klog.Fatalf("using remote config, without specifying file extension")
	}
	configType := ext[1:]

	klog.Infof("Using Remote Config Provider: '%s', Endpoint: '%s', Path: '%s', ConfigType: '%s'", provider, endpoint, path, configType)
	if err := c.Viper.AddRemoteProvider(provider, endpoint, path); err != nil {
		klog.Fatalf("error adding remote provider %s", err)
	}

	c.Viper.SetConfigType(configType)
}

func (c *Config) SetDefaultValue() {
	c.Viper.SetDefault("global.unset", "default(viper)")
}

func (c *Config) WatchRemoteConf() {
	for {
		time.Sleep(time.Second * 5) // delay after each request

		// currently, only tested with etcd support
		err := c.Viper.WatchRemoteConfig()
		if err != nil {
			klog.Errorf("unable to read remote config: %v", err)
			continue
		}

		// unmarshal new config into our runtime config struct. you can also use channel
		// to implement a signal to notify the system of the changes
		// runtime_viper.Unmarshal(&runtime_conf)
		klog.Info("Watching Remote Config")
		klog.Infof("Global.Source: '%s'", c.Viper.GetString("Global.Source"))
		klog.Infof("Global.ChangeMe: '%s'", c.Viper.GetString("Global.ChangeMe"))
	}
}

// ZapLogConfig 读取Log的配置文件，并返回
func (c *Config) ZapLogConfig() []byte {
	log := c.Viper.Sub("Log")
	logConfig, err := json.Marshal(log.AllSettings())
	if err != nil {
		klog.Fatalf("error marshalling log config %s", err)
	}
	return logConfig
}

func ConfigInit(envPrefix string, cfgName string) Config {
	pflag.Parse()

	v := viper.New()
	config := Config{Viper: v}
	Viper := config.Viper

	Viper.BindPFlags(pflag.CommandLine)
	config.SetDefaultValue()

	// read from env
	Viper.AutomaticEnv()
	// so that client.foo maps to MYAPP_CLIENT_FOO
	Viper.SetEnvPrefix(envPrefix)
	Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configVar != "" {
		/*
			如果设置了--config参数，尝试从这里解析
			它可能是一个Remote Config，来自etcd或consul
			也可能是一个本地文件
		*/
		u, err := url.Parse(configVar)
		if err != nil {
			klog.Fatalf("error parsing: '%s'", configVar)
		}

		if u.Scheme != "" {
			// 看起来是个remote config
			config.SetRemoteConfig(u)
			isRemoteConfig = true
		} else {
			Viper.SetConfigFile(configVar)
		}
	} else {
		/*
			尝试搜索若干默认路径，先后顺序如下:
			- /etc/tiktok/config/userConfig.<ext>
			- ~/.tiktok/userConfig.<ext>
			- ./userConfig.<ext>

			其中<ext> 是 viper所支持的文件类型，如yml，json等
		*/

		Viper.SetConfigName(cfgName) // name of config file (without extension)
		Viper.AddConfigPath("/etc/go_project/byte_dance_5th/config")
		Viper.AddConfigPath("$HOME/.byte_dance_5th/")
		Viper.AddConfigPath("./config")
		Viper.AddConfigPath("../../config")
		Viper.AddConfigPath("../../../config")
	}

	if isRemoteConfig {
		if err := Viper.ReadRemoteConfig(); err != nil {
			klog.Fatalf("error reading config: %s", err)
		}
		klog.Infof("Using Remote Config: '%s'", configVar)

		Viper.WatchRemoteConfig()
		// 另启动一个协程来监测远程配置文件
		go config.WatchRemoteConf()

	} else {
		if err := Viper.ReadInConfig(); err != nil {
			klog.Fatalf("error reading config: %s", err)
		}
		klog.Infof("Using configuration file '%s'", Viper.ConfigFileUsed())

		Viper.WatchConfig()
		Viper.OnConfigChange(func(e fsnotify.Event) {
			klog.Info("Config file changed:", e.Name)
			klog.Infof("Global.Source: '%s'", Viper.GetString("Global.Source"))
			klog.Infof("Global.ChangeMe: '%s'", Viper.GetString("Global.ChangeMe"))
		})

	}

	return config
}
