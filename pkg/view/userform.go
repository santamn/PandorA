package view

import (
	"image/color"
	"pandora/pkg/account"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// MakeUserForm EscIDとパスワードを入力する部分を作成する
func MakeUserForm(parent fyne.Window) fyne.CanvasObject {
	ecsIDentry := widget.NewEntry()
	ecsIDentry.PlaceHolder = "ecsID"

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.PlaceHolder = "p@ssword"

	videoCheck := widget.NewCheck("Video", func(_ bool) {})
	audioCheck := widget.NewCheck("Audio", func(_ bool) {})
	excelCheck := widget.NewCheck("Excel", func(_ bool) {})
	powerPointCheck := widget.NewCheck("Power Point", func(_ bool) {})
	wordCheck := widget.NewCheck("Word", func(_ bool) {})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "EcsID", Widget: ecsIDentry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "", Widget: videoCheck},
			{Text: "", Widget: audioCheck},
			{Text: "", Widget: excelCheck},
			{Text: "", Widget: powerPointCheck},
			{Text: "", Widget: wordCheck},
		},
		OnSubmit: func() {
			// 入力された内容をファイルに保存してウィンドウを閉じる
			id := ecsIDentry.Text
			password := passwordEntry.Text

			if id != "" && password != "" {
				if err := account.WriteAccountInfo(id, password); err != nil {
					dialog.NewError(err, parent)
				}
				parent.Close()
			}
		},
		OnCancel: func() {
			// 入力された内容を消去する
			ecsIDentry.SetText("")
			passwordEntry.SetText("")
		},
		SubmitText: "Save",
	}

	return form
}

// MakeForm フォームを作成する関数
func MakeForm(parent fyne.Window) fyne.CanvasObject {
	ecsIDentry := widget.NewEntry()
	ecsIDentry.PlaceHolder = "ecsID"
	ecsContainer := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		canvas.NewText("ECS-ID", color.White),
		ecsIDentry,
	)

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.PlaceHolder = "p@ssword"
	passwordContainer := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		canvas.NewText("Password", color.White),
		passwordEntry,
	)

	videoCheck := widget.NewCheck("Video", func(_ bool) {})
	videoCheck.Checked = true
	audioCheck := widget.NewCheck("Audio", func(_ bool) {})
	audioCheck.Checked = true
	excelCheck := widget.NewCheck("Excel", func(_ bool) {})
	excelCheck.Checked = true
	powerPointCheck := widget.NewCheck("Power Point", func(_ bool) {})
	powerPointCheck.Checked = true
	wordCheck := widget.NewCheck("Word", func(_ bool) {})
	wordCheck.Checked = true

	save := widget.NewButton("Save", func() {
		// 入力された内容をファイルに保存してウィンドウを閉じる
		id := ecsIDentry.Text
		password := passwordEntry.Text

		if id != "" && password != "" {
			if err := account.WriteAccountInfo(id, password); err != nil {
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
		save,
		cancel,
	)

	header := canvas.NewText("Enter ypur information", color.White)
	header.Alignment = fyne.TextAlignCenter

	middle := canvas.NewText("Select file types not to download", color.White)
	middle.Alignment = fyne.TextAlignLeading

	base := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		ecsContainer,
		passwordContainer,
		videoCheck,
		audioCheck,
		excelCheck,
		powerPointCheck,
		wordCheck,
		buttonContainer,
	)

	return base
}
