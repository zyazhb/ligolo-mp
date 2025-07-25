package pages

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	forms "github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/forms"
	modals "github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/modals"
	widgets "github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/widgets"
	"github.com/ttpreport/ligolo-mp/v2/internal/certificate"
	"github.com/ttpreport/ligolo-mp/v2/internal/config"
	"github.com/ttpreport/ligolo-mp/v2/internal/operator"
)

type AdminPage struct {
	tview.Pages

	flex      *tview.Flex
	server    *widgets.ServerWidget
	operators *widgets.OperatorsWidget
	certs     *widgets.CertificatesWidget

	setFocus func(tview.Primitive)

	getMetadata     func() (*config.Config, *operator.Operator, error)
	getOperators    func() ([]*operator.Operator, error)
	getCertificates func() ([]*certificate.Certificate, error)
	switchback      func()

	exportOperator  func(string, string) (string, error)
	addOperator     func(string, bool, string) (*operator.Operator, error)
	delOperator     func(string) error
	promoteOperator func(string) error
	demoteOperator  func(string) error
	regenCert       func(string) error

	operator *operator.Operator
}

func NewAdminPage() *AdminPage {
	admin := &AdminPage{
		Pages: *tview.NewPages(),

		flex:      tview.NewFlex(),
		server:    widgets.NewServerWidget(),
		operators: widgets.NewOperatorsWidget(),
		certs:     widgets.NewCertificatesWidget(),
	}

	admin.initOperatorsWidget()
	admin.initCertsWidget()

	firstRow := tview.NewFlex()
	firstRow.SetDirection(tview.FlexColumn)
	firstRow.AddItem(admin.operators, 0, 50, true)
	firstRow.AddItem(admin.certs, 0, 50, false)

	admin.flex.SetDirection(tview.FlexRow)
	admin.flex.AddItem(admin.server, 3, 0, false)
	admin.flex.AddItem(firstRow, 0, 100, true)

	admin.Reset()

	return admin
}

func (admin *AdminPage) Reset() {
	for _, page := range admin.GetPageNames(false) {
		admin.RemovePage(page)
	}

	admin.AddAndSwitchToPage("main", admin.flex, true)
}

func (admin *AdminPage) initOperatorsWidget() {
	admin.operators.SetSelectedFunc(func(elem *widgets.OperatorsWidgetElem) {
		menu := modals.NewMenuModal(fmt.Sprintf("Operator — %s", elem.Operator.Name))
		cleanup := func() {
			admin.RemovePage(menu.GetID())
			admin.setFocus(admin.operators)
			admin.RefreshData()
		}

		menu.AddItem(modals.NewMenuModalElem("Export", func() {
			export := forms.NewExportForm()
			export.SetSubmitFunc(func(path string) {
				admin.DoWithLoader("Exporting operator...", func() {
					fullPath, err := admin.exportOperator(elem.Operator.Name, path)
					if err != nil {
						admin.ShowError(fmt.Sprintf("Could not export operator: %s", err), nil)
						return
					}

					admin.RemovePage(export.GetID())
					admin.ShowInfo(fmt.Sprintf("Exported operator to %s", fullPath), cleanup)
					admin.RefreshData()
				})
			})
			export.SetCancelFunc(func() {
				admin.RemovePage(export.GetID())
			})
			admin.AddPage(export.GetID(), export, true, true)
		}))

		if !elem.Operator.IsAdmin {
			menu.AddItem(modals.NewMenuModalElem("Promote", func() {
				admin.DoWithLoader("Promoting operator...", func() {
					err := admin.promoteOperator(elem.Operator.Name)
					if err != nil {
						admin.ShowError(fmt.Sprintf("Could not promote operator: %s", err), nil)
						return
					}

					admin.ShowInfo("Operator promoted", cleanup)
				})
			}))
		} else {
			menu.AddItem(modals.NewMenuModalElem("Demote", func() {
				admin.DoWithLoader("Demoting operator...", func() {
					err := admin.demoteOperator(elem.Operator.Name)
					if err != nil {
						admin.ShowError(fmt.Sprintf("Could not demote operator: %s", err), nil)
						return
					}

					admin.ShowInfo("Operator demoted", cleanup)
				})
			}))
		}

		menu.AddItem(modals.NewMenuModalElem("Remove", func() {
			admin.DoWithLoader("Removing operator...", func() {
				err := admin.delOperator(elem.Operator.Name)
				if err != nil {
					admin.ShowError(fmt.Sprintf("Could not remove operator: %s", err), nil)
					return
				}

				admin.ShowInfo("Operator removed", cleanup)
			})
		}))

		menu.SetCancelFunc(cleanup)

		admin.AddPage(menu.GetID(), menu, true, true)
	})
}

func (admin *AdminPage) initCertsWidget() {
	admin.certs.SetSelectedFunc(func(elem *widgets.CertificatesWidgetElem) {
		menu := modals.NewMenuModal(fmt.Sprintf("Certificate — %s", elem.Certificate.Name))
		cleanup := func() {
			admin.RemovePage(menu.GetID())
			admin.setFocus(admin.certs)
			admin.RefreshData()
		}

		menu.AddItem(modals.NewMenuModalElem("Regenerate", func() {
			name := elem.Certificate.Name
			admin.DoWithConfirm(fmt.Sprintf("Regenerate certificate %s?", name), func() {
				admin.DoWithLoader("Regenerating certificate...", func() {
					err := admin.regenCert(name)
					if err != nil {
						admin.ShowError(fmt.Sprintf("Could not regenerate certificate: %s", err), nil)
						return
					}

					admin.ShowInfo("Certificate regenerated", cleanup)
					admin.RefreshData()
				})
			})
		}))

		menu.SetCancelFunc(cleanup)

		admin.AddPage(menu.GetID(), menu, true, true)
	})
}

func (admin *AdminPage) GetID() string {
	return "admin"
}

func (admin *AdminPage) GetNavBar() []widgets.NavBarElem {
	return []widgets.NavBarElem{
		widgets.NewNavBarElem(tcell.KeyCtrlA, "Back"),
		widgets.NewNavBarElem(tcell.KeyCtrlN, "New operator"),
	}
}

func (admin *AdminPage) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()

		if admin.flex.HasFocus() {
			switch key {
			case tcell.KeyTab:
				focusOrder := []tview.Primitive{
					admin.operators,
					admin.certs,
				}

				for id, pane := range focusOrder {
					if pane.HasFocus() {
						nextId := (id + 1) % len(focusOrder)
						setFocus(focusOrder[nextId])
						break
					}
				}
			case tcell.KeyCtrlA:
				admin.switchback()
			case tcell.KeyCtrlN:
				gen := forms.NewOperatorForm()
				gen.SetSubmitFunc(func(name string, isAdmin bool, server string) {
					admin.DoWithLoader("Creating operator...", func() {
						oper, err := admin.addOperator(name, isAdmin, server)
						if err != nil {
							admin.ShowError(fmt.Sprintf("Could not create operator: %s", err), nil)
							return
						}

						admin.RemovePage(gen.GetID())
						admin.ShowInfo(fmt.Sprintf("Created operator %s", oper.Name), nil)
						admin.RefreshData()
					})
				})
				gen.SetCancelFunc(func() {
					admin.RemovePage(gen.GetID())
				})
				admin.AddPage(gen.GetID(), gen, true, true)
			default:
				defaultHandler := admin.Pages.InputHandler()
				defaultHandler(event, setFocus)
			}
		} else {
			switch key {
			case tcell.KeyEscape:
				if admin.GetPageCount() > 1 {
					frontPage, _ := admin.GetFrontPage()
					admin.RemovePage(frontPage)
				}
			default:
				defaultHandler := admin.Pages.InputHandler()
				defaultHandler(event, setFocus)
			}
		}

	}
}

func (admin *AdminPage) Focus(delegate func(p tview.Primitive)) {
	admin.setFocus = delegate
	admin.Pages.Focus(delegate)
}

func (admin *AdminPage) RefreshData() {
	if !admin.operator.IsAdmin {
		return
	}

	opers, err := admin.getOperators()
	if err != nil {
		admin.ShowError(fmt.Sprintf("Could not refresh operators: %s", err), nil)
		return
	}
	admin.operators.SetData(opers)

	certs, err := admin.getCertificates()
	if err != nil {
		admin.ShowError(fmt.Sprintf("Could not refresh certs: %s", err), nil)
		return
	}
	admin.certs.SetData(certs)

	config, operator, err := admin.getMetadata()
	if err != nil {
		admin.ShowError(fmt.Sprintf("Could not fetch metadata: %s", err), nil)
		return
	}

	admin.server.SetData(config, operator)
}

func (admin *AdminPage) SetMetadataFunc(f func() (*config.Config, *operator.Operator, error)) {
	admin.getMetadata = f
}

func (admin *AdminPage) SetOperator(oper *operator.Operator) {
	admin.operator = oper
}

func (admin *AdminPage) SetExportOperatorFunc(f func(string, string) (string, error)) {
	admin.exportOperator = f
}

func (admin *AdminPage) SetAddOperatorFunc(f func(string, bool, string) (*operator.Operator, error)) {
	admin.addOperator = f
}

func (admin *AdminPage) SetDelOperatorFunc(f func(string) error) {
	admin.delOperator = f
}

func (admin *AdminPage) SetPromoteOperatorFunc(f func(string) error) {
	admin.promoteOperator = f
}

func (admin *AdminPage) SetDemoteOperatorFunc(f func(string) error) {
	admin.demoteOperator = f
}

func (admin *AdminPage) SetRegenCertFunc(f func(string) error) {
	admin.regenCert = f
}

func (admin *AdminPage) SetSwitchbackFunc(f func()) {
	admin.switchback = f
}

func (admin *AdminPage) SetOperatorsFunc(f func() ([]*operator.Operator, error)) {
	admin.getOperators = f
}

func (admin *AdminPage) SetCertificatesFunc(f func() ([]*certificate.Certificate, error)) {
	admin.getCertificates = f
}

func (admin *AdminPage) ShowError(text string, done func()) {
	modal := modals.NewErrorModal()
	modal.SetText(text)
	modal.SetDoneFunc(func(_ int, _ string) {
		admin.RemovePage(modal.GetID())

		if done != nil {
			done()
		}
	})
	admin.AddPage(modal.GetID(), modal, true, true)
}

func (admin *AdminPage) ShowInfo(text string, done func()) {
	modal := modals.NewInfoModal()
	modal.SetText(text)
	modal.SetDoneFunc(func(_ int, _ string) {
		admin.RemovePage(modal.GetID())

		if done != nil {
			done()
		}
	})
	admin.AddPage(modal.GetID(), modal, true, true)
}

func (admin *AdminPage) DoWithLoader(text string, action func()) {
	go func() {
		modal := modals.NewLoaderModal()
		modal.SetText(text)
		admin.AddPage(modal.GetID(), modal, true, true)
		action()
		admin.RemovePage(modal.GetID())
	}()
}

func (admin *AdminPage) DoWithConfirm(text string, action func()) {
	modal := modals.NewConfirmModal()
	modal.SetText(text)
	modal.SetDoneFunc(func(confirmed bool) {
		if confirmed {
			action()
		}

		admin.RemovePage(modal.GetID())
	})
	admin.AddPage(modal.GetID(), modal, true, true)
}
