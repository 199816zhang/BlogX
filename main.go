package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"fmt"
)

func main() {
	flags.Parse()
	fmt.Println(flags.FlagOptions)
	core.ReadConf()
}
