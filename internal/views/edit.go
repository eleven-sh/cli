package views

type EditViewDataContent struct {
	Message string
}

type EditViewData struct {
	Error   *ViewableError
	Content EditViewDataContent
}

type EditView struct {
	BaseView
}

func NewEditView(baseView BaseView) EditView {
	return EditView{
		BaseView: baseView,
	}
}

func (e EditView) View(data EditViewData) {
	if data.Error == nil {
		e.ShowSuccessView(data.Content.Message, "")
		return
	}

	e.ShowErrorView(data.Error)
}
