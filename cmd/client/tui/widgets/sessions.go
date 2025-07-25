package widgets

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/style"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/utils"
	"github.com/ttpreport/ligolo-mp/v2/internal/session"
)

type SessionsWidget struct {
	tview.Table
	data            []*SessionsWidgetElem
	selectedFunc    func(*session.Session)
	selectedSession *session.Session
}

func NewSessionsWidget() *SessionsWidget {
	widget := &SessionsWidget{
		Table: *tview.NewTable(),
		data:  nil,
	}

	widget.Table.SetSelectable(true, false)
	widget.Table.SetBackgroundColor(style.BgColor)
	widget.Table.SetTitle(fmt.Sprintf("[::b]%s", strings.ToUpper("sessions")))
	widget.Table.SetBorderColor(style.BorderColor)
	widget.Table.SetTitleColor(style.FgColor)
	widget.Table.SetBorder(true)

	widget.SetFocusFunc(func() {
		widget.SetSelectable(true, false)
		widget.ResetSelector()
	})
	widget.SetBlurFunc(func() {
		widget.SetSelectable(false, false)
	})

	return widget
}

func (widget *SessionsWidget) SetSelectedSession(sess *session.Session) {
	widget.selectedSession = sess
	widget.Refresh()
}

func (widget *SessionsWidget) FetchSession(row int) *SessionsWidgetElem {
	id := max(0, row-1)
	if len(widget.data) > id {
		return widget.data[id]
	}

	return nil
}

func (widget *SessionsWidget) FetchRow(sess *session.Session) int {
	for row, elem := range widget.data {
		if elem.Session.ID == sess.ID {
			return row + 1
		}
	}

	return 0
}

func (widget *SessionsWidget) SetSelectionChangedFunc(f func(*session.Session)) {
	widget.Table.SetSelectionChangedFunc(func(row, _ int) {
		item := widget.FetchSession(row)
		if item != nil {
			f(item.Session)
		}
	})
}

func (widget *SessionsWidget) SetSelectedFunc(f func(*session.Session)) {
	widget.Table.SetSelectedFunc(func(row, _ int) {
		item := widget.FetchSession(row)
		if item != nil {
			f(item.Session)
		}
	})
}

func (widget *SessionsWidget) SetData(data []*session.Session) {
	widget.Clear()

	widget.data = nil
	for _, session := range data {
		widget.data = append(widget.data, NewSessionsWidgetElem(session))
	}

	widget.Refresh()
	widget.ResetSelector()
}

func (widget *SessionsWidget) ResetSelector() {
	if len(widget.data) > 0 {
		row := 1
		if widget.selectedSession != nil {
			row = widget.FetchRow(widget.selectedSession)
		}

		if row > 0 {
			widget.Select(row, 0)
		}
	}
}

func (widget *SessionsWidget) Refresh() {
	headers := []string{"Alias", "Hostname", "Connected", "Relaying", "First Seen", "Last Seen", ""}
	for colNo, header := range headers {
		header := fmt.Sprintf("[::b]%s", strings.ToUpper(header))
		widget.SetCell(0, colNo, tview.NewTableCell(header).SetExpansion(1).SetSelectable(false)).SetFixed(1, 0)
	}

	rowId := 1
	for _, elem := range widget.data {
		if elem.IsSelected(widget.selectedSession) {
			elem.Highlight(true)
		} else {
			elem.Highlight(false)
		}

		widget.SetCell(rowId, 0, elem.Alias())
		widget.SetCell(rowId, 1, elem.Hostname())
		widget.SetCell(rowId, 2, elem.IsConnected())
		widget.SetCell(rowId, 3, elem.IsRelaying())
		widget.SetCell(rowId, 4, elem.FirstSeen())
		widget.SetCell(rowId, 5, elem.LastSeen())
		widget.SetCell(rowId, 6, elem.Status().SetSelectable(false).SetAlign(tview.AlignCenter))

		rowId++
	}
}

type SessionsWidgetElem struct {
	Session *session.Session
	bgcolor tcell.Color
}

func NewSessionsWidgetElem(session *session.Session) *SessionsWidgetElem {
	return &SessionsWidgetElem{
		Session: session,
		bgcolor: style.BgColor,
	}
}

func (elem *SessionsWidgetElem) IsSelected(sess *session.Session) bool {
	if sess == nil {
		return false
	}

	if sess.ID != elem.Session.ID {
		return false
	}

	return true
}

func (elem *SessionsWidgetElem) Highlight(h bool) {
	if h {
		elem.bgcolor = style.HighlightColor
	} else {
		elem.bgcolor = style.BgColor
	}
}

func (elem *SessionsWidgetElem) Alias() *tview.TableCell {
	val := elem.Session.Alias
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

func (elem *SessionsWidgetElem) Hostname() *tview.TableCell {
	val := elem.Session.Hostname
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

func (elem *SessionsWidgetElem) IsConnected() *tview.TableCell {
	val := utils.HumanBool(elem.Session.IsConnected)
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

func (elem *SessionsWidgetElem) IsRelaying() *tview.TableCell {
	val := utils.HumanBool(elem.Session.IsRelaying)
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

func (elem *SessionsWidgetElem) FirstSeen() *tview.TableCell {
	val := utils.HumanTime(elem.Session.FirstSeen)
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

func (elem *SessionsWidgetElem) LastSeen() *tview.TableCell {
	val := utils.HumanTimeSince(elem.Session.LastSeen)
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

func (elem *SessionsWidgetElem) Status() *tview.TableCell {
	val := "⚑"
	if !elem.Session.IsConnected {
		return tview.NewTableCell(val).SetTextColor(tcell.ColorRed)
	}

	if !elem.Session.IsRelaying {
		return tview.NewTableCell(val).SetTextColor(tcell.ColorBlue)
	}

	return tview.NewTableCell(val).SetTextColor(tcell.ColorGreen)
}
