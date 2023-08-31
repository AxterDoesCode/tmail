package tmailclient

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main", -1, -1, maxX, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Loading Emails...")
        if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
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
	g.SetManagerFunc(layout)

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

    //Initial email fetch on load
    c.MsgChangePageChan <- struct{}{}

	// this mainloop is blocking
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panic(err)
	}
}

func (c *Client) redrawCui(g *gocui.Gui) error {
	v, err := g.View("main")
	if err != nil {
		return err
	}
	v.Clear()

	for _, val := range c.Cache.MsgCacheDisplay {
		fmt.Fprintf(v, "%s\n", val.Subject)
	}
	return nil
}

