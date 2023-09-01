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
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, c.refreshEmails); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'j', gocui.ModNone, c.cursorMovement(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'k', gocui.ModNone, c.cursorMovement(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, c.cursorMovement(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, c.cursorMovement(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'J', gocui.ModNone, c.cursorMovement(10)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'K', gocui.ModNone, c.cursorMovement(-10)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, focusMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEsc, gocui.ModNone, focusSide); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyTab, gocui.ModNone, c.nextTab); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyBacktab, gocui.ModNone, c.prevTab); err != nil {
		return err
	}
	return nil
}
