package views

type UninstallViewDataContent struct {
	ShowAsWarning bool
	Message       string
	Subtext       string
}

type UninstallViewData struct {
	Error   *ViewableError
	Content UninstallViewDataContent
}

type UninstallView struct {
	BaseView
}

func NewUninstallView(baseView BaseView) UninstallView {
	return UninstallView{
		BaseView: baseView,
	}
}

func (u UninstallView) View(data UninstallViewData) {
	if data.Error == nil {
		if data.Content.ShowAsWarning {
			u.ShowWarningView(
				data.Content.Message,
				data.Content.Subtext,
			)
			return
		}

		u.ShowSuccessView(
			data.Content.Message,
			data.Content.Subtext,
		)
		return
	}

	u.ShowErrorView(data.Error)
}
