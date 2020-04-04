package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/jessevdk/go-flags"
)

func main() {
	options := NewOptions()
	parser := flags.NewParser(&options, flags.Default)
	parser.NamespaceDelimiter = "-"
	_, err := parser.Parse()
	if err != nil {
		if !flags.WroteHelp(err) {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	if len(options.Tasks) == 0 {
		fmt.Println("Please supply a file or an environment variable for configuration")
		os.Exit(1)
	}

	c := &cfclient.Config{
		ApiAddress:        options.API,
		Username:          options.Username,
		Password:          options.Password,
		SkipSslValidation: options.SkipSslValidation,
	}
	client, err := cfclient.NewClient(c)
	if err != nil {
		panic(err)
	}
	q := url.Values{}
	apps, err := client.ListAppsByQuery(q)
	if err != nil {
		panic(err)
	}
	for _, app := range apps {
		fmt.Printf("Name: %s (%d instances)\n", app.Name, app.Instances)
	}

	// 3 worker pools, configured by env/flags for size:
	//   <-(Filter, Org)
	//   <-(Filter, Space)
	//   <-(Filter, App)
	// Per task worker:
	//   <-(Tick) and delivers to proper pool based on starting filter
}
