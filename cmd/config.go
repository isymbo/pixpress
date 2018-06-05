package cmd

import (
	"fmt"

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

	fmt.Println("Web server port: ", setting.WebPort)

	fmt.Println("Database host: ", setting.Database.Host)
	fmt.Println("Database port: ", setting.Database.Port)

	fmt.Println("Ldap host: ", setting.Ldap.Host)
	fmt.Println("Ldap port: ", setting.Ldap.Port)
	fmt.Println("Ldap base: ", setting.Ldap.Base)
	fmt.Println("Ldap bind dn: ", setting.Ldap.BindDn)
	fmt.Println("Ldap password: ", setting.Ldap.Password)

	fmt.Println("Other show footer template load time: ", setting.Other.ShowFooterTemplateLoadTime)

	return nil
}
