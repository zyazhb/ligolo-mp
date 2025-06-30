package forms

import (
	"github.com/rivo/tview"
)

var (
	traceroute_addr = FormVal[string]{
		Hint: "Address to traceroute\n\nExample:\n1.2.3.4",
	}
)

type TracerouteForm struct {
	tview.Flex
	form      *tview.Form
	submitBtn *tview.Button
	cancelBtn *tview.Button
}

func NewTracerouteForm() *TracerouteForm {
	page := &TracerouteForm{
		Flex:      *tview.NewFlex(),
		form:      tview.NewForm(),
		submitBtn: tview.NewButton("Submit"),
		cancelBtn: tview.NewButton("Cancel"),
	}

	hintBox := tview.NewTextView()
	hintBox.SetTitle("HINT")
	hintBox.SetTitleAlign(tview.AlignCenter)
	hintBox.SetBorder(true)
	hintBox.SetBorderPadding(1, 1, 1, 1)

	page.form.SetTitle("Traceroute").SetTitleAlign(tview.AlignCenter)
	page.form.SetBorder(true)
	page.form.SetButtonsAlign(tview.AlignCenter)

	addrField := tview.NewInputField()
	addrField.SetLabel("Address")
	addrField.SetText(traceroute_addr.Last)
	addrField.SetFocusFunc(func() {
		hintBox.SetText(traceroute_addr.Hint)
	})
	addrField.SetChangedFunc(func(text string) {
		traceroute_addr.Last = text
	})
	page.form.AddFormItem(addrField)

	page.form.AddButton("Submit", nil)
	page.form.AddButton("Cancel", nil)

	formFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(page.form, 11, 1, true).
		AddItem(hintBox, 8, 1, false)

	page.Flex.AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(formFlex, 0, 1, true).
			AddItem(nil, 0, 1, false),
			0, 1, true).
		AddItem(nil, 0, 1, false)

	return page
}

func (page *TracerouteForm) GetID() string {
	return "addroute_page"
}

func (page *TracerouteForm) SetSubmitFunc(f func(string)) {
	btnId := page.form.GetButtonIndex("Submit")
	submitBtn := page.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(func() {
		f(traceroute_addr.Last)
	})
}

func (page *TracerouteForm) SetCancelFunc(f func()) {
	btnId := page.form.GetButtonIndex("Cancel")
	submitBtn := page.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(f)
}
