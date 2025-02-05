package app

import "net/http"

const (
	addr = ":9000"
)

type App struct {
	HttpServer *http.ServeMux
}

func New() App {
	return App{
		HttpServer: http.NewServeMux(),
	}
}
