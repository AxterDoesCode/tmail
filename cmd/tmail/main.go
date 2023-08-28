package main

import (
	"fmt"
	"log"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
)

func main() {
	user := NewUser()

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

	messages, err := user.Srv.Users.Messages.List("me").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
		return
	}

	fmt.Println("Message subjects:")
	fmt.Printf("Number of messages %v\n", len(messages.Messages))

    //Make this concurrent, thinking need to store message ID's to scrape in an array then create a scraper like in BlogAggregator
	for _, m := range messages.Messages {
		msg, err := user.Srv.Users.Messages.Get("me", m.Id).Do()
		if err != nil {
			fmt.Printf("Error retrieving message: %v", err)
			return
		}

		MessageEntry := tmailcache.MsgCacheEntry{}
		for _, h := range msg.Payload.Headers {
			switch h.Name {
            case "Subject" : MessageEntry.Subject = h.Value
            case "To" : MessageEntry.To = h.Value
            case "From" : MessageEntry.From = h.Value
			}
		}
        fmt.Println(MessageEntry.Subject)
	}

}
