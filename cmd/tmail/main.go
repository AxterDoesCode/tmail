package main

import (
	tmailuser "github.com/AxterDoesCode/tmail/pkg/tmailUser"
)


func main() {
	user := tmailuser.NewUser()

    go user.Listen()
    user.MsgNextPageChan <- true
    user.MsgNextPageChan <- true
    for {

    }
}
