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

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrency)
	for _, m := range messages.Messages {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(m *gmail.Message) {
			defer func() { <-semaphore }()
			msgEntry, err := c.scrapeMessage(m, &wg)
			if err != nil {
				log.Println(err)
			}
            c.MsgCacheMu.Lock()
            c.AddToMessageCache(msgEntry)
            c.MsgCacheMu.Unlock()
		}(m)
	}
	wg.Wait()
    c.RefreshGuiChan <- struct{}{}
}

func (c *Client) scrapeMessage(m *gmail.Message, wg *sync.WaitGroup) (*tmailcache.MsgCacheEntry, error) {
	defer wg.Done()
	msg, err := c.Srv.Users.Messages.Get("me", m.Id).Do()
	if err != nil {
		log.Printf("Error retrieving message: %v", err)
		return nil, err
	}

	MessageEntry := tmailcache.MsgCacheEntry{}

	MessageEntry.Id = msg.Id
	decodedBody, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
	if err != nil {
		return nil, err
	}
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

	return &MessageEntry, nil
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
