package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/isymbo/pixpress/cmd"
	"github.com/isymbo/pixpress/setting"
)

func main() {

	app := cli.NewApp()
	app.Name = setting.APP_NAME
	app.Usage = setting.APP_USAGE
	app.Version = setting.APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Config,
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "specify `FILE` for configuration of " + setting.APP_NAME,
			Value: setting.CfgPath,
		},
	}

	app.Before = setting.LoadConfig

	app.Run(os.Args)
}
