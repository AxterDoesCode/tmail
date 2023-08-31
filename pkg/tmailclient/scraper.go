package tmailclient

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	"google.golang.org/api/gmail/v1"
)

func (c *Client) messageScraper(concurrency int, maxResults int64) {
	messages, err := c.Srv.Users.Messages.List("me").PageToken(c.MsgPageToken).MaxResults(maxResults).Do()
	if err != nil {
		log.Printf("Unable to retrieve messages: %v\n", err)
		return
	}
	c.MsgPageToken = messages.NextPageToken
	//fmt.Printf("Number of messages %v\n", len(messages.Messages))

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrency)
	for _, m := range messages.Messages {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(m *gmail.Message) {
			defer func() { <-semaphore }()
			c.scrapeMessage(m, &wg)
		}(m)
	}
	wg.Wait()
}

func (c *Client) scrapeMessage(m *gmail.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	msg, err := c.Srv.Users.Messages.Get("me", m.Id).Do()
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

	c.MsgRecvChan <- MessageEntry
}

func (c *Client) getRawMessageData(m *gmail.Message) {
	res, err := c.Srv.Users.Messages.Get("me", m.Id).Format("RAW").Do()
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