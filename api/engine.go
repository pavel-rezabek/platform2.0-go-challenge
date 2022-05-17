package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateEngine(db *gorm.DB) *gin.Engine {
	engine := gin.Default()
	uc := UserController{db: db, SessionConfig: &gorm.Session{}}
	ac := AssetController{db: db, SessionConfig: &gorm.Session{}}

	// Allow user creation without authorization
	insecure := engine.Group("/api/v1")
	// TODO: token refresh endpoint
	insecure.POST("/users", uc.PostUsers)
	insecure.POST("/token", uc.PostToken)

	// Require JWT token for this group
	secure := engine.Group("/api/v1")
	secure.Use(AuthMiddleware())

	// User methods
	secure.GET("/users", uc.GetUsers) // TODO: size query param
	secure.GET("/users/:id", uc.GetUserByID)
	secure.DELETE("/users/:id", uc.DeleteUserByID)

	// Favourites methods
	secure.GET("/users/:id/favourites", uc.GetFavourites) // TODO: size query param
	secure.POST("/users/:id/favourites", uc.PostFavourites)
	secure.GET("/users/:id/favourites/:favId", uc.GetFavouriteByID)
	secure.DELETE("/users/:id/favourites/:favId", uc.DeleteFavouriteByID)

	// Asset management
	secure.GET("/assets", ac.GetAssets) // TODO: size query param
	secure.POST("/assets", ac.PostAssets)
	secure.GET("/assets/:id", ac.GetAssetByID)
	secure.PUT("/assets/:id", ac.PutAssetByID)
	secure.PATCH("/assets/:id", ac.PatchAssetByID)
	secure.DELETE("/assets/:id", ac.DeleteAssetByID)
	return engine
}
