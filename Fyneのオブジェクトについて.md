# Fyne パッケージについて

## 概要

- App  
  基底部分
- Window  
  ウィンドウ。複数個のウィンドウを表示する場合は別のスレッドで`Window.Show()`を呼ぶ。
- Canvas  
  アプリケーションが描画される場所。全てのWindowはCanvasを持っているが、直接`Window.Canvas()`でアクセスしないのが普通。一つのCanvasObjectを表示するのに用いられる。
- CanvasObject  
  Fyneにおいて描画される物は全てこの型。
- Container  
  CanvasObject型の一種。複数のCanvasObjectを配置するのに用いられる。
- Widget  
  ロジック付きの部品(CanvasObject)

## 主要な部品

### canvasパッケージ

- [Rectangle](https://tour.fyne.io/canvas/rectangle.html)  
  長方形。最も単純なキャンバスオブジェクト。`FillColor`フィールドを用いて色を指定することができる。
- [Text](https://tour.fyne.io/canvas/text.html)  
  テキスト。`Alignment`と`TextStyle`を用いて設定を変えることができる。
- [Line](https://tour.fyne.io/canvas/line.html)  
  線。`Position1`,`Position2`フィールドや`Move()`,`Resize()`メソッドを用いることで線を設定することができる。width:0は垂直な線を、height:0は水平な線を表す。
- [Circle](https://tour.fyne.io/canvas/circle.html)  
  円。使わないので省略。
- [Image](https://tour.fyne.io/canvas/image.html)  
  大きさ可変な画像を表すオブジェクト。
- [Raster](https://tour.fyne.io/canvas/raster.html)  
  1pixelずつ描画できるオブジェクト。
- [Gradient](https://tour.fyne.io/canvas/gradient.html)  
  色の勾配を表現できるオブジェクト。

### layoutパッケージ

- [Box Layout](https://tour.fyne.io/layout/boxlayout.html)  
  最も一般的なレイアウト。垂直に要素を配置するHboxと水平に配置するVboxが存在する。`layout.NewSpacer()`を用いて要素間にスペースをとることができる。
  - Hbox  
    `layout.NewHBoxLayout()`で生成できる。この中に配置される要素の幅は全てその要素の最小幅に設定される。また高さについては、Hbox内の要素がもつ最大の`MinSize().Height`に統一される。`widget.NewHBox()`でも似たようなことができる。
  - Vbox  
    Hboxを横に倒した版。

- [Grid Layout](https://tour.fyne.io/layout/gridlayout.html)  
  格子状に要素を配置するレイアウト。`layout.NewGridLayout(cols)`によって生成され、colの数を満たすまで一列に並べ、それ以上の数は次の行へ送られる。`fyne.NewContainerWithLayout(...)`の第一引数に渡す。コンテナをリサイズすると、全てのセルが利用可能なスペースを平等に分け合うようにリサイズされる。

- [Fixed Grid Layout](https://tour.fyne.io/layout/fixedgridlayout.html)  
  各セルが同じ大きさをもち、ウィンドウの大きさに合わせて自動的に配置が変更されるような格子状のレイアウト。`layout.NewFixedGridLayout(size)`で生成される。

- [Border Layout](https://tour.fyne.io/layout/borderlayout.html)  
  `layout.NewBorderLayout(top, bottom, left, right)`で生成される、上下左右に要素を配置するレイアウト。

- [Form Layout](https://tour.fyne.io/layout/formlayout.html)  
  入力フォームをつくるレイアウト。2列のグリッドと似ているが、横幅を拡張する点は異なる。`layout.NewFormLayout()`で生成される。普通`widget.Form`の中で用いられる。

- [Center Layout](https://tour.fyne.io/layout/centerlayout.html)  
  全ての要素を全て中心に配置するレイアウト。要素は全て最小サイズに設定される。`layout.NewCenterLayout()`で生成される。

- [Max Layout](https://tour.fyne.io/layout/maxlayout.html)  
  全ての要素をコンテナと同じ大きさに設定するレイアウト。`layout.NewMaxLayout()`で生成される。

### widgetパッケージ

- [Label](https://tour.fyne.io/widget/label.html)  
  `widget.NewLabel("some text")`で生成される、フォーマットが可能なテキストオブジェクト。

- [Button](https://tour.fyne.io/widget/button.html)  
  `widget.NewButton()`や`widget.NewButtonWithIcon()`で生成されるボタンオブジェクト。

- [Box](https://tour.fyne.io/widget/box.html)  
  `widget.NewHBox()`や`widget.NewVBox()`で生成されるボックスオブジェクト。

- [Entry](https://tour.fyne.io/widget/entry.html)  
  `widget.NewEntry()`で生成される、文章入力インテーフェースオブジェクト。`NewPasswordEntry()`を用いてパスワード用のエントリーも作ることができる。

- [Choices](https://tour.fyne.io/widget/choices.html)  
  `widget.NewCheck(..)`、`widget.NewRadio(...)`、`widget.NewSelect(...)`で生成される選択肢オブジェクト。

- [Form](https://tour.fyne.io/widget/form.html)  
  入力部分を配置するウィジェット。`widget.NewForm(...)`か`&widget.Form{}`で生成することができる。

- [ProgressBar](https://tour.fyne.io/widget/progressbar.html)  
  `widget.NewProgressBar()`や`widget.NewProgressBarInfinite()`で生成されるプログレスバー。

- [TabContainer](https://tour.fyne.io/widget/tabcontainer.html)  
  様々なパネルを切り替えるために使われるタブ。`widget.NewTabContainer(...)`で生成される。

- [Toolbar](https://tour.fyne.io/widget/toolbar.html)  
  `widget.NewToolbar(...)`で生成されるツールバー。
