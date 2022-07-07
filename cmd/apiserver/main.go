package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"

	"premier_api/internal/app/apiserver"

	pm2io "github.com/keymetrics/pm2-io-apm-go"
	"github.com/keymetrics/pm2-io-apm-go/structures"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
	pm2 := pm2io.Pm2Io{
		Config: &structures.Config{
			PublicKey:  "myPublic",   // define the public key given in the dashboard
			PrivateKey: "myPrivate",  // define the private key given in the dashboard
			Name:       "MainServer", // define an application name
		},
	}
	pm2.Start()

}
