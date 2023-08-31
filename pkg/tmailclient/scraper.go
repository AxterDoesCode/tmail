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
	messages, err := c.Srv.Users.Messages.List("me").PageToken(c.MsgPageTokenMap[c.MsgPageTokenIndex]).MaxResults(maxResults).Do()
	if err != nil {
		log.Printf("Unable to retrieve messages: %v\n", err)
		return
	}

    if _, ok := c.MsgPageTokenMap[c.MsgPageTokenIndex + 1]; !ok {
        c.MsgPageTokenMap[c.MsgPageTokenIndex + 1] = messages.NextPageToken
    }

    //Reset the content to be displayed
	c.MsgCacheDisplay = make(map[string]tmailcache.MsgCacheEntry)
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrency)
	for _, m := range messages.Messages {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(m *gmail.Message) {
			defer func() { <-semaphore }()
			msgEntry, new, err := c.fetchMessage(m, &wg)
			if err != nil {
				log.Println(err)
				return
			}
			c.MsgCacheMu.Lock()
			if new {
				c.AddToMessageCache(msgEntry)
			}
			c.AddToMessageCacheDisplay(msgEntry)
			c.MsgCacheMu.Unlock()
		}(m)
	}
	wg.Wait()
	c.RefreshGuiChan <- struct{}{}
}

func (c *Client) fetchMessage(m *gmail.Message, wg *sync.WaitGroup) (*tmailcache.MsgCacheEntry, bool, error) {
	defer wg.Done()

    //Checks if the entry already is in cache ()
    //Maybe this needs to be changed because im creating a pointer of a value???
    if k, ok := c.MsgCache[m.Id]; ok{
        return &k, false, nil
    }

	msg, err := c.Srv.Users.Messages.Get("me", m.Id).Do()
	if err != nil {
		log.Printf("Error retrieving message: %v", err)
		return nil, false, err
	}

	MessageEntry := tmailcache.MsgCacheEntry{}

	MessageEntry.Id = msg.Id
	decodedBody, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
	if err != nil {
		return nil, false, err
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

	return &MessageEntry, true, nil
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
