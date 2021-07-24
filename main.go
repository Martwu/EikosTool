package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
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
	item := mw.projModel.items[i]
	fmt.Println("CurrentIndex: ", i)
	fmt.Println("CurrentEnvVarName: ", item)
	mw.proj = item
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

func (mw *MyMainWin) loadK4840Data(infunc func(*csv.Reader)) func() {
	f := func() {
		fs, err := os.OpenFile(strings.TrimSpace(mw.txtPathTE.Text()), os.O_RDONLY, 4)
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
	return f
}

func (mw *MyMainWin) doLoadprojs(reader *csv.Reader) {
	pS := projSet{m: map[string]interface{}{}}
	for {
		row, err := reader.Read()
		if err != nil && err != io.EOF {
			continue
		}
		if err == io.EOF {
			break
		}
		pS.Add(row[0])
	}
	// 清空列表中原来的数据，并把现在加载到的数据加进去
	mw.projModel.items = pS.ToStrings()
	mw.projModel.PublishItemsReset()
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
		// row[0]为数据所属的项目
		if row[0] == mw.proj {
			// row[1]为实验对象
			if len(data[row[1]]) == 0 {
				data[row[1]] = []([]string){row[2:12]}
			} else {
				data[row[1]] = append(data[row[1]], row[2:12])
			}
		}
	}
	xlsfile := excelize.NewFile()
	xlsfile.SetActiveSheet(xlsfile.NewSheet(mw.proj))
	posLine := 0
	POSCOL := 'A'
	for item := range data {
		// 不重复的实验数据不展示
		if len(data[item]) <= 1 {
			continue
		}
		// 重复的数据逐行展示
		// 遍历同一个项目中的若干行数据
		for eachline := range data[item] {
			// 从第一行开始
			posLine += 1
			// 遍历一行的若干列
			pos := fmt.Sprintf("A%d", posLine)
			fmt.Println(pos, " -- ", item)
			xlsfile.SetCellValue(mw.proj, pos, item)
			for i := range data[item][eachline] {
				// 每一行逐个元素
				pos = fmt.Sprintf("%s%d", string(POSCOL+rune(i+1)), posLine)
				fmt.Println(pos, " -- ", data[item][eachline][i])
				xlsfile.SetCellValue(mw.proj, pos, data[item][eachline][i])
			}
		}
	}
	if err := xlsfile.SaveAs(`C:\Users\blackfat\Downloads\Book1.xlsx`); err != nil {
		fmt.Println(err)
	}
	mw.doPrompt("excel数据导出来啦！")

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
					PushButton{Text: "加载数据", OnClicked: mw.loadK4840Data(mw.doLoadprojs)},
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
							PushButton{Text: "导出多次实验的数据", OnClicked: mw.loadK4840Data(mw.doExportDupData)},
						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet2")
	// Set value of a cell.
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs(`C:\Users\blackfat\Book1.xlsx`); err != nil {
		fmt.Println(err)
	}
}
