package modals

import (
	"fmt"

	"github.com/rivo/tview"
)

type TextModal struct {
	tview.Flex
	title string
	form  *tview.Form
}

func NewTextModal(title string, body string) *TextModal {
	page := &TextModal{
		Flex:  *tview.NewFlex(),
		title: title,
		form:  tview.NewForm(),
	}

	page.form.SetTitle(title).SetTitleAlign(tview.AlignCenter)
	page.form.SetBorder(true)
	page.form.SetButtonsAlign(tview.AlignCenter)
	page.form.AddButton("OK", nil)
	page.form.AddTextView("", body, 0, 0, true, true)

	page.AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(page.form, 12, 1, true).
			AddItem(nil, 0, 1, false),
			0, 1, true).
		AddItem(nil, 0, 1, false)

	return page
}

func (page *TextModal) GetID() string {
	return fmt.Sprintf("text_%s", page.title)
}

func (page *TextModal) SetDoneFunc(f func()) {
	btnId := page.form.GetButtonIndex("OK")
	submitBtn := page.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(func() {
		if f != nil {
			f()
		}
	})
}
