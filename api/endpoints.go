package api

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/GlobalWebIndex/platform2.0-go-challenge/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: optionally wrap these functions under struct to give them access
//		 to the `db` struct

// POST /users
func PostUsers(c *gin.Context) {
	var apiUser User

	if err := c.ShouldBindJSON(&apiUser); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid input",
			},
		)
		return
	}

	dbUser, err := apiUser.getDBUser()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs",
			},
		)
		return
	}
	session := db.GetDB()
	result := session.Create(&dbUser)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		urlPath := c.Request.URL.Path
		paths := append([]string{urlPath}, fmt.Sprint(dbUser.ID))
		urlPath = path.Join(paths...)
		c.Header("Location", urlPath)
		c.PureJSON(http.StatusCreated, dbUser)
		return
	}
}

// GET /users
func GetUsers(c *gin.Context) {
	var dbUsers []db.User

	session := db.GetDB()
	result := session.Find(&dbUsers)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		c.PureJSON(http.StatusOK, dbUsers)
		return
	}
}

// GET /users/:id
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	var dbUser = db.User{ID: uint(userId)}

	session := db.GetDB()
	// Prevent ErrRecordNotFound
	result := session.Limit(1).Find(&dbUser)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		c.PureJSON(http.StatusOK, dbUser)
		return
	}
}

// DELETE /users/:id
func DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	var dbUser = db.User{ID: uint(userId)}

	session := db.GetDB()
	result := session.Delete(&dbUser)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		c.Status(http.StatusNoContent)
		return
	}
}

// POST /users/:id/favourites
// PostFavourites adds an Asset to favourite Assets of a User
func PostFavourites(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	var dbUser = db.User{ID: uint(userId)}

	var apiFavourite Favourite

	if err := c.ShouldBindJSON(&apiFavourite); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid input",
			},
		)
		return
	}

	dbAsset, err := apiFavourite.getDBAsset()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs",
			},
		)
		return
	}
	session := db.GetDB()

	err = session.Model(&dbUser).Association("Favourites").Append(&dbAsset)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "DB problem.",
			},
		)
		return
	} else {
		urlPath := c.Request.URL.Path
		paths := append([]string{urlPath}, fmt.Sprint(dbAsset.ID))
		urlPath = path.Join(paths...)
		c.Header("Location", urlPath)
		c.PureJSON(http.StatusCreated, dbAsset)
		return
	}
}

// GET /users/:id/favourites
func GetFavourites(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	var dbUser = db.User{ID: uint(userId)}
	var dbAssets []*db.Asset

	session := db.GetDB()
	// Prevent ErrRecordNotFound
	err := session.Preload(clause.Associations).Model(&dbUser).Association("Favourites").Find(&dbAssets)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if len(dbAssets) == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		c.PureJSON(http.StatusOK, dbAssets)
		return
	}
}

// GET /users/:id/favourites/:favId
func GetFavouriteByID(c *gin.Context) {
	id := c.Param("id")
	favId := c.Param("favId")
	userId, _ := strconv.Atoi(id)
	assetId, _ := strconv.Atoi(favId)
	var dbUser = db.User{ID: uint(userId)}
	var dbAsset = db.Asset{ID: uint(assetId)}

	session := db.GetDB()
	// Prevent ErrRecordNotFound
	err := session.Preload(clause.Associations).Model(&dbUser).Association("Favourites").Find(&dbAsset)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "DB problem.",
			},
		)
		return
	} else {
		c.PureJSON(http.StatusOK, dbAsset)
		return
	}
}

// DELETE /users/:id/favourites/:favId
func DeleteFavouriteByID(c *gin.Context) {
	id := c.Param("id")
	favId := c.Param("favId")
	userId, _ := strconv.Atoi(id)
	assetId, _ := strconv.Atoi(favId)
	var dbUser = db.User{ID: uint(userId)}
	var dbAsset = db.Asset{ID: uint(assetId)}

	session := db.GetDB()
	err := session.Model(&dbUser).Association("Favourites").Delete(&dbAsset)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "DB problem.",
			},
		)
		return
	} else {
		c.Status(http.StatusNoContent)
		return
	}
}

// GET /assets
func GetAssets(c *gin.Context) {
	var dbAssets []db.Asset

	session := db.GetDB()
	result := session.Preload(clause.Associations).Preload("Audience.Characteristics").Find(&dbAssets)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		c.PureJSON(http.StatusOK, dbAssets)
		return
	}
}

// POST /assets
func PostAssets(c *gin.Context) {
	var apiAsset Asset

	if err := c.ShouldBindJSON(&apiAsset); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid input",
			},
		)
		return
	}

	dbAsset, err := apiAsset.getDBAsset()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs",
			},
		)
		return
	}
	session := db.GetDB()
	result := session.Create(&dbAsset)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	} else {
		urlPath := c.Request.URL.Path
		paths := append([]string{urlPath}, fmt.Sprint(dbAsset.ID))
		urlPath = path.Join(paths...)
		c.Header("Location", urlPath)
		c.PureJSON(http.StatusCreated, dbAsset)
		return
	}
}

// GET /assets/:id
func GetAssetByID(c *gin.Context) {
	id := c.Param("id")
	assetId, _ := strconv.Atoi(id)
	var dbAsset = db.Asset{ID: uint(assetId)}

	session := db.GetDB()
	// Prevent ErrRecordNotFound
	result := session.Preload(clause.Associations).Preload("Audience.Characteristics").Find(&dbAsset)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		c.PureJSON(http.StatusOK, dbAsset)
		return
	}
}

// PUT /assets/:id
func PutAssetByID(c *gin.Context) {
	id := c.Param("id")
	assetId, _ := strconv.Atoi(id)
	var apiAsset Asset

	if err := c.ShouldBindJSON(&apiAsset); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid input",
			},
		)
		return
	}

	dbAsset, err := apiAsset.getDBAsset()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs",
			},
		)
		return
	}

	session := db.GetDB()
	result := session.Where(db.Asset{ID: uint(assetId)}).FirstOrCreate(&dbAsset)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		urlPath := c.Request.URL.Path
		paths := append([]string{urlPath}, fmt.Sprint(dbAsset.ID))
		urlPath = path.Join(paths...)
		c.Header("Location", urlPath)
		c.PureJSON(http.StatusCreated, dbAsset)
		return
	}
}

// PATCH /assets/:id
func PatchAssetByID(c *gin.Context) {
	id := c.Param("id")
	assetId, _ := strconv.Atoi(id)
	var apiAsset Asset

	if err := c.ShouldBindJSON(&apiAsset); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid input",
			},
		)
		return
	}

	newDbAsset, err := apiAsset.getDBAsset()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "Invalid inputs",
			},
		)
		return
	}

	session := db.GetDB()
	// Get asset from db, modify changed fields, save
	dbAsset := db.Asset{ID: uint(assetId)}
	result := session.Preload(clause.Associations).Preload("Audience.Characteristics").Find(&dbAsset)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	if newDbAsset.Chart != nil {
		if dbAsset.Chart != nil {
			newDbAsset.Chart.ID = dbAsset.Chart.ID
		}
		dbAsset.Chart = newDbAsset.Chart
	}
	if newDbAsset.Insight != nil {
		if dbAsset.Insight != nil {
			newDbAsset.Insight.ID = dbAsset.Insight.ID
		}
		dbAsset.Insight = newDbAsset.Insight
	}
	if newDbAsset.Audience != nil {
		// TODO: when adding characteristics, see if they are not already in DB
		// TODO: when removing characteristics, delete orphaned ones
		if dbAsset.Audience != nil {
			newDbAsset.Audience.ID = dbAsset.Audience.ID
		}
		dbAsset.Audience = newDbAsset.Audience
	}

	result = session.Session(&gorm.Session{FullSaveAssociations: true}).Save(&dbAsset)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		urlPath := c.Request.URL.Path
		paths := append([]string{urlPath}, fmt.Sprint(dbAsset.ID))
		urlPath = path.Join(paths...)
		c.Header("Location", urlPath)
		c.PureJSON(http.StatusCreated, dbAsset)
		return
	}
}

// DELETE /assets/:id
func DeleteAssetByID(c *gin.Context) {
	id := c.Param("id")
	assetId, _ := strconv.Atoi(id)
	dbAsset := db.Asset{ID: uint(assetId)}

	session := db.GetDB()
	result := session.Delete(&dbAsset)
	if result.Error != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":   result.Error.Error(),
				"message": "DB problem.",
			},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	} else {
		c.Status(http.StatusNoContent)
		return
	}
}
