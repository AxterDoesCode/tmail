package tmailclient

import (
	"context"
	"log"
	"os"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	"github.com/awesome-gocui/gocui"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Client struct {
	Srv *gmail.Service
	*gocui.Gui
	tmailcache.Cache
	MsgRecvChan     chan tmailcache.MsgCacheEntry
	MsgNextPageChan chan struct{}
	MsgPageToken    string
	RefreshGuiChan  chan struct{}
}

func NewClient() Client {
	ctx := context.Background()
	b, err := os.ReadFile("./credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return Client{}
	}

	return Client{
		Srv:             srv,
		Cache:           tmailcache.NewCache(),
		MsgRecvChan:     make(chan tmailcache.MsgCacheEntry),
		MsgNextPageChan: make(chan struct{}),
		RefreshGuiChan:  make(chan struct{}),
		Gui:             g,
	}
}

func (c *Client) Listen() {
	//goroutine which adds messages received to the cache sequentially
	go c.listenForNextPage()
	go c.listenForMsgReceive()
}

func (c *Client) listenForNextPage() {
	for {
		select {
		case <-c.MsgNextPageChan:
			c.messageScraper(10, 20)
		}
	}
}

func (c *Client) listenForMsgReceive() {
	for {
		select {
		case msg := <-c.MsgRecvChan:
			c.Cache.AddToMessageCache(msg)
			c.RefreshGuiChan <- struct{}{}
		}
	}
}
