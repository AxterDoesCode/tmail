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
	y0offset := 2
	c.MaxResults = maxY - y0offset - 1

	if v, err := g.SetView("labels", 0, 0, maxX-1, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
        fmt.Fprintln(v, "Label Label Label")
	}

	if v, err := g.SetView("side", 0, y0offset, 40, maxY, 0); err != nil {
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

	if v, err := g.SetView("main", 40, y0offset, maxX-1, maxY, 0); err != nil {
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

    //Listening for message aggregation completion
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
    v, err := g.View("labels")
    if err != nil {
        return err
    }
    v.Clear()

    for _,l := range c.Labels {
        if c.CurrentLabel == l {
            fmt.Fprintf(v, "\033[34;7m%s\033[0m ", l)
        } else {
            fmt.Fprintf(v, "%s ", l)
        }
    }


	v, err = g.View("side")
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
            fmt.Fprintf(v, "\033[34;1m%s\033[0m\n", val.Subject)
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
