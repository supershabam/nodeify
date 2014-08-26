package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/supershabam/nodeify"
)

func main() {
	source := flag.String("source", "", "couch db http resource to poll")
	period := flag.Duration("period", time.Minute, "polling period")
	since := flag.Duration("since", time.Minute, "start requesting data from now minus since time ago")

	flag.Parse()

	u, err := url.Parse(*source)
	if err != nil {
		log.Fatal(err)
	}
	consumer := nodeify.Consumer{
		Fetcher: &nodeify.HTTPFetcher{
			URL: u,
		},
		Since:  time.Now().Add(-*since),
		Period: *period,
	}
	for module := range consumer.Consume() {
		fmt.Printf("%+v\n", module)
	}
	if err := consumer.Err(); err != nil {
		log.Fatal(err)
	}
}
