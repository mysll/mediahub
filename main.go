package main

import (
	"flag"
	"github.com/mysll/toolkit"
	"mediahub/internal/conf"
	"mediahub/server"
	"runtime/debug"
)

var (
	dev = flag.Bool("dev", false, "dev mode")
)

func main() {
	flag.Parse()
	debug.SetTraceback("single")
	server.Start(conf.LoadOption(conf.WithDataPath("./data")))
	toolkit.WaitForQuit()
	server.Close()
}
