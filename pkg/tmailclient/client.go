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
	MsgChangePageChan chan struct{}
	MsgPageTokenMap   map[int]string
	MsgPageTokenIndex int
	MaxResults        int
	RefreshGuiChan    chan struct{}
	GuiStarted        bool
	Labels            []string
	CurrentLabel      string
}

func NewClient() *Client {
	ctx := context.Background()
	b, err := os.ReadFile("./credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	ret := &Client{
		Srv:               srv,
		Cache:             tmailcache.NewCache(),
		MsgChangePageChan: make(chan struct{}),
		RefreshGuiChan:    make(chan struct{}),
		MsgPageTokenMap:   make(map[int]string),
		GuiStarted:        false,
		Labels:            []string{"INBOX", "IMPORTANT", "SENT", "SPAM", "TRASH"},
		CurrentLabel:      "INBOX",
	}
	ret.Listen()
	return ret
}

func (c *Client) Listen() {
	//goroutine which adds messages received to the cache sequentially
	go c.listenForPageChange()
}

func (c *Client) listenForPageChange() {
	for {
		select {
		case <-c.MsgChangePageChan:
			c.messageScraper()
		}
	}
}
