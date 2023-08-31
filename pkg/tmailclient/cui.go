package tmailclient

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func (c *Client) layout(g *gocui.Gui) error {
	if !c.GuiStarted {
		//Initial fetch
		c.MsgChangePageChan <- struct{}{}
		c.GuiStarted = true
	}

	//Setting the number of results to be the max rows of the terminal
	maxX, maxY := g.Size()
	c.MaxResults = maxY

	if v, err := g.SetView("side", -1, 10, 40, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Loading Emails...")
		if _, err := g.SetCurrentView("side"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("main", 40, 10, maxX, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
        v.Wrap = true
        fmt.Fprintln(v, "Select an email")
	}

	return nil
}

func (c *Client) StartCui() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return
	}

	defer g.Close()
	c.Gui = g

	g.Cursor = true
	g.SetManagerFunc(c.layout)

	if err := c.keybindings(g); err != nil {
		log.Panicln(err)
	}

	go func() {
		for {
			select {
			case <-c.RefreshGuiChan:
				g.UpdateAsync(c.redrawCui)
			}
		}
	}()

	// this mainloop is blocking
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panic(err)
	}
}

func (c *Client) redrawCui(g *gocui.Gui) error {
	v, err := g.View("side")
	if err != nil {
		return err
	}
	v.Clear()

	for _, val := range c.Cache.MsgCacheDisplay {
		fmt.Fprintf(v, "%s\n", val.Subject)
	}
	return nil
}
