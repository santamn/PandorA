package screens

import (
	"fmt"
	"image/color"
	"pandora/pkg/account"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// ecsIDとパスワードを入力するフォームを生成
func showUpdateForm(base fyne.Window) {
	// ページ上部に表示するテキスト
	top := canvas.NewText(
		"Enter your ECS-ID and Password. If blank, the values set previously will be used.",
		color.White,
	)
	top.Alignment = fyne.TextAlignCenter

	// 入力フォームと説明文
	ecsID := widget.NewEntry()
	ecsID.PlaceHolder = "ECS-ID"
	password := widget.NewPasswordEntry()
	password.PlaceHolder = "Password"

	content := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		top,
		ecsID,
		password,
	)

	dialog.ShowCustomConfirm(
		"Update Information",
		"Update",
		"Cancel",
		content,
		func(ok bool) {
			if ok {
				update(ecsID, password, base)
			}
		},
		base,
	)
}

// updateボタンを押した時に作動する関数
func update(ecsID, password *widget.Entry, base fyne.Window) {
	if ecsID.Text != "" || password.Text != "" {
		if err := account.WriteAccountInfo(ecsID.Text, password.Text); err != nil {
			dialog.ShowError(err, base)
		}
	}

	id, pass, err := account.ReadAccountInfo()
	if err != nil {
		dialog.ShowError(err, base)
	}

	// TODO:idとpassを使ってpandaから課題を取得する
	fmt.Println(id, pass)
}
