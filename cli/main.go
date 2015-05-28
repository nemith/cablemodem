package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/nemith/cablemodem"
)

const (
	defaultHost = "192.168.100.1"
)

var cm cablemodem.Modem

func cmdChannel(c *cli.Context) {
	signal := cm.SignalData()

	fmt.Println("Upstream Channels")
	for _, c := range signal.Upstream {
		fmt.Printf("  Channel: %d\n", c.ID)
		fmt.Printf("    Frequency          : %d Mhz\n", c.Freq/1000/1000)
		fmt.Printf("    Ranging Service ID : %d\n", c.RangingServiceID)
		fmt.Printf("    Modulation         : %s\n", c.Modulation)
		fmt.Printf("    Ranging Status     : %s\n", c.RangingStatus)
		fmt.Printf("    Power Level        : %d dBmV\n", c.Power)
	}
	fmt.Println("")

	fmt.Println("Downstream Channels")
	for _, c := range signal.Downstream {
		fmt.Printf("  Channel: %d\n", c.ID)
		fmt.Printf("    Frequency      : %d Mhz\n", c.Freq/1000/1000)
		fmt.Printf("    Modulation     : %s\n", c.Modulation)
		fmt.Printf("    SNR            : %d dB\n", c.SNR)
		fmt.Printf("    Power Level    : %d dBmV\n", c.Power)
		fmt.Printf("    Codeword Stats : %d/%d/%d (Total/Correctable/Uncorrectable)\n",
			c.UnerroredCodewords, c.CorrectableCodewords, c.UncorrectableCodewords)
	}
	fmt.Println("")

}

func defaultString(c *cli.Context, flag, defValue string) string {
	val := c.String(flag)
	if val == "" {
		return defValue
	}
	return val
}

func main() {
	app := cli.NewApp()
	app.Name = "cm"
	app.Usage = "interface cable modem data"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Usage: fmt.Sprintf("Hostname/IP for cablemodem (defaults to %s)", defaultHost),
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "channels",
			Aliases: []string{"c"},
			Usage:   "channel stats",
			Action:  cmdChannel,
		},
	}
	app.Before = func(c *cli.Context) error {
		cm = cablemodem.NewSurfboardCM(defaultString(c, "host", defaultHost))
		return nil
	}
	app.Run(os.Args)
}
