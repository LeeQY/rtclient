package main

import (
	"errors"
	"os"

	"github.com/drkaka/rtclient/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rtclient"
	app.Usage = "command line application for RescueTime"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "key, k",
			Usage: "Specify the API key.",
		},
	}

	app.Before = func(c *cli.Context) error {
		k := c.GlobalString("key")
		if len(k) == 0 {
			return errors.New("please specify you API key to continue")
		}
		return nil
	}

	app.Commands = []cli.Command{
		cmd.NewListCMD(),
	}

	app.Run(os.Args)
}
