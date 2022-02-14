package widget

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/awesome-gocui/gocui"
	"strings"
)

var (
	HelpGlobal = ""
)

type HelpWidget struct {
	name string
	x, y int
	w, h int
	body string
}

func NewHelpWidget() {
	initContent()

	maxX, _ := commConsts.Cui.Size()

	lines := strings.Split(HelpGlobal, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines)
	w = w + 2

	help := HelpWidget{name: "help", x: maxX - w - 3, y: commConsts.MainViewHeight, w: w, h: h, body: HelpGlobal}
	help.Layout()
}

func (w *HelpWidget) Layout() error {
	v, err := commConsts.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		fmt.Fprint(v, w.body)
	}
	return nil
}

func ShowHelpFromView(g *gocui.Gui, v *gocui.View) error {
	return ShowHelp()
}
func ShowHelp() error {
	help, _ := commConsts.Cui.View("help")

	if help != nil {
		HideHelp()
	} else {
		NewHelpWidget()
	}

	return nil
}

func HideHelp() error {
	help, _ := commConsts.Cui.View("help")

	if help != nil {
		if err := commConsts.Cui.DeleteView("help"); err != nil {
			return err
		}
	}

	return nil
}

func initContent() {
	HelpGlobal = fmt.Sprintf("%s \n %s \n %s \n",
		i118Utils.Sprintf("help_key_bind"),
		i118Utils.Sprintf("help_show"),
		i118Utils.Sprintf("help_exit"))
}
