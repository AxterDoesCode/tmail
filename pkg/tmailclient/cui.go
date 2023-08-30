package tmailclient

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func layout(g *gocui.Gui) error {
	_, maxY := g.Size()

	if v, err := g.SetView("main", -1, -1, 30, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Testing")
	}
	return nil
}
func (c *Client) StartCui() {
	g := c.Gui
	defer g.Close()
	g.Cursor = true
	g.SetManagerFunc(layout)

	go func() {
		for {
			select {
			case <-c.RefreshGuiChan:
				g.Update(c.test)
			}
		}
	}()

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

}

func (c *Client) test(g *gocui.Gui) error {
	if v, err := g.View("main"); err == nil {
		v.Clear()
		for _, val := range c.Cache.MsgCache {
			fmt.Fprintf(v, "%s\n", val.Subject)
		}
	}
    return nil
}
