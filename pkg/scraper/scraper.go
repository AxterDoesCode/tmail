package scraper

import (
	"fmt"
	"log"
	"sync"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	tmailuser "github.com/AxterDoesCode/tmail/pkg/tmailUser"
	"google.golang.org/api/gmail/v1"
)

func MessageScraper(user *tmailuser.User) {
	messages, err := user.Srv.Users.Messages.List("me").MaxResults(500).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
		return
	}

	fmt.Println("Message subjects:")
	fmt.Printf("Number of messages %v\n", len(messages.Messages))

	var wg sync.WaitGroup
	concurrencyLimit := 5
	semaphore := make(chan struct{}, concurrencyLimit)

	go func() {
		for {
			select {
			case msg := <-user.MsgRecvChan:
				user.Cache.AddToMessageCache(msg)
				fmt.Println("Added message", msg)
			}
		}
	}()

	for _, m := range messages.Messages {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(m *gmail.Message) {
			defer func() { <-semaphore }()
			scrapeMessage(m, user, &wg)
		}(m)
	}
	wg.Wait()
}

func scrapeMessage(m *gmail.Message, user *tmailuser.User, wg *sync.WaitGroup) {
	defer wg.Done()
	msg, err := user.Srv.Users.Messages.Get("me", m.Id).Do()
	if err != nil {
		fmt.Printf("Error retrieving message: %v", err)
		return
	}

	MessageEntry := tmailcache.MsgCacheEntry{}
	MessageEntry.Id = m.Id
	for _, h := range msg.Payload.Headers {
		switch h.Name {
		case "Subject":
			MessageEntry.Subject = h.Value
		case "To":
			MessageEntry.To = h.Value
		case "From":
			MessageEntry.From = h.Value
		}
	}
	user.MsgRecvChan <- MessageEntry
}
