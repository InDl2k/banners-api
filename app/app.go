package app

import (
	"banners/controllers"
	"banners/internal/config"
	"banners/internal/utils/cache"
	u "banners/internal/utils/jwt"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type App struct {
	Router      *mux.Router
	BannerCache cache.BannerCache
}

func (a *App) Initialize() {

	cfg := config.MustLoad()

	a.BannerCache = cache.NewRedisCache(cfg.Redis.Address, 0, 5*60)
	controllers.NewBannerController(a.BannerCache)

	a.Router = mux.NewRouter()
	//TOKEN GET
	a.Router.HandleFunc("/token/{username}", controllers.GetToken).Methods("GET")

	//BASE API
	//FOR ADMIN
	adminRouter := a.Router.PathPrefix("").Subrouter()
	adminRouter.Use(u.JWTAuthAdmin)
	adminRouter.HandleFunc("/banner", controllers.CreateBanner).Methods("POST")
	adminRouter.HandleFunc("/banner", controllers.GetBanners).Methods("GET")
	adminRouter.HandleFunc("/banner/{id}", controllers.UpdateBanner).Methods("PATCH")
	adminRouter.HandleFunc("/banner/{id}", controllers.DeleteBanner).Methods("DELETE")
	//FOR USER
	userRouter := a.Router.PathPrefix("").Subrouter()
	userRouter.Use(u.JWTAuthUser)
	userRouter.HandleFunc("/user_banner", controllers.GetUserBanner).Methods("GET")

}

func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
