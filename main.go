package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"github.com/sirupsen/logrus"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	logrus.Errorf("yyyyy")
	logrus.Debug("zzzz")
}
