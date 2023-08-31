package tmailclient

import "github.com/awesome-gocui/gocui"

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Client) nextPage(g *gocui.Gui, v *gocui.View) error {
	c.MsgPageTokenIndex++
	c.MsgChangePageChan <- struct{}{}
	return nil
}

func (c *Client) prevPage(g *gocui.Gui, v *gocui.View) error {
	if c.MsgPageTokenIndex > 0 {
		c.MsgPageTokenIndex--
		c.MsgChangePageChan <- struct{}{}
	}
	return nil
}
