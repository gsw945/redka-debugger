package gui

import (
	"fmt"
	"redka-debugger/redkacore"
	"redka-debugger/utils"

	"github.com/aerogu/tvchooser"
	"github.com/aerogu/tvchooser/tvclang"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TuiDemo() {
	// Set your preferred language (If is enlish you can remove this line because it is the default)
	tvclang.SetTranslations(tvclang.LangEnglish())

	app := tview.NewApplication().EnableMouse(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	text := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	text.SetText("Select a file or directory")

	flex.AddItem(text, 0, 1, false)

	form := tview.NewForm()
	form.AddButton("File", func() {
		path := tvchooser.FileChooser(app, false)
		if path == "" {
			text.SetText("No file selected")
		} else {
			text.SetText(path)
		}
	})
	form.AddButton("Directory", func() {
		path := tvchooser.DirectoryChooser(app, false)
		if path == "" {
			text.SetText("No directory selected")
		} else {
			text.SetText(path)
		}
	})

	flex.SetBackgroundColor(tcell.ColorRed)
	flex.AddItem(form, 4, 0, false)

	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	flex.AddItem(box, 0, 1, false)

	button := tview.NewButton("Quit").SetSelectedFunc(func() {
		app.Stop()
	})
	button.SetBorder(true).SetRect(0, 0, 22, 3)
	button.SetBackgroundColor(tcell.ColorRed)
	flex.AddItem(button, 0, 1, false)

	app.SetRoot(flex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

var logText *tview.TextView
var list *tview.List
var table *tview.Table

func showLog(template string, args ...any) {
	logText.SetText(fmt.Sprintf(template, args...))
}

func onKeySelected(filePath string, key string, keyType int) {
	showLog("Selected: %s - %v", key, redkacore.RedkaTypeName(keyType))
	db := redkacore.LoadDB(filePath)
	defer db.Close()
	table.Clear()
	// header
	isHash := keyType == 4
	cols := []string{"Value"}
	if isHash {
		cols = []string{"Field", "Value"}
	}
	for c := 0; c < len(cols); c++ {
		tableCell := tview.NewTableCell(cols[c]).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter)
		table.SetCell(0, c, tableCell)
	}
	// rows
	rows := [][]string{}
	switch keyType {
	case 0:
		rows = append(rows, []string{"not supported"})
	case 1:
		valueString, err := db.Str().Get(key)
		if err != nil {
			panic(err)
		}
		rows = append(rows, []string{valueString.String()})
	case 2:
		valueList, err := db.List().Range(key, 0, -1)
		if err != nil {
			panic(err)
		}
		for _, v := range valueList {
			rows = append(rows, []string{v.String()})
		}
	case 3:
		valueSet, err := db.Set().Items(key)
		if err != nil {
			panic(err)
		}
		for _, v := range valueSet {
			rows = append(rows, []string{v.String()})
		}
	case 4:
		valueHash, err := db.Hash().Items(key)
		if err != nil {
			panic(err)
		}
		for k, v := range valueHash {
			rows = append(rows, []string{k, v.String()})
		}
	case 5:
		valueZSet, err := db.ZSet().Range(key, 0, -1)
		if err != nil {
			panic(err)
		}
		for _, v := range valueZSet {
			rows = append(rows, []string{v.Elem.String()})
		}
	default:
		rows = append(rows, []string{"not supported"})
	}
	// data-rows
	for r, row := range rows {
		for c, col := range row {
			tableCell := tview.NewTableCell(col).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter)
			table.SetCell(r+1, c, tableCell)
		}
	}
}

func loadDbKeys(filePath string) {
	db := redkacore.LoadDB(filePath)
	defer db.Close()
	pattern := "*"
	keys, err := db.Key().Keys(pattern)
	if err != nil {
		panic(err)
	}
	list.Clear()
	if len(keys) == 0 {
		list.AddItem("No keys", "", 0, nil)
		return
	}
	for i, key := range keys {
		list.AddItem(key.Key, "", 0, func() {
			onKeySelected(filePath, key.Key, int(key.Type))
		})
		if i == 0 {
			onKeySelected(filePath, key.Key, int(key.Type))
		}
	}
	list.SetCurrentItem(0)
}

func buildHeader(app *tview.Application) *tview.Flex {
	// header-file
	fileText := tview.NewInputField().
		SetLabel("DB File: ").
		SetPlaceholder("DB File Path").
		SetFieldWidth(0).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorDeepSkyBlue).
		SetPlaceholderTextColor(tcell.ColorGray)
	// header-button
	pickButton := tview.NewButton("Pick").
		SetBackgroundColorActivated(tcell.ColorBlack).
		SetLabelColorActivated(tcell.ColorBlue)
	loadButton := tview.NewButton("Load").
		SetBackgroundColorActivated(tcell.ColorBlack).
		SetLabelColorActivated(tcell.ColorBlue)
	// header
	flexHeader := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(fileText, 0, 1, false).
		AddItem(pickButton, 6, 1, false).
		AddItem(loadButton, 6, 1, false)

	fileText.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			showLog("DB File: %s", fileText.GetText())
			app.SetFocus(logText)
		}
	})
	pickButton.SetSelectedFunc(func() {
		path := tvchooser.FileChooser(app, false)
		if len(path) > 0 {
			fileText.SetText(path)
			showLog("DB File: %s", fileText.GetText())
			app.SetFocus(logText)
		}
	})
	loadButton.SetSelectedFunc(func() {
		filePath := fileText.GetText()
		if !utils.FileExists(filePath) {
			showLog("File not exists: %s", filePath)
			return
		}
		loadDbKeys(filePath)
	})

	return flexHeader
}

func buildMiddle(_ *tview.Application) *tview.Flex {
	// middle-list
	list = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(true).
		SetHighlightFullLine(true).
		AddItem("Waiting Load Keys", "", 0, nil)
	/*
		list.SetChangedFunc(func(index int, _ string, _ string, _ rune) {
			mt, _ := list.GetItemText(index)
			showLog("Selected: %s", mt)
			count := rune(list.GetItemCount())
			newShortcut := count + 1
			newText := fmt.Sprintf("List item %d", newShortcut)
			list.AddItem(newText, "", 0, nil)
		})
	*/
	list.SetTitle("Keys List").SetBorder(true)
	// middle-sidebar
	flexSidebar := tview.NewFlex().
		SetDirection(tview.FlexRow).
		// AddItem(tview.NewBox().SetBorder(true).SetTitle("Sidebar"), 0, 1, false)
		AddItem(list, 0, 1, false)
	// middle-right
	// valueBox := tview.NewBox().
	// 	SetBorder(true).
	// 	SetTitle("Right")
	table = tview.NewTable().
		SetBorders(true).
		SetBordersColor(tcell.ColorDefault).
		SetSelectable(false, false)
	table.SetTitle("Value Table").SetBorder(true)
	/*
		for r := 0; r < 10; r++ {
			color := tcell.ColorWhite
			if r < 1 {
				color = tcell.ColorYellow
			}
			for c := 0; c < 10; c++ {
				tableCell := tview.NewTableCell(fmt.Sprintf("%d,%d", r, c)).
					SetTextColor(color).
					SetAlign(tview.AlignCenter)
				table.SetCell(r, c, tableCell)
			}
		}
		table.Select(0, 0).
			SetFixed(1, 1).
			SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEscape {
					table.SetSelectable(false, false)
				}
				if key == tcell.KeyEnter {
					table.SetSelectable(true, true)
				}
			})
		table.SetSelectedFunc(func(row int, column int) {
			for i := 0; i < table.GetColumnCount(); i++ {
				table.GetCell(row, i).SetTextColor(tcell.ColorRed)
			}
			table.SetSelectable(false, false)
		})
	*/
	flexRight := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, false)
		// AddItem(valueBox, 0, 1, false)
	// middle
	flexMiddle := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(flexSidebar, 0, 1, false).
		AddItem(flexRight, 0, 1, false)

	return flexMiddle
}

func buildFooter(_ *tview.Application, pages *tview.Pages) *tview.Flex {
	// footer-button
	button := tview.NewButton("Quit").SetSelectedFunc(func() {
		// pages.SwitchToPage("infobox")
		pages.ShowPage("infobox")
	})
	// footer-text
	logText = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	// footer
	flexFooter := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(button, 10, 1, false).
		AddItem(logText, 0, 1, false)

	return flexFooter
}

func TuiDemo2() {
	app := tview.NewApplication().
		EnableMouse(true).
		EnablePaste(true)

	pages := tview.NewPages()

	flexHeader := buildHeader(app)
	flexMiddle := buildMiddle(app)
	flexFooter := buildFooter(app, pages)

	// flex
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(flexHeader, 1, 1, false).
		AddItem(flexMiddle, 0, 1, false).
		AddItem(flexFooter, 1, 1, false)
	flex.SetBackgroundColor(tcell.ColorBlack)

	infobox := tview.NewModal().
		AddButtons([]string{"Quit", "Cancel"}).
		SetText("Do you want to quit?").
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			} else if buttonLabel == "Cancel" {
				// pages.SwitchToPage("main")
				pages.HidePage("infobox")
			}
		})

	pages.AddPage("main", flex, true, true).
		AddPage("infobox", infobox, false, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			// pages.SwitchToPage("infobox")
			pages.ShowPage("infobox")
			return nil
		}
		return event
	})

	app.SetRoot(pages, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
