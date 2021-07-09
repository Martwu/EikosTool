package main

import (
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var inTE, outTE *walk.TextEdit // 声明两个文本编辑控件

	//主窗口对象
	MainWindow{
		Title:   "Eiko利刃",       // 窗口标题设置
		MinSize: Size{600, 400}, //窗体的大小
		Layout:  VBox{},         // 窗体的布局形式
		//定义vbox的所有控件
		Children: []Widget{ //定义控件
			HSplitter{ //水平分割控件
				Children: []Widget{ //定义子控件
					TextEdit{AssignTo: &inTE},
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			PushButton{ //按钮控件
				Text: "确定",
				OnClicked: func() {
					outTE.SetText(strings.ToUpper(inTE.Text()))
				},
			},
		},
	}.Run()
}
