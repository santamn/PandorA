package main

import (
	"image/color"
	"pandora/pkg/account"
	"pandora/pkg/resource"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

//フォームを作成する関数
func makeForm(parent fyne.Window) fyne.CanvasObject {
	ecsID, password, rejectable, err := account.ReadAccountInfo()

	ecsIDentry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	videoCheck := widget.NewCheck("Video", func(_ bool) {})
	audioCheck := widget.NewCheck("Audio", func(_ bool) {})
	excelCheck := widget.NewCheck("Excel", func(_ bool) {})
	powerPointCheck := widget.NewCheck("Power Point", func(_ bool) {})
	wordCheck := widget.NewCheck("Word", func(_ bool) {})

	if err == nil {
		ecsIDentry.Text = ecsID
		passwordEntry.Text = password
		videoCheck.Checked = rejectable.Video
		audioCheck.Checked = rejectable.Audio
		excelCheck.Checked = rejectable.Excel
		powerPointCheck.Checked = rejectable.PowerPoint
		wordCheck.Checked = rejectable.Word
	} else {
		ecsIDentry.PlaceHolder = "ecsID"
		passwordEntry.PlaceHolder = "p@ssword"
		videoCheck.Checked = true
		audioCheck.Checked = true
		excelCheck.Checked = true
		powerPointCheck.Checked = true
		wordCheck.Checked = true
	}

	accountFormContainer := fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithColumns(2),
		canvas.NewText("ECS-ID", color.White),
		ecsIDentry,
		canvas.NewText("Password", color.White),
		passwordEntry,
	)

	save := widget.NewButton("Save", func() {
		// 入力された内容をファイルに保存してウィンドウを閉じる
		id := ecsIDentry.Text
		password := passwordEntry.Text

		rejectable := new(resource.RejectableType)
		rejectable.Video = videoCheck.Checked
		rejectable.Audio = audioCheck.Checked
		rejectable.Excel = excelCheck.Checked
		rejectable.PowerPoint = powerPointCheck.Checked
		rejectable.Word = wordCheck.Checked

		if id != "" && password != "" {
			if err := account.WriteAccountInfo(id, password, rejectable); err != nil {
				dialog.NewError(err, parent)
			}
			parent.Close()
		}
	})

	cancel := widget.NewButton("Cancel", func() {
		// 入力された内容を消去する
		ecsIDentry.SetText("")
		passwordEntry.SetText("")
	})

	buttonContainer := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		save,
		cancel,
		layout.NewSpacer(),
	)

	header := canvas.NewText("Enter your account information", color.White)
	header.Alignment = fyne.TextAlignCenter

	middle := canvas.NewText("Select file types not to download", color.White)
	middle.Alignment = fyne.TextAlignCenter

	base := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		header,
		accountFormContainer,
		middle,
		videoCheck,
		audioCheck,
		excelCheck,
		powerPointCheck,
		wordCheck,
		buttonContainer,
	)

	return base
}
