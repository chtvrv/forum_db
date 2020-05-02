package main

import (
	"flag"
	"log"

	serv "github.com/chtvrv/forum_db/app/server"
	"github.com/spf13/viper"
)

var opts struct {
	configPath string
}

func main() {
	flag.StringVar(&opts.configPath, "c", "", "path to configuration file")
	flag.StringVar(&opts.configPath, "config", "", "path to configuration file")

	flag.Parse()
	viper.SetConfigFile(opts.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	server := new(serv.Server)

	server.Run()
}
