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

func (c *Client) cursorMovement(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		dir := 1
		if d < 0 {
			dir = -1
		}
		distance := int(math.Abs(float64(d)))
		for ; distance > 0; distance-- {
			if lineBelow(v, distance*dir) {
				v.MoveCursor(0, distance*dir)
				c.printMessageBody(g, v)
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

// This function doesnt need an async call since data is already stored in cache
// Note that this func requires the view to be "side"
// Prints message body to main view
func (c *Client) printMessageBody(g *gocui.Gui, v *gocui.View) error {
    //Gets the cursor poisition of the "side" view
	_, y := v.Cursor()
	v, err := g.View("main")
	if err != nil {
		return err
	}
	v.Clear()
	currentMessage := c.MsgCacheDisplay[y]
	fmt.Fprintf(v, "Date: %s\nFrom: %s\nType: %s\n\n", currentMessage.Date ,currentMessage.From, currentMessage.ContentType)
    fmt.Fprintf(v, "%v\n", currentMessage.LabelIds)
	fmt.Fprintln(v, currentMessage.Body)
	return nil
}

func focusMain(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView("main")
	if err != nil {
		return err
	}
	return nil
}

func focusSide(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView("side")
	if err != nil {
		return err
	}
	return nil
}
