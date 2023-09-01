package tmailclient

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	"google.golang.org/api/gmail/v1"
)

func (c *Client) messageScraper() {
	messages, err := c.Srv.Users.Messages.List("me").PageToken(c.MsgPageTokenMap[c.MsgPageTokenIndex]).MaxResults(int64(c.MaxResults)).LabelIds(c.CurrentLabel).Do()
	if err != nil {
		log.Printf("Unable to retrieve messages: %v\n", err)
		return
	}

	c.MsgPageTokenMap[c.MsgPageTokenIndex+1] = messages.NextPageToken

	//Reset the content to be displayed
	c.MsgCacheDisplay = []tmailcache.MsgCacheEntry{}
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, c.MaxResults)
	for _, m := range messages.Messages {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(m *gmail.Message) {
			defer func() { <-semaphore }()
			msgEntry, err := c.fetchMessage(m, &wg)
			if err != nil {
				log.Println(err)
				return
			}
			c.MsgCacheMu.Lock()
			c.AddToMessageCache(msgEntry)
			c.AddToMessageCacheDisplay(msgEntry)
			c.MsgCacheMu.Unlock()
		}(m)
	}
	wg.Wait()
	c.RefreshGuiChan <- struct{}{}
}

func (c *Client) fetchMessage(m *gmail.Message, wg *sync.WaitGroup) (*tmailcache.MsgCacheEntry, error) {
	defer wg.Done()

	//Checks if the entry already is in cache ()
	if k, ok := c.MsgCache[m.Id]; ok {
		msg, err := c.Srv.Users.Messages.Get("me", m.Id).Format("minimal").Do()
		if err != nil {
			log.Printf("Error retrieving message: %v", err)
			return nil, err
		}
		k.LabelIds = msg.LabelIds

		return &k, nil
	}

	//Msg isn't in cache so fetch the whole body/do raw data parsing
	msg, err := c.Srv.Users.Messages.Get("me", m.Id).Do()
	if err != nil {
		log.Printf("Error retrieving message: %v", err)
		return nil, err
	}

	MessageEntry := tmailcache.MsgCacheEntry{}

	MessageEntry.Id = msg.Id
	MessageEntry.InternalDate = msg.InternalDate
	MessageEntry.LabelIds = msg.LabelIds

	decodedBody, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
	if err != nil {
		return nil, err
	}
	MessageEntry.Body = string(decodedBody)

	//Store relevant payload headers in cache
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
		case "Date":
			MessageEntry.Date = h.Value
		case "Reply-To":
			MessageEntry.ReplyTo = h.Value
		case "Return-Path":
			MessageEntry.ReturnPath = h.Value
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
