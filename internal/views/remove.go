package views

type RemoveViewDataContent struct {
	Message string
}

type RemoveViewData struct {
	Error   *ViewableError
	Content RemoveViewDataContent
}

type RemoveView struct {
	BaseView
}

func NewRemoveView(baseView BaseView) RemoveView {
	return RemoveView{
		BaseView: baseView,
	}
}

func (r RemoveView) View(data RemoveViewData) {
	if data.Error == nil {
		r.ShowSuccessView(data.Content.Message, "")
		return
	}

	r.ShowErrorView(data.Error)
}
