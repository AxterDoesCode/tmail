package main

import (
	"github.com/AxterDoesCode/tmail/pkg/scraper"
	tmailuser "github.com/AxterDoesCode/tmail/pkg/tmailUser"
)

const concurrency int = 5

func main() {
	user := tmailuser.NewUser()

	go scraper.MessageScraper(&user, concurrency)
    go user.Listen()
    for {

    }
}
