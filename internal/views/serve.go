package views

type ServeViewDataContent struct {
	Message string
}

type ServeViewData struct {
	Error   *ViewableError
	Content ServeViewDataContent
}

type ServeView struct {
	BaseView
}

func NewServeView(baseView BaseView) ServeView {
	return ServeView{
		BaseView: baseView,
	}
}

func (s ServeView) View(data ServeViewData) {
	if data.Error == nil {
		s.ShowSuccessView(data.Content.Message, "")
		return
	}

	s.ShowErrorView(data.Error)
}
