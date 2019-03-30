package svrcfg

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

const (
	configFile = "config.ini"
)

type config struct {
	AppMode string
	// [server]
	PidFile      string
	UseTLS       bool
	CertPath     string
	CertKeyPath  string
	ServerAddr   string
	ServerPort   int
	HostOverride string
}

// Config is the exported configuration
var Config = &config{}

// sets up the exported configuration for any packages that import this config pkg
func init() {
	cfg, err := ini.Load(configFile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// [DEFAULT]
	Config.AppMode = cfg.Section("").Key("app_mode").String()

	// [server]
	serverSection := cfg.Section("server")
	Config.PidFile = serverSection.Key("pidfile").String()
	Config.UseTLS = serverSection.Key("use_tls").MustBool(false)
	Config.CertPath = serverSection.Key("cert_path").String()
	Config.CertKeyPath = serverSection.Key("cert_key_path").String()
	Config.ServerAddr = serverSection.Key("server_addr").String()
	Config.ServerPort = serverSection.Key("server_port").MustInt(10000)
	Config.HostOverride = serverSection.Key("host_override").String()

	fmt.Printf("%v\n", Config)
}
