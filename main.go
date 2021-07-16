package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type LsBoxModel struct {
	walk.ListModelBase
	items []string
}

func (m *LsBoxModel) ItemCount() int {
	return len(m.items)
}

// 列表中展示的值
func (m *LsBoxModel) Value(index int) interface{} {
	return m.items[index]
}

type projSet struct {
	m map[string]interface{}
}

func (ps *projSet) Add(item string) {
	ps.m[item] = true
}

func (ps *projSet) ToStrings() []string {
	list := []string{}
	for item := range ps.m {
		list = append(list, item)
	}
	return list
}

type MyMainWin struct {
	*walk.MainWindow                // 窗口
	fontsize         int            //字体大小
	txtPathTE        *walk.LineEdit // txt文件的路径
	projLsBox        *walk.ListBox  // 数据项目的选择列表
	projModel        *LsBoxModel    // 存数据有哪些项目
	proj             string         // 存储已选定的数据项目
	dataPath         string         // 存储已加载的数据的路径
}

// 打开数据文件，选择文件
func (mw *MyMainWin) selectFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择数据文件"
	dlg.Filter = "文本文件(*.txt)|*.txt"

	mw.txtPathTE.SetText("")
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.doPrompt("选择文件过程出错，请联系老头！")
		return
	} else if !ok {
		mw.doPrompt("未选择文件")
		return
	}
	s := fmt.Sprintf(" %s\r\n", dlg.FilePath)

	mw.txtPathTE.SetText(s)
}

//定义项目的单击行为
func (mw *MyMainWin) doCurrentIndexChanged() {
	i := mw.projLsBox.CurrentIndex()
	item := &mw.projModel.items[i]

	fmt.Println("CurrentIndex: ", i)
	fmt.Println("CurrentEnvVarName: ", item)
}

// 定义项目的双击行为
func (mw *MyMainWin) doItemActivated() {
	value := mw.projModel.items[mw.projLsBox.CurrentIndex()]
	walk.MsgBox(mw, "value", value, walk.MsgBoxIconInformation)
}

// 老头温馨提示
func (mw *MyMainWin) doPrompt(msg string) {
	walk.MsgBox(mw, "老头提示", msg, walk.MsgBoxOK)
}

func (mw *MyMainWin) doDropFiles(filepath []string) {
	if len(filepath) > 1 {
		mw.doPrompt("一次只可以处理一个文件而已哦！")
		return
	}
	if !strings.HasSuffix(filepath[0], ".txt") {
		mw.doPrompt("目前只支持txt文件啦。")
		return
	}
	mw.txtPathTE.SetText(filepath[0])
}

func (mw *MyMainWin) loadK4840Data(path string, infunc func(*csv.Reader)) func() {
	return func() {
		fs, err := os.OpenFile(strings.TrimSpace(path), os.O_RDONLY, 4)
		if err != nil {
			mw.doPrompt(fmt.Sprintf("打不开这个文件喔, 错误信息是：\n\r%+v", err))
			return
		}
		defer fs.Close()
		reader := csv.NewReader(fs)
		reader.Comma = '\t'
		reader.FieldsPerRecord = 13
		infunc(reader)
	}
}

func (mw *MyMainWin) doLoadprojs(reader *csv.Reader) {
	mw.doPrompt(mw.txtPathTE.Text())
	pS := projSet{m: map[string]interface{}{}}
	mw.doPrompt(fmt.Sprintln(pS))
	for {
		row, err := reader.Read()
		mw.doPrompt(fmt.Sprintln(row))
		mw.doPrompt(fmt.Sprintln(err))
		if err != nil && err != io.EOF {
			continue
		}
		if err == io.EOF {
			break
		}
		mw.doPrompt(fmt.Sprintln(row))
		pS.Add(row[0])
	}
	mw.doPrompt("循环完了.")
	// 清空列表中原来的数据，并把现在加载到的数据加进去
	mw.projModel.items = pS.ToStrings()
	mw.doPrompt("pS.ToStrings done.")
	mw.projModel.PublishItemsReset()
	mw.doPrompt("Refresh ListBox done.")
	mw.dataPath = mw.txtPathTE.Text()
}

func (mw *MyMainWin) doExportDupData(reader *csv.Reader) {
	// 存储所有数据，并且做聚合。
	data := map[string]([]([]string)){}

	for {
		row, err := reader.Read()
		if err != nil && err != io.EOF {
			continue
		}
		if err == io.EOF {
			break
		}
		if row[0] == mw.proj {
			if len(data[row[1]]) == 0 {
				data[row[1]] = []([]string){row[2:12]}
			} else {
				data[row[1]] = append(data[row[1]], row[2:12])
			}
		}
	}
	fmt.Println(data)
}

func main() {
	//主窗口对象
	mw := &MyMainWin{projModel: &LsBoxModel{items: []string{"-"}}, fontsize: 8}
	icon, _ := walk.Resources.Icon("./resources/main.ico")

	if _, err := (MainWindow{
		Font:        Font{PointSize: mw.fontsize},
		Size:        Size{Width: 400, Height: 300},
		Layout:      VBox{MarginsZero: true},
		Title:       "Eiko的工具箱",
		AssignTo:    &mw.MainWindow,
		Icon:        icon,
		OnDropFiles: mw.doDropFiles,
		Children: []Widget{
			GroupBox{
				Layout: HBox{},
				Title:  "导入数据",
				Font:   Font{PointSize: mw.fontsize},
				Children: []Widget{
					Label{Text: "文件:", Font: Font{PointSize: 10}},
					LineEdit{AssignTo: &mw.txtPathTE},
					PushButton{Text: "浏览", OnClicked: mw.selectFile},
					PushButton{Text: "加载数据", OnClicked: mw.loadK4840Data(mw.txtPathTE.Text(), mw.doLoadprojs)},
				},
			},
			GroupBox{
				Layout: HBox{},
				Title:  "数据处理",
				Font:   Font{PointSize: 10},
				Children: []Widget{
					HSplitter{
						Children: []Widget{
							ListBox{
								AssignTo:              &mw.projLsBox,
								Model:                 mw.projModel,
								OnCurrentIndexChanged: mw.doCurrentIndexChanged,
								OnItemActivated:       mw.doItemActivated,
							},
							PushButton{Text: "导出多次实验的数据", OnClicked: mw.loadK4840Data(mw.txtPathTE.Text(), mw.doExportDupData)},
						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}
