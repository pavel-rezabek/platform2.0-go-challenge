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

type UserController struct {
	db            *gorm.DB
	SessionConfig *gorm.Session
}

func (uc *UserController) GetSession() *gorm.DB {
	return uc.db.Session(uc.SessionConfig)
}

// POST /token
func (uc *UserController) PostToken(c *gin.Context) {
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
				"message": "Invalid input",
			},
		)
		return
	}
	session := uc.GetSession()
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
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error":   "Unauthorized",
				"message": "Invalid credentials",
			},
		)
		return
	}

	credentialError := dbUser.CheckPassword(apiUser.Password)
	if credentialError != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error":   "Unauthorized",
				"message": "Invalid credentials",
			},
		)
		return
	}

	tokenString, err := GenerateJWT(dbUser.ID)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error":   "Unauthorized",
				"message": "Invalid credentials",
			},
		)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"id":         dbUser.ID,
		"token":      tokenString,
		"token_type": "Bearer",
		"expires_in": TokenExpiration.Seconds(),
	})
}

// POST /users
func (uc *UserController) PostUsers(c *gin.Context) {
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
	session := uc.GetSession()
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
func (uc *UserController) GetUsers(c *gin.Context) {
	var dbUsers []db.User

	session := uc.GetSession()
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
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	var dbUser = db.User{ID: uint(userId)}

	session := uc.GetSession()
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
// User can only delete itself
func (uc *UserController) DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	err := VerifyID(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error":   "Forbidden",
				"message": "You do not have access to this resource.",
			},
		)
		return
	}
	var dbUser = db.User{ID: uint(userId)}

	session := uc.GetSession()
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
func (uc *UserController) PostFavourites(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	err := VerifyID(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error":   "Forbidden",
				"message": "You do not have access to this resource.",
			},
		)
		return
	}
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
	session := uc.GetSession()

	err = session.Preload(clause.Associations).Model(&dbUser).Association("Favourites").Append(&dbAsset)
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
func (uc *UserController) GetFavourites(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	err := VerifyID(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error":   "Forbidden",
				"message": "You do not have access to this resource.",
			},
		)
		return
	}
	var dbUser = db.User{ID: uint(userId)}
	var dbAssets []*db.Asset

	session := uc.GetSession()
	// Prevent ErrRecordNotFound
	err = session.Preload(clause.Associations).Model(&dbUser).Association("Favourites").Find(&dbAssets)
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
func (uc *UserController) GetFavouriteByID(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	err := VerifyID(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error":   "Forbidden",
				"message": "You do not have access to this resource.",
			},
		)
		return
	}
	favId := c.Param("favId")
	assetId, _ := strconv.Atoi(favId)
	var dbUser = db.User{ID: uint(userId)}
	var dbAssets []*db.Asset

	session := uc.GetSession()
	// Prevent ErrRecordNotFound
	// Load only assets that belong to the user and have the correct id (which should be 0 or 1)
	result := session.Debug().Model(&db.Asset{ID: uint(assetId)}).Joins("INNER JOIN user_assets ua ON ua.asset_id = assets.id AND ua.user_id = ?", dbUser.ID).Preload(clause.Associations).Find(&dbAssets)
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
		c.PureJSON(http.StatusOK, dbAssets[0])
		return
	}
}

// DELETE /users/:id/favourites/:favId
func (uc *UserController) DeleteFavouriteByID(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)
	err := VerifyID(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error":   "Forbidden",
				"message": "You do not have access to this resource.",
			},
		)
		return
	}
	favId := c.Param("favId")
	assetId, _ := strconv.Atoi(favId)
	var dbUser = db.User{ID: uint(userId)}
	var dbAsset = db.Asset{ID: uint(assetId)}

	session := uc.GetSession()
	err = session.Model(&dbUser).Association("Favourites").Delete(&dbAsset)
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

type AssetController struct {
	db            *gorm.DB
	SessionConfig *gorm.Session
}

func (ac *AssetController) GetSession() *gorm.DB {
	return ac.db.Session(ac.SessionConfig)
}

// GET /assets
func (ac *AssetController) GetAssets(c *gin.Context) {
	var dbAssets []db.Asset

	session := ac.GetSession()
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
func (ac *AssetController) PostAssets(c *gin.Context) {
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
	session := ac.GetSession()
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
func (ac *AssetController) GetAssetByID(c *gin.Context) {
	id := c.Param("id")
	assetId, _ := strconv.Atoi(id)
	var dbAsset = db.Asset{ID: uint(assetId)}

	session := ac.GetSession()
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
func (ac *AssetController) PutAssetByID(c *gin.Context) {
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

	session := ac.GetSession()
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
func (ac *AssetController) PatchAssetByID(c *gin.Context) {
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

	session := ac.GetSession()
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
func (ac *AssetController) DeleteAssetByID(c *gin.Context) {
	id := c.Param("id")
	assetId, _ := strconv.Atoi(id)
	dbAsset := db.Asset{ID: uint(assetId)}

	session := ac.GetSession()
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
