package view

import (
	"pandora/pkg/account"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

// MakeUserForm EscIDとパスワードを入力する部分を作成する
func MakeUserForm(parent fyne.Window) fyne.CanvasObject {
	ecsIDentry := widget.NewEntry()
	ecsIDentry.PlaceHolder = "ecsID"

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.PlaceHolder = "p@ssword"

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "EcsID", Widget: ecsIDentry},
			{Text: "Password", Widget: passwordEntry},
		},
		OnSubmit: func() {
			// 入力された内容をファイルに保存する
			id := ecsIDentry.Text
			password := passwordEntry.Text

			if id != "" && password != "" {
				if err := account.WriteAccountInfo(id, password); err != nil {
					dialog.NewError(err, parent)
				}
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
