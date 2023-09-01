package tmailclient

import (
	"errors"
	"fmt"
	"log"
	"sort"

	tmailcache "github.com/AxterDoesCode/tmail/pkg/tmailCache"
	"github.com/awesome-gocui/gocui"
)

func (c *Client) layout(g *gocui.Gui) error {
	//This kind of is looped through?
	if !c.GuiStarted {
		//Initial fetch
		c.MsgChangePageChan <- struct{}{}
		c.GuiStarted = true
	}
	//Setting the number of results to be the max rows of the terminal
	maxX, maxY := g.Size()
	y0offset := 5
	c.MaxResults = maxY - y0offset - 1

	if v, err := g.SetView("side", -1, y0offset, 40, maxY, 0); err != nil {
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

	if v, err := g.SetView("main", 40, y0offset, maxX, maxY, 0); err != nil {
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

// Redraws the cui after an api call to fetch emails
func (c *Client) redrawCui(g *gocui.Gui) error {
	v, err := g.View("side")
	if err != nil {
		return err
	}
	v.Clear()

	//Sorting the slice by their internal date (epoch time ms)
	sort.SliceStable(c.MsgCacheDisplay, func(i, j int) bool {
		return c.MsgCacheDisplay[i].InternalDate > c.MsgCacheDisplay[j].InternalDate
	})

	for _, val := range c.Cache.MsgCacheDisplay {
		if messageUnread(val) {
			fmt.Fprintf(v, "\x1b[0;36m%s\n", val.Subject)
		} else {
			fmt.Fprintf(v, "%s\n", val.Subject)
		}
	}

	//Prints the message body to main
	err = c.printMessageBody(g, v)
	if err != nil {
		return err
	}

	return nil
}

func messageUnread(m tmailcache.MsgCacheEntry) bool {
	for _, l := range m.LabelIds {
		if l == "UNREAD" {
			return true
		}
	}
	return false
}
