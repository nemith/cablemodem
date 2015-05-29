package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/nemith/cablemodem"
)

const (
	defaultHost = "192.168.100.1"
)

var cm cablemodem.Modem

func printJSON(v interface{}) {
	output, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't encode data: %s", err)
		return
	}
	fmt.Printf("%s", output)

}

func cmdStatus(c *cli.Context) {
	status, err := cm.Status()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get cable modem status: %s", err)
		return
	}

	if c.GlobalBool("json") {
		printJSON(status)
		return
	}

	fmt.Printf("Uptime: %s\n", status.Uptime)
}

func cmdChannel(c *cli.Context) {
	signal, err := cm.SignalData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get signal data: %s", err)
		return
	}

	if c.GlobalBool("json") {
		printJSON(signal)
		return
	}

	fmt.Println("Upstream Channels:")
	for _, c := range signal.Upstream {
		fmt.Printf("  Channel: %d\n", c.ID)
		fmt.Printf("    Frequency          : %d Mhz\n", c.Freq/1000/1000)
		fmt.Printf("    Ranging Service ID : %d\n", c.RangingServiceID)
		fmt.Printf("    Modulation         : %s\n", c.Modulation)
		fmt.Printf("    Ranging Status     : %s\n", c.RangingStatus)
		fmt.Printf("    Power Level        : %d dBmV\n", c.Power)
	}
	fmt.Println("")

	fmt.Println("Downstream Channels:")
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
		cli.BoolFlag{
			Name:  "json",
			Usage: "Output data in json",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "channels",
			Aliases: []string{"c"},
			Usage:   "channel stats",
			Action:  cmdChannel,
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "cable modem status",
			Action:  cmdStatus,
		},
	}
	app.Before = func(c *cli.Context) error {
		cm = cablemodem.NewSurfboardCM(defaultString(c, "host", defaultHost))
		return nil
	}
	app.Run(os.Args)
}
