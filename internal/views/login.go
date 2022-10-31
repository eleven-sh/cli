package views

type LoginViewDataContent struct {
	Message string
}

type LoginViewData struct {
	Error   *ViewableError
	Content LoginViewDataContent
}

type LoginView struct {
	BaseView
}

func NewLoginView(baseView BaseView) LoginView {
	return LoginView{
		BaseView: baseView,
	}
}

func (l LoginView) View(data LoginViewData) {
	if data.Error == nil {
		l.ShowSuccessView(data.Content.Message, "")
		return
	}

	l.ShowErrorView(data.Error)
}
