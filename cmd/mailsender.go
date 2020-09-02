package main

import (
	"flag"
	"fmt"
	"os"

	"mailsender"
)

// usage constant provide help message
const usage = "Usage:\n  {-flags} \nExample: ./mailsender -p 8080"

func main() {
	// initialize command-line flags
	flagSet := flag.NewFlagSet("set", flag.ExitOnError)

	p := flagSet.Int("p", 8080, "-p {number} port for API server (default 8080)")
	s := flagSet.String("s", "mx.spamorez.ru:25", "-s {string} smtp server address (default mx.spamorez.ru:25)")
	e := flagSet.String("e", "info@uddug.com", "-e {string} contact form email address (default info@uddug.com)")

	// parse command-line flags
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		flagSet.PrintDefaults()
		os.Exit(1)
	}

	server := mailsender.NewServer(fmt.Sprintf(":%d", *p), *s, *e)
	server.Serve()
}
