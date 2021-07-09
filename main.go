package main

import (
	"fmt"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyMainWin struct {
	*walk.MainWindow                // 窗口
	txtPathTE        *walk.LineEdit // txt文件的路径
	dataProj         []*string      // 数据项目的列表

}

// 打开数据文件，选择文件
func (mw *MyMainWin) selectFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择数据文件"
	dlg.Filter = "文本文件(*.txt)|*.txt"

	mw.txtPathTE.SetText("")
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.txtPathTE.SetText("选择文件过程出错，请联系老头！\r\n")
		return
	} else if !ok {
		mw.txtPathTE.SetText("未选择文件\r\n")
		return
	}
	s := fmt.Sprintf(" %s\r\n", dlg.FilePath)
	mw.txtPathTE.SetText(s)

}

func main() {
	//主窗口对象
	mw := MyMainWin{}

	err := MainWindow{
		Accessibility:    Accessibility{},
		Background:       nil,
		ContextMenuItems: []MenuItem{},
		DoubleBuffering:  false,
		Enabled:          nil,
		Font:             Font{},
		MaxSize:          Size{},
		MinSize:          Size{Width: 600, Height: 400},
		Name:             "",
		OnBoundsChanged: func() {
		},
		OnKeyDown: func(key walk.Key) {
		},
		OnKeyPress: func(key walk.Key) {
		},
		OnKeyUp: func(key walk.Key) {
		},
		OnMouseDown: func(x int, y int, button walk.MouseButton) {
		},
		OnMouseMove: func(x int, y int, button walk.MouseButton) {
		},
		OnMouseUp: func(x int, y int, button walk.MouseButton) {
		},
		OnSizeChanged: func() {
		},
		Persistent:         false,
		RightToLeftLayout:  false,
		RightToLeftReading: false,
		ToolTipText:        nil,
		Visible:            nil,
		Children: []Widget{
			GroupBox{
				Layout: HBox{},
				Title:  "导入数据",
				Font:   Font{PointSize: 10},
				Children: []Widget{
					Label{Text: "数据文件:", Font: Font{PointSize: 10}},
					LineEdit{AssignTo: &mw.txtPathTE},
					PushButton{Text: "浏览", OnClicked: mw.selectFile},
					PushButton{Text: "导入"},
				},
			},
		},
		DataBinder: DataBinder{},
		Layout:     VBox{},
		Icon:       nil,
		Size:       Size{},
		Title:      "Eiko的工具箱",
		AssignTo:   &mw.MainWindow,
		Bounds:     Rectangle{},
		MenuItems:  []MenuItem{},
		OnDropFiles: func([]string) {
		},
		StatusBarItems:    []StatusBarItem{},
		SuspendedUntilRun: false,
		ToolBar:           ToolBar{},
		ToolBarItems:      []MenuItem{},
	}.Create()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	mw.Run()
}
