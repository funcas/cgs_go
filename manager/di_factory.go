package manager

import "github.com/sarulabs/di/v2"

var app di.Container

func InitDiFactory(container di.Container) {
	app = container
}

func GetDiFactory() di.Container {
	return app
}

func Destroy() {
	app.Delete()
}
