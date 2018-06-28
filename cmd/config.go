package cmd

import (
	"fmt"

	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/urfave/cli"

	"github.com/isymbo/pixpress/setting"
)

// Config show application configurations
var Config = cli.Command{
	Name:        "config",
	Usage:       "Show application configurations",
	Description: `Show application configuration on demand.`,
	Action:      showConfig,
}

func showConfig(c *cli.Context) error {

	fmt.Println("Web server port: ", setting.Server.HTTPPort)

	database, _ := prettyjson.Marshal(setting.Database)
	fmt.Printf("Database:\n%+v\n", string(database))

	ldap, _ := prettyjson.Marshal(setting.Ldap)
	fmt.Printf("LDAP:\n%+v\n", string(ldap))

	log, _ := prettyjson.Marshal(setting.Log)
	fmt.Printf("Log:\n%+v\n", string(log))

	other, _ := prettyjson.Marshal(setting.Other)
	fmt.Printf("Other:\n%+v\n", string(other))

	server, _ := prettyjson.Marshal(setting.Server)
	fmt.Printf("Server:\n%+v\n", string(server))

	return nil
}
