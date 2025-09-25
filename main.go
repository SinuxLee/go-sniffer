package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/sinuxlee/go-sniffer/core"
	"github.com/urfave/cli"
)

var (
	version = "No Build Info"
	app     *cli.App
)

func init() {
	app = cli.NewApp()
	app.Name = "go-sniffer"
	app.Usage = "A tool to capture and analyze network packets"
	app.Version = version

	cli.VersionPrinter = func(c *cli.Context) {
		_, _ = fmt.Fprintf(c.App.Writer, "%v build info: %v\n", c.App.Name, strings.ReplaceAll(c.App.Version, "_", "\n"))
	}
}

func main() {
	app.Commands = []cli.Command{
		{
			Name:  "device",
			Usage: "Show all devices",
			Action: func(c *cli.Context) error {
				ifaces, err := net.Interfaces()
				if err != nil {
					return err
				}

				for _, iface := range ifaces {
					addrs, _ := iface.Addrs()
					for _, a := range addrs {
						if ipnet, ok := a.(*net.IPNet); ok {
							if ip4 := ipnet.IP.To4(); ip4 != nil {
								fmt.Println("[device] : " + iface.Name + " : " + iface.HardwareAddr.String() + "  " + ip4.String())
							}
						}
					}
				}

				return nil
			},
		}, {
			Name:  "list",
			Usage: "Show all plugins",
			Action: func(c *cli.Context) error {
				plug, err := core.NewPlug()
				if err != nil {
					return err
				}
				plug.PrintList()

				return nil
			},
		}, {
			Name:  "env",
			Usage: "Show plugin's dir",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "plugin",
			Value: "http",
			Usage: "Specify the plugin to use",
		},
		cli.StringFlag{
			Name:  "device",
			Value: "eth0",
			Usage: "Specify the device to use",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "80",
			Usage: "Specify the port to use",
		},
	}

	app.Action = func(c *cli.Context) error {
		plugin := c.String("plugin")
		device := c.String("device")
		port := c.String("port")

		plug, err := core.NewPlug()
		if err != nil {
			return err
		}

		err = plug.SetOption(plugin, []string{"port", port})
		if err != nil {
			return err
		}

		core.NewDispatch(plug, device).Capture()
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
