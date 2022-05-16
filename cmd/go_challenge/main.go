package main

import (
	"github.com/GlobalWebIndex/platform2.0-go-challenge/api"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/db"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	// TODO: can be changed when deploying to prod
	engine.Use(location.Default())

	// User methods
	engine.POST("/users", api.PostUsers)
	engine.GET("/users", api.GetUsers) // TODO: size query param
	engine.GET("/users/:id", api.GetUserByID)
	engine.DELETE("/users/:id", api.DeleteUserByID)
	// Favourites methods
	engine.POST("/users/:id/favourites", api.PostFavourites)
	engine.GET("/users/:id/favourites", api.GetFavourites) // TODO: size query param
	engine.GET("/users/:id/favourites/:favId", api.GetFavouriteByID)
	engine.DELETE("/users/:id/favourites/:favId", api.DeleteFavouriteByID)
	// Asset management // TODO: size query param
	engine.GET("/assets", api.GetAssets)
	engine.POST("/assets", api.PostAssets)
	engine.GET("/assets/:id", api.GetAssetByID)
	engine.PUT("/assets/:id", api.PutAssetByID)
	engine.PATCH("/assets/:id", api.PatchAssetByID)
	engine.DELETE("/assets/:id", api.DeleteAssetByID)

	db.FillDB()

	engine.Run("localhost:8080")
}

/*
	Users have favourite assets
	Asset:
		chart - title, axes titles, data
		insight - description
		audience - list of characteristics
			gender
			birth country
			age groups (range)
			number hours on social media

	Endpoints:
		Post endpoints:
			Location header
			201 Status
			Body {"id": <id>, <data about the resource>}

		>Auth:
		POST	/login
				authorization token
				refresh token
		POST	/refresh

		>User:
		POST 	/users
		GET		/users
		GET 	/users/:id
		DELETE 	/users/:id
		POST 	/users/:id/favourites
		GET 	/users/:id/favourites (q_par for size)
		GET 	/users/:id/favourites/:fav_id
		DELETE 	/users/:id/favourites/:fav_id
		PATCH 	/users/:id/favourites/:fav_id

		>Asset management:
		GET   	/assets
		POST	/assets
		GET		/assets/:id
		PUT		/assets/:id
		PATCH	/assets/:id
		DELETE	/assets/:id

	[*] Sqlite db
	Security (user authorization so they cannot access other user's data)
	Tests
	Dockerfile
*/
