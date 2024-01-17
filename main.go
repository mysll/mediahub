package main

import (
	"github.com/mysll/toolkit"
	"mediahub/internal/conf"
	"mediahub/server"
)

func main() {
	server.Start(conf.LoadOption(conf.WithDataPath("./data")))
	toolkit.WaitForQuit()
	server.Close()
}
