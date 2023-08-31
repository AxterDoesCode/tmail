package tmailclient

import "github.com/awesome-gocui/gocui"

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Client) nextPage (g *gocui.Gui, v *gocui.View) error {
    c.MsgNextPageChan <- struct{}{}
    return nil
}
