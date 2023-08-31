package tmailclient

import (
	"errors"
	"fmt"
	"log"
	"os"

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

func keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        return err
    }
    return nil
}

func (c *Client) StartCui() {
    g := c.Gui
	defer g.Close()
	g.Cursor = true
	g.SetManagerFunc(layout)

    if err := keybindings(g); err != nil {
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
	} else if errors.Is(err, gocui.ErrQuit) {
        os.Exit(0)
    }

}

func (c *Client) redrawCui(g *gocui.Gui) error {
	v, err := g.View("main")
    if err != nil {
        return err
    }
    v.Clear()
    for _, val := range c.Cache.MsgCache {
        fmt.Fprintf(v, "%s\n", val.Subject)
    }
    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    //Need to fix cursor being in weird position after
    //And need to clear the view once it's exited the app
    v, err := g.View("")
    if err != nil {
        fmt.Println("here")
        return err
    }
    v.Clear()
    return gocui.ErrQuit
}
