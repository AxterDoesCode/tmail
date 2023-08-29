package tmailuser

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	"google.golang.org/api/gmail/v1"
)

func (user *User) messageScraper(concurrency int, maxResults int64) {
	messages, err := user.Srv.Users.Messages.List("me").PageToken(user.MsgPageToken).MaxResults(maxResults).Do()
	if err != nil {
		log.Printf("Unable to retrieve messages: %v\n", err)
		return
	}
	user.MsgPageToken = messages.NextPageToken
	//fmt.Printf("Number of messages %v\n", len(messages.Messages))

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrency)
    counter := 0
	for _, m := range messages.Messages {
        counter++
		wg.Add(1)
		semaphore <- struct{}{}
		go func(m *gmail.Message) {
			defer func() { <-semaphore }()
			user.scrapeMessage(m, &wg)
		}(m)
	}
	wg.Wait()
    fmt.Println("Counter is: ", counter)
}

func (user *User) scrapeMessage(m *gmail.Message, wg *sync.WaitGroup) {
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

func (user *User) getRawMessageData(m *gmail.Message) {
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
