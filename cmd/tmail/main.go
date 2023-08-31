package main

import (
	"github.com/AxterDoesCode/tmail/pkg/tmailclient"
)

func main() {
	client := tmailclient.NewClient()
    client.StartCui()
	//This shows that next page trigger works, should refactor to be struct{}{} instead of true
	//Also should probably do a waitgroup instead of a shitty for loop actually maybe not because the cli app will be blocking
}
