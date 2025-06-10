package route

import (
	"github.com/rivo/tview"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/forms"
	"github.com/ttpreport/ligolo-mp/v2/internal/route"
)

type EditRouteForm struct {
	tview.Flex
	form      *tview.Form
	submitBtn *tview.Button
	cancelBtn *tview.Button
}

var (
	edit_route_cidr = forms.FormVal[string]{
		Hint: "A CIDR that will be routed via this session.\n\nExample:\n10.10.5.0/24",
	}

	edit_route_loopback = forms.FormVal[bool]{
		Hint: "If checked, specified CIDR will address the machine running the agent itself, i.e. localhost. Use this instead of port forwarding.",
	}
)

func NewEditRouteForm(route *route.Route) *EditRouteForm {
	edit_route_cidr.Last = route.Cidr.String()
	edit_route_loopback.Last = route.IsLoopback

	form := &EditRouteForm{
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

	form.form.SetTitle("Add route").SetTitleAlign(tview.AlignCenter)
	form.form.SetBorder(true)
	form.form.SetButtonsAlign(tview.AlignCenter)

	cidrField := tview.NewInputField()
	cidrField.SetLabel("CIDR")
	cidrField.SetText(edit_route_cidr.Last)
	cidrField.SetFocusFunc(func() {
		hintBox.SetText(edit_route_cidr.Hint)
	})
	cidrField.SetChangedFunc(func(text string) {
		edit_route_cidr.Last = text
	})
	form.form.AddFormItem(cidrField)

	loopbackField := tview.NewCheckbox()
	loopbackField.SetLabel("Loopback")
	loopbackField.SetChecked(edit_route_loopback.Last)
	loopbackField.SetFocusFunc(func() {
		hintBox.SetText(edit_route_loopback.Hint)
	})
	loopbackField.SetChangedFunc(func(checked bool) {
		edit_route_loopback.Last = checked
	})
	loopbackField.SetBlurFunc(func() {
		hintBox.Clear()
	})
	form.form.AddFormItem(loopbackField)

	form.form.AddButton("Submit", nil)
	form.form.AddButton("Cancel", nil)

	formFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form.form, 9, 1, true).
		AddItem(hintBox, 8, 1, false)

	form.Flex.AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(formFlex, 0, 1, true).
			AddItem(nil, 0, 1, false),
			0, 1, true).
		AddItem(nil, 0, 1, false)

	return form
}

func (form *EditRouteForm) GetID() string {
	return "editroute_form"
}

func (form *EditRouteForm) SetSubmitFunc(f func(string, bool)) {
	btnId := form.form.GetButtonIndex("Submit")
	submitBtn := form.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(func() {
		f(edit_route_cidr.Last, edit_route_loopback.Last)
	})
}

func (form *EditRouteForm) SetCancelFunc(f func()) {
	btnId := form.form.GetButtonIndex("Cancel")
	submitBtn := form.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(f)
}
