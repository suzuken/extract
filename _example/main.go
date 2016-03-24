package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/suzuken/extract"
)

func main() {
	var (
		rawurl = flag.String("url", "http://example.com", "url for extract")
	)
	flag.Parse()
	ex := extract.New()
	if rule := ex.Match(*rawurl); rule == nil {
		log.Printf("%s doesn't match in rule", *rawurl)
		os.Exit(0)
	}
	c, err := ex.ExtractURL(*rawurl)
	if err != nil {
		log.Fatalf("extract failed: %s", err)
	}
	fmt.Printf("content: %v", c)
}
