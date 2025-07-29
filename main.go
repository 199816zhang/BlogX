package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()
	core.InitIPDB()
	fmt.Println(core.GetIpAddr("113.90.151.170"))
}
