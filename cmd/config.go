package cmd

import (
	"fmt"

	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/isymbo/pixpress/setting"
	"github.com/urfave/cli"
)

// Config show application configurations
var Config = cli.Command{
	Name:        "config",
	Usage:       "Show application configurations",
	Description: `Show application configuration on demand.`,
	Action:      showConfig,
}

func showConfig(c *cli.Context) error {

	server, _ := prettyjson.Marshal(setting.Server)
	fmt.Printf("Server:\n%+v\n", string(server))

	session, _ := prettyjson.Marshal(setting.Session)
	fmt.Printf("Session:\n%+v\n", string(session))

	security, _ := prettyjson.Marshal(setting.Security)
	fmt.Printf("Security:\n%+v\n", string(security))

	ldap, _ := prettyjson.Marshal(setting.Ldap)
	fmt.Printf("LDAP:\n%+v\n", string(ldap))

	database, _ := prettyjson.Marshal(setting.Database)
	fmt.Printf("Database:\n%+v\n", string(database))

	service, _ := prettyjson.Marshal(setting.Service)
	fmt.Printf("Service:\n%+v\n", string(service))

	log, _ := prettyjson.Marshal(setting.Log)
	fmt.Printf("Log:\n%+v\n", string(log))

	logconsole, _ := prettyjson.Marshal(setting.LogConsole)
	fmt.Printf("LogConsole:\n%+v\n", string(logconsole))

	logfile, _ := prettyjson.Marshal(setting.LogFile)
	fmt.Printf("LogFile:\n%+v\n", string(logfile))

	logslack, _ := prettyjson.Marshal(setting.LogSlack)
	fmt.Printf("LogSlack:\n%+v\n", string(logslack))

	logdiscord, _ := prettyjson.Marshal(setting.LogDiscord)
	fmt.Printf("LogDiscord:\n%+v\n", string(logdiscord))

	logxorm, _ := prettyjson.Marshal(setting.LogXorm)
	fmt.Printf("LogXorm:\n%+v\n", string(logxorm))

	other, _ := prettyjson.Marshal(setting.Other)
	fmt.Printf("Other:\n%+v\n", string(other))

	return nil
}
