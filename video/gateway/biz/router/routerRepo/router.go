package routerRepo

import "github.com/cloudwego/hertz/pkg/app/server"

type Router interface {
	Router(r *server.Hertz)
}

var routers []Router

func Register(rs ...Router) {
	routers = append(routers, rs...)
}
func InitRouters(r *server.Hertz) {
	for _, router := range routers {
		router.Router(r)
	}
}
