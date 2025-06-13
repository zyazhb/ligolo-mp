package route

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/forms"
	"github.com/ttpreport/ligolo-mp/v2/internal/session"
)

type MoveRouteForm struct {
	tview.Flex
	form      *tview.Form
	submitBtn *tview.Button
	cancelBtn *tview.Button
}

var (
	move_target_session_id = forms.FormVal[forms.FormSelectVal]{
		Hint: "???\n???",
	}
)

func NewMoveRouteForm(sessions []*session.Session) *MoveRouteForm {
	form := &MoveRouteForm{
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

	form.form.SetTitle("Move route").SetTitleAlign(tview.AlignCenter)
	form.form.SetBorder(true)
	form.form.SetButtonsAlign(tview.AlignCenter)

	targetField := tview.NewDropDown()
	targetField.SetLabel("Target session")
	targetField.SetFocusFunc(func() {
		hintBox.SetText(move_target_session_id.Hint)
	})
	var sessionNames []string
	sessionMap := make(map[string]string)
	for index, session := range sessions {
		sessionNames = append(sessionNames, session.GetName())
		sessionMap[fmt.Sprintf("%s%d", session.GetName(), index)] = session.ID
	}
	targetField.SetOptions(sessionNames, func(option string, index int) {
		move_target_session_id.Last.ID = index
		move_target_session_id.Last.Value = sessionMap[fmt.Sprintf("%s%d", option, index)]
	})
	targetField.SetCurrentOption(move_target_session_id.Last.ID)
	form.form.AddFormItem(targetField)

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
			0, 2, true).
		AddItem(nil, 0, 1, false)

	return form
}

func (form *MoveRouteForm) GetID() string {
	return "editroute_form"
}

func (form *MoveRouteForm) SetSubmitFunc(f func(string)) {
	btnId := form.form.GetButtonIndex("Submit")
	submitBtn := form.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(func() {
		f(move_target_session_id.Last.Value)
	})
}

func (form *MoveRouteForm) SetCancelFunc(f func()) {
	btnId := form.form.GetButtonIndex("Cancel")
	submitBtn := form.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(f)
}
