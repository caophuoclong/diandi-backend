package middlewares

type Middlewares []Middleware

type Middleware interface {
	SetUp()
}

func NewMiddlewares() {
}

func (mw Middlewares) SetUp() {
	for _, middleware := range mw {
		middleware.SetUp()
	}
}
