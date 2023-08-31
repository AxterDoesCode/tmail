package tmailclient

import (
	"fmt"
	"math"

	"github.com/awesome-gocui/gocui"
)

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

func (c *Client) refreshEmails(g *gocui.Gui, v *gocui.View) error {
	c.MsgPageTokenIndex = 0
	c.MsgChangePageChan <- struct{}{}
	return nil
}

func cursorMovement(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		dir := 1
		if d < 0 {
			dir = -1
		}
		distance := int(math.Abs(float64(d)))
		for ; distance > 0; distance-- {
			if lineBelow(v, distance*dir) {
				v.MoveCursor(0, distance*dir)
				return nil
			}
		}
		return nil
	}
}

func lineBelow(v *gocui.View, d int) bool {
	_, y := v.Cursor()
	line, err := v.Line(y + d)
	return err == nil && line != ""
}

func (c *Client) getBody(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
    v, err := g.View("main")
    if err != nil {
        return err
    }
    v.Clear()
    fmt.Fprintln(v, c.MsgCacheDisplay[y].Body)
	return nil
}
