package main

import (
	"fmt"
	"log"

	"github.com/AxterDoesCode/tmail/pkg/scraper"
	tmailuser "github.com/AxterDoesCode/tmail/pkg/tmailUser"
)

func main() {
	user := tmailuser.NewUser()

	r, err := user.Srv.Users.Labels.List("me").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
    scraper.MessageScraper(&user)
    fmt.Println("Finished scraper exec")
}
