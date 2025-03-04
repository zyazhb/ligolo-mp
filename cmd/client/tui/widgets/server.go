package widgets

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	"github.com/ttpreport/ligolo-mp/v2/internal/config"
	"github.com/ttpreport/ligolo-mp/v2/internal/operator"
)

type ServerWidget struct {
	*tview.TextView
	operator     *operator.Operator
	serverConfig *config.Config
	fetchConfig  func()
}

func NewServerWidget() *ServerWidget {
	widget := &ServerWidget{
		TextView: tview.NewTextView(),
	}

	widget.SetTitle(fmt.Sprintf("[::b]%s", strings.ToUpper("server")))
	widget.SetBorder(true)
	widget.SetTextAlign(tview.AlignCenter)

	return widget
}

func (widget *ServerWidget) SetData(config *config.Config, oper *operator.Operator) {
	widget.Clear()
	widget.operator = oper
	widget.serverConfig = config
	widget.Refresh()
}

func (widget *ServerWidget) Refresh() {
	if widget.operator != nil {
		access := "operator"
		if widget.operator.IsAdmin {
			access = "admin"
		}

		text := fmt.Sprintf("Operator: %s@%s (%s) | Agent server: %s", widget.operator.Name, widget.operator.Server, access, widget.serverConfig.ListenInterface)
		widget.SetText(text)
	}
}
