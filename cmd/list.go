package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/urfave/cli"
)

// NewListCMD to list information.
func NewListCMD() cli.Command {
	return cli.Command{
		Name:   "list",
		Usage:  "list information",
		Action: listAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "start, s",
				Usage: "The start day to search like YYYY-MM-DD",
			},
			cli.StringFlag{
				Name:  "end, e",
				Usage: "The end day to search like YYYY-MM-DD",
			},
		},
	}
}

func newCli() *http.Client {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func listAction(c *cli.Context) error {
	start := c.String("start")
	if len(start) == 0 {
		return errors.New("please specify a start day")
	}

	end := c.String("end")
	if len(end) == 0 {
		return errors.New("please specify an end day")
	}

	v := url.Values{}
	v.Add("rb", start)
	v.Add("re", end)
	v.Add("key", c.GlobalString("key"))
	v.Add("format", "json")
	v.Add("rs", "day")

	p := fmt.Sprintf("https://www.rescuetime.com/anapi/data?%s", v.Encode())
	cli := newCli()
	resp, err := cli.Get(p)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var respMsg map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&respMsg); err != nil {
			return err
		}
		fmt.Println(respMsg["messages"])
	} else {
		var respMsg map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&respMsg); err != nil {
			return err
		}
		fmt.Printf("from: %s, to: %s\n", start, end)
		items, ok := respMsg["rows"].([]interface{})
		if !ok {
			fmt.Println("parse items wrong")
			return nil
		}
		for i := 0; i < len(items); i++ {
			one, ok := items[i].([]interface{})
			if !ok {
				fmt.Println("parse item wrong")
				return nil
			}
			seconds, ok := one[1].(float64)
			if !ok {
				fmt.Println("parse seconds wrong")
				return nil
			}
			timeString := formatSeconds(int64(seconds))

			name, ok := one[3].(string)
			if !ok {
				fmt.Println("parse name wrong")
				return nil
			}
			fmt.Printf("%d: %s, %s\n", i+1, name, timeString)
		}
	}

	return nil
}

// formatSeconds to format seconds to human readable format
func formatSeconds(seconds int64) string {
	days := seconds / 86400
	left := seconds - days*86400
	hours := left / 3600
	left = left - hours*3600
	minutes := left / 60
	left = left - minutes*60

	return fmt.Sprintf("%2dD%2dH%2dM%2dS", days, hours, minutes, left)
}
