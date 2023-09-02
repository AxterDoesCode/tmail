package tmailclient

import (
	"errors"
	"fmt"
	"math"

	"github.com/awesome-gocui/gocui"
	"google.golang.org/api/gmail/v1"
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

func (c *Client) nextTab(g *gocui.Gui, v *gocui.View) error {
	for i, l := range c.Labels {
		if c.CurrentLabel == l {
			labelIndex := (i + 1 + len(c.Labels)) % len(c.Labels)
			c.CurrentLabel = c.Labels[labelIndex]
			break
		}
	}
	c.refreshEmails(g, v)
	return nil
}

func (c *Client) prevTab(g *gocui.Gui, v *gocui.View) error {
	for i, l := range c.Labels {
		if c.CurrentLabel == l {
			labelIndex := (i - 1 + len(c.Labels)) % len(c.Labels)
			c.CurrentLabel = c.Labels[labelIndex]
			break
		}
	}
	c.refreshEmails(g, v)
	return nil
}

func (c *Client) selectTab(t int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		c.CurrentLabel = c.Labels[t]
		c.refreshEmails(g, v)
		return nil
	}
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


func scrollMessage(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		dir := 1
		if d < 0 {
			dir = -1
		}
		distance := int(math.Abs(float64(d)))
		for ; distance > 0; distance-- {
			if lineBelow2(v, distance*dir) {
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

func lineBelow2(v *gocui.View, d int) bool {
	_, y := v.Cursor()
	_, err := v.Line(y + d)
	return err == nil 
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
	c.CurrentMessage = currentMessage
	fmt.Fprintf(v, "ID: %s\nDate: %s\nFrom: %s\nType: %s\n\n", currentMessage.Id, currentMessage.Date, currentMessage.From, currentMessage.ContentType)
	fmt.Fprintf(v, "Reply-To: %s\nReturn-Path: %s\n", currentMessage.ReplyTo, currentMessage.ReturnPath)
	//Add reply to and return path and check the output
	fmt.Fprintf(v, "%v\n", currentMessage.LabelIds)
	fmt.Fprintln(v, currentMessage.Body)
	return nil
}

func (c *Client) focusMain(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	line, err := v.Line(y)
	if err != nil {
		return err
	}
	v.SetLine(y, line)
	go func() {
		mrq := gmail.ModifyMessageRequest{}
		mrq.RemoveLabelIds = []string{"UNREAD"}
		c.Srv.Users.Messages.Modify("me", c.CurrentMessage.Id, &mrq).Do()
	}()
	_, err = setCurrentViewOnTop(g, "main")
	if err != nil {
		return err
	}
	return nil
}

func focusSide(g *gocui.Gui, v *gocui.View) error {
	_, err := setCurrentViewOnTop(g, "side")
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) openReplyView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	_, err := g.SetView("reply", 10, 10, maxX-10, maxY-10, 0)
	if !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	v, err = setCurrentViewOnTop(g, "reply")
	v.Editable = true
	v.Wrap = true
	if err != nil {
		return err
	}

	return nil
}

func closeReplyView(g *gocui.Gui, v *gocui.View) error {
	err := g.DeleteView("reply")
	if err != nil {
		return err
	}
	_, err = setCurrentViewOnTop(g, "side")
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) sendMessage(g *gocui.Gui, v *gocui.View) error {
	msg := gmail.Message{}
	//msgContent := v.ViewBuffer()

	//msg.Header.Add("To", value string)
	//use
	//Need to base64 encode message and some other stuff
	c.Srv.Users.Messages.Send("me", &msg)
	return nil
}
