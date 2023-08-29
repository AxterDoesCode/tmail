package scraper

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	tmailuser "github.com/AxterDoesCode/tmail/pkg/tmailUser"
	"google.golang.org/api/gmail/v1"
)

func MessageScraper(user *tmailuser.User, concurrency int) {
	messages, err := user.Srv.Users.Messages.List("me").PageToken(user.MsgPageToken).MaxResults(5).Do()
	if err != nil {
		log.Printf("Unable to retrieve messages: %v\n", err)
		return
	}
    user.MsgPageToken = messages.NextPageToken
	//fmt.Printf("Number of messages %v\n", len(messages.Messages))

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrency)

    //goroutine which adds messages received to the cache sequentially
	go func() {
		for {
			select {
			case msg := <-user.MsgRecvChan:
				user.Cache.AddToMessageCache(msg)
				fmt.Println(msg.Id)
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
    fmt.Println("Next pagetoken is: ", user.MsgPageToken)
}

func scrapeMessage(m *gmail.Message, user *tmailuser.User, wg *sync.WaitGroup) {
	defer wg.Done()
	msg, err := user.Srv.Users.Messages.Get("me", m.Id).Do()
	if err != nil {
		fmt.Printf("Error retrieving message: %v", err)
		return
	}

	MessageEntry := tmailcache.MsgCacheEntry{}

	MessageEntry.Id = msg.Id
	decodedBody, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
	MessageEntry.Body = string(decodedBody)

	for _, h := range msg.Payload.Headers {
		switch h.Name {
		case "Subject":
			MessageEntry.Subject = h.Value
		case "To":
			MessageEntry.To = h.Value
		case "From":
			MessageEntry.From = h.Value
		case "Content-Type":
			MessageEntry.ContentType = h.Value
		}
	}
	//I need to handle non plaintext content sometime
	//probably by getting the raw data and parsing the html

	user.MsgRecvChan <- MessageEntry
}

func getRawMessageData(m *gmail.Message, user *tmailuser.User) {
	res, err := user.Srv.Users.Messages.Get("me", m.Id).Format("RAW").Do()
	if err != nil {
		log.Println("Error when getting raw mail content: ", err)
		return
	}
	decodedData, err := base64.URLEncoding.DecodeString(res.Raw)
	if err != nil {
		log.Println("Error decoding raw message body: ", err)
		return
	}

	fmt.Printf("- %s\n", decodedData)
}
