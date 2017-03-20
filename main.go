package main

import (
	"github.com/zuriu/mobalysis-snowcap-cloud-api/snowcapapi"

	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	exampleFlag := flag.String("example", "", "{query-events|csv-events-by-index|csv-events-by-time}")
	publicKeyFlag := flag.String("publickey", "", "public key provided by Mobalysis")
	privateKeyFlag := flag.String("privatekey", "", "private key provided by Mobalysis")
	hostFlag := flag.String("host", "", "host to query (e.g. myresort.snowcapalerts.com)")
	useTLSFlag := flag.Bool("usetls", true, "whether to use connect using TLS")
	flag.Parse()

	example := *exampleFlag
	publicKey := *publicKeyFlag
	privateKey := *privateKeyFlag
	host := *hostFlag
	prefix := "http"
	if *useTLSFlag {
		prefix = "https"
	}

	if len(host) < 1 {
		flag.Usage()
		return
	}

	switch example {
	case "query-events":

		res, err := snowcapapi.Get(fmt.Sprintf("%s://%s/api/events?offset=0&limit=10", prefix, host), publicKey, privateKey)
		if err != nil {
			log.Println(fmt.Sprintf("could not get events: %s", err))
			return
		}
		if res.Error != "" {
			log.Println(res.Error)
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Println(fmt.Sprintf("could not get events due to unexpected response status code (StatusCode=%d)", res.StatusCode))
			return
		}

		events, err := res.ParseEvents(privateKey)
		if err != nil {
			log.Println(fmt.Sprintf("could not parse events: %s", err))
			return
		}

		log.Println(fmt.Sprintf("There are a total of %d events. Here are the 10 most recent events:", events.TotalCount))
		for _, event := range events.Results {
			log.Println(event)
		}

		log.Println(res.StatusCode, res.Body)

	case "csv-events-by-index":
		res, err := snowcapapi.Get(fmt.Sprintf("%s://%s/api/events.csv?offset=2&limit=200&skip=0&fields=created_on,depth", prefix, host), publicKey, privateKey)
		if err != nil {
			log.Println(fmt.Sprintf("could not get events: %s", err))
			return
		}
		if res.Error != "" {
			log.Println(fmt.Sprintf("%s (StatusCode=%d)", res.Error, res.StatusCode))
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Println(fmt.Sprintf("could not get events due to unexpected response status code (StatusCode=%d)", res.StatusCode))
			return
		}

		log.Println(res.StatusCode, res.Body)

	case "csv-events-by-time":
		res, err := snowcapapi.Get(fmt.Sprintf("%s://%s/api/events.csv?offset=2017-03-09T00:00:00Z&limit=2017-03-10T00:00:00Z&fields=created_on,depth", prefix, host), publicKey, privateKey)
		if err != nil {
			log.Println(fmt.Sprintf("could not get events: %s", err))
			return
		}
		if res.Error != "" {
			log.Println(fmt.Sprintf("%s (StatusCode=%d)", res.Error, res.StatusCode))
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Println(fmt.Sprintf("could not get events due to unexpected response status code (StatusCode=%d)", res.StatusCode))
			return
		}

		log.Println(res.StatusCode, res.Body)

	default:
		flag.Usage()
	}
}
