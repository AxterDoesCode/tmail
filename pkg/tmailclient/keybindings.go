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
	if err := g.SetKeybinding("main", 'j', gocui.ModNone, scrollMessage(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'k', gocui.ModNone, scrollMessage(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'J', gocui.ModNone, c.cursorMovement(10)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'K', gocui.ModNone, c.cursorMovement(-10)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, c.focusMain); err != nil {
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
	if err := g.SetKeybinding("side", '1', gocui.ModNone, c.selectTab(0)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", '2', gocui.ModNone, c.selectTab(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", '3', gocui.ModNone, c.selectTab(2)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", '4', gocui.ModNone, c.selectTab(3)); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", '5', gocui.ModNone, c.selectTab(4)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'r', gocui.ModNone, c.openReplyView); err != nil {
		return err
	}
	if err := g.SetKeybinding("reply", gocui.KeyEsc, gocui.ModNone, closeReplyView); err != nil {
		return err
	}
	if err := g.SetKeybinding("reply", gocui.KeyCtrlS, gocui.ModNone, c.sendMessage); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'D', gocui.ModNone, c.deleteMessage); err != nil {
		return err
	}

	return nil
}
