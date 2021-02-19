package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wangzewang/esman/config"
	"github.com/wangzewang/esman/es"
	"github.com/wangzewang/esman/server"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	es.Init()
	server.Init()
}
