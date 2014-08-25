package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/supershabam/nodeify/nodeify"
)

func main() {
	source := flag.String("source", "", "couch db http resource to poll")
	period := flag.Duration("period", time.Minute, "polling period")
	since := flag.Duration("since", time.Minute, "start requesting data from now minus since time ago")

	flag.Parse()

	fetcher, err := nodeify.NewFetcher(*source)
	if err != nil {
		log.Fatal(err)
	}
	last := time.Now().Add(-*since)
	for {
		now := last
		last = time.Now()
		modules, err := fetcher.Fetch(now)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", modules)
		time.Sleep(*period)
	}
}
