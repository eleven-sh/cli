package views

type UnserveViewDataContent struct {
	Message string
}

type UnserveViewData struct {
	Error   *ViewableError
	Content UnserveViewDataContent
}

type UnserveView struct {
	BaseView
}

func NewUnserveView(baseView BaseView) UnserveView {
	return UnserveView{
		BaseView: baseView,
	}
}

func (u UnserveView) View(data UnserveViewData) {
	if data.Error == nil {
		u.ShowSuccessView(data.Content.Message, "")
		return
	}

	u.ShowErrorView(data.Error)
}
