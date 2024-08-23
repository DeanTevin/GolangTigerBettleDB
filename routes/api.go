package routes

import (
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	TigerBettle "goravel/app/http/controllers/TigerBettleController"
)

func Api() {
	userController := controllers.NewUserController()
	facades.Route().Get("/users/{id}", userController.Show)

	TigerBettleController := TigerBettle.NewTigerBettleController()
	facades.Route().Post("/test/create-tb-user", TigerBettleController.PostCreateUserTB)
	facades.Route().Post("/test/query-user", TigerBettleController.QueryUserTB)
	facades.Route().Get("/test/account-balance", TigerBettleController.AccountBalances)
}
