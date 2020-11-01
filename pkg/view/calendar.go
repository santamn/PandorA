package view

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const (
	// 日付の文字サイズ
	dateSize = 20
	// 曜日の文字サイズ
	daySize = 15
	// セルの下部にあるフォームの一辺の大きさ
	formSize = 40
)

// TODO:科目情報を受け取って科目ページを作る関数と科目ページへ遷移するためのボタンを生成する関数

// カレンダーのセルの中身を生成する
func makeCellContent(date *canvas.Text, isToday bool) *fyne.Container {
	// セル上部の日付が書いてある部分
	top := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), date)
	// しきりとなる線
	line := canvas.NewLine(color.White)
	line.StrokeWidth = 1
	// セル下部の科目名などが入るフォーム
	rec := canvas.NewRectangle(color.Transparent)
	rec.SetMinSize(fyne.NewSize(formSize, formSize))
	form := fyne.NewContainerWithLayout(layout.NewMaxLayout(), rec)

	var back *canvas.Rectangle
	if isToday {
		back = canvas.NewRectangle(color.RGBA{R: 50, G: 50, A: 10})
	} else {
		back = canvas.NewRectangle(color.Transparent)
	}

	return fyne.NewContainerWithLayout(
		layout.NewMaxLayout(),
		back,
		fyne.NewContainerWithLayout(layout.NewVBoxLayout(), top, form, line),
	)
}

// カレンダーのセルを生成する
func makeCell(t time.Time, isToday bool) (cell *fyne.Container) {
	_, _, day := t.Date()
	date := canvas.NewText(fmt.Sprint(day), color.White)
	date.TextSize = dateSize

	switch t.Weekday() {
	case 0:
		// 日曜日の場合は日付を赤に設定
		date.Color = color.RGBA{R: 255, A: 255}
	case 6:
		// 土曜日の場合は日付を青に設定
		date.Color = color.RGBA{B: 255, A: 50}
	}

	return makeCellContent(date, isToday)
}

// 曜日名表示を生成する
func makedWeekDayLabel() (week *fyne.Container) {
	days := make([]*canvas.Text, 0, 7)
	week = fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(7))

	// 各曜日のテキストオブジェクトを作成する
	days = append(days, canvas.NewText("Sun.", color.RGBA{R: 255, A: 255}))
	for _, text := range []string{"Mon.", "Tue.", "Wed.", "Thu.", "Fri."} {
		days = append(days, canvas.NewText(text, color.White))
	}
	days = append(days, canvas.NewText("Sat.", color.RGBA{B: 255, A: 50}))

	for _, day := range days {
		day.Alignment = fyne.TextAlignCenter
		day.TextSize = 15
		week.AddObject(
			fyne.NewContainerWithLayout(layout.NewVBoxLayout(), day, canvas.NewLine(color.White)),
		)
	}

	return
}

// 一月分のセルを格子状に配置したものを生成する
func makeMonthGrid(year int, month time.Month) (monthGrid *fyne.Container) {
	loc, _ := time.LoadLocation("Local")
	// 月の初日
	day := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	// 今日の日付
	now := time.Now()

	// 日曜始まりになるようにカレンダーにパディングを追加
	count := int(day.Weekday())
	padding := make([]fyne.CanvasObject, 0, 6)
	for i := 0; i < count; i++ {
		padding = append(
			padding,
			makeCellContent(&canvas.Text{Color: color.Transparent, Text: "", TextSize: dateSize}, false),
		)
	}

	cells := make([]fyne.CanvasObject, 0, 42)
	cells = append(cells, padding...)
	// ある月の全てのセルを取得する
	for day.Month() == month {
		y, m, d := now.Date()
		isToday := (y == day.Year() && m == day.Month() && d == day.Day())
		cells = append(cells, makeCell(day, isToday))
		day = day.AddDate(0, 0, 1)
	}

	return fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		makedWeekDayLabel(),
		fyne.NewContainerWithLayout(layout.NewGridLayout(7), cells...),
	)
}

// SemesterScreen 半期分のカレンダーを表示する関数
func SemesterScreen(base fyne.Window) fyne.CanvasObject {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	var semester *widget.TabContainer
	switch month {
	case 4, 5, 6, 7, 8, 9:
		// 前期のカレンダーを作成
		semester = widget.NewTabContainer(
			widget.NewTabItem("April", makeMonthGrid(year, time.April)),
			widget.NewTabItem("May", makeMonthGrid(year, time.May)),
			widget.NewTabItem("June", makeMonthGrid(year, time.June)),
			widget.NewTabItem("July", makeMonthGrid(year, time.July)),
			widget.NewTabItem("August", makeMonthGrid(year, time.August)),
			widget.NewTabItem("September", makeMonthGrid(year, time.September)),
		)
		// 初期位置を今月に設定
		semester.SelectTabIndex(month - 4)
	case 10, 11, 12, 1, 2, 3:
		// 後期のカレンダーを作成
		semester = widget.NewTabContainer(
			widget.NewTabItem("October", makeMonthGrid(year, time.October)),
			widget.NewTabItem("November", makeMonthGrid(year, time.November)),
			widget.NewTabItem("December", makeMonthGrid(year, time.December)),
			widget.NewTabItem("January", makeMonthGrid(year, time.January)),
			widget.NewTabItem("February", makeMonthGrid(year, time.February)),
			widget.NewTabItem("March", makeMonthGrid(year, time.March)),
		)

		semester.SelectTabIndex((month + 2) % 12)
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			showUpdateForm(base)
		}),
	)

	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar,
		semester,
	)
}
