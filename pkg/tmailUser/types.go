package tmailuser

import (
	"context"
	"fmt"
	"log"
	"os"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type User struct {
	Srv             *gmail.Service
	Cache           tmailcache.Cache
	MsgRecvChan     chan tmailcache.MsgCacheEntry
	MsgNextPageChan chan bool
	MsgPageToken    string
}

func NewUser() User {
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

	return User{
		Srv:             srv,
		Cache:           tmailcache.NewCache(),
		MsgRecvChan:     make(chan tmailcache.MsgCacheEntry),
		MsgNextPageChan: make(chan bool),
	}
}

func (u *User) Listen() {
	//goroutine which adds messages received to the cache sequentially
    go u.listenForNextPage()
    go u.listenForMsgReceive()
}

func (u *User) listenForNextPage() {
	for {
		if <-u.MsgNextPageChan {
			u.messageScraper(10, 50)
		}
	}
}

func (u *User) listenForMsgReceive() {
	for {
		select {
		case msg := <-u.MsgRecvChan:
			u.Cache.AddToMessageCache(msg)
			fmt.Println(msg.Subject)
		}
	}
}
