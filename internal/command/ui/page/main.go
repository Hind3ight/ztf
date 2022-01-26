package page

import (
	"github.com/aaronchen2k/deeptest/internal/command/ui"
	"github.com/aaronchen2k/deeptest/internal/command/ui/widget"
	constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"
	"github.com/aaronchen2k/deeptest/internal/command/utils/vari"
	"github.com/awesome-gocui/gocui"
	"log"
)

func InitMainPage() error {
	maxX, maxY := vari.Cui.Size()
	if maxX < constant.MinWidth {
		maxX = constant.MinWidth
	}
	if maxY < constant.MinHeight {
		maxY = constant.MinHeight
	}
	vari.MainViewHeight = maxY - constant.CmdViewHeight

	mainView := widget.NewPanelWidget("main", 0, 0, maxX-2, vari.MainViewHeight, "")

	ui.ViewMap["root"] = append(ui.ViewMap["root"], mainView.Name())

	cmdView := widget.NewPanelWidget("cmd", 0, vari.MainViewHeight, maxX-2, constant.CmdViewHeight-1, "")
	ui.ViewMap["root"] = append(ui.ViewMap["root"], cmdView.Name())

	widget.NewHelpWidget()
	MainPageKeyBindings()

	return nil
}

func MainPageKeyBindings() error {
	if err := vari.Cui.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, widget.ShowHelpFromView); err != nil {
		log.Panicln(err)
	}
	if err := vari.Cui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		return err
	}

	ui.SupportScroll("cmd")

	v, _ := vari.Cui.View("cmd")
	v.Autoscroll = true

	return nil
}
