package deployment

import (
	"log"
	"os"
)

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return hostname
}
