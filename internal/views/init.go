package views

type InitViewData struct {
	Error   *ViewableError
	Content InitViewDataContent
}

type InitViewDataContent struct {
	ShowAsWarning bool
	Message       string
	Subtext       string
}

type InitView struct {
	BaseView
}

func NewInitView(baseView BaseView) InitView {
	return InitView{
		BaseView: baseView,
	}
}

func (i InitView) View(data InitViewData) {
	if data.Error == nil {
		if data.Content.ShowAsWarning {
			i.ShowWarningView(
				data.Content.Message,
				data.Content.Subtext,
			)
			return
		}

		i.ShowSuccessView(
			data.Content.Message,
			data.Content.Subtext,
		)
		return
	}

	i.ShowErrorView(data.Error)
}
