package tmailclient

import "github.com/awesome-gocui/gocui"

func (c *Client) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, c.nextPage); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, c.prevPage); err != nil {
		return err
	}
	return nil
}
