package tmailuser

import (
	"context"
	"log"
	"os"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type User struct {
	Srv   *gmail.Service
	Cache tmailcache.Cache
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
        Srv: srv,
        Cache: tmailcache.NewCache() ,
    }
}
