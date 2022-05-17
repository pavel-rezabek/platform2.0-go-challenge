// Package db defines database schema using gorm struct models.
package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// FillDB adds test data to the database
func FillDB(database *gorm.DB) {
	// TODO: try this in one go
	user := User{Username: "test"}
	user.SetPassword("testpass")
	database.Create(&user)

	// Create the assets
	var assets = []Asset{{ID: 1}, {ID: 2}, {ID: 3}}
	database.Create(&assets)
	// Fill the assets
	var chart = Chart{Title: "test"}
	chart.Asset = assets[0]
	var insight = Insight{Description: "Test insight"}
	insight.Asset = assets[1]
	var audience = Audience{
		Asset:           assets[2],
		Characteristics: []*Characteristic{{Gender: "M", BirthCountry: "Czech Republic"}},
	}
	database.Create(&chart)
	database.Create(&insight)
	database.Create(&audience)
}

// Migrate migrates changes in the models into databse `db`
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Asset{})
	db.AutoMigrate(&Chart{})
	db.AutoMigrate(&Insight{})
	db.AutoMigrate(&Characteristic{})
	db.AutoMigrate(&Audience{})
}

type User struct {
	ID           uint     `gorm:"primarykey;not null;autoIncrement:true" json:"id"`
	Username     string   `gorm:"unique;not null" json:"username"`
	PasswordHash string   `gorm:"column:password;not null" json:"-"`
	Favourites   []*Asset `gorm:"many2many:user_assets;" json:"-"`
}

// Note: When Asset gets deleted, all of the linked types get deleted
type Asset struct {
	ID       uint      `gorm:"primarykey;not null;autoIncrement:true" json:"id"`
	Users    []*User   `gorm:"many2many:user_assets;" json:"-"`
	Chart    *Chart    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"chart,omitempty"`
	Insight  *Insight  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"insight,omitempty"`
	Audience *Audience `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"audience,omitempty"`
}

type Chart struct {
	ID      uint        `gorm:"primaryKey;not null;autoIncrement:true" json:"-"`
	AssetID uint        `gorm:"unique;not null" json:"-"`
	Asset   Asset       `json:"-"`
	Title   string      `gorm:"size:256" json:"title"`
	TitleX  string      `gorm:"size:256" json:"title_x"`
	TitleY  string      `gorm:"size:256" json:"title_y"`
	Data    interface{} `gorm:"type:bytes;serializer:gob" json:"data"`
}

type Insight struct {
	ID          uint   `gorm:"primaryKey;not null" json:"-"`
	AssetID     uint   `gorm:"unique;not null" json:"-"`
	Asset       Asset  `json:"-"`
	Description string `gorm:"size:1024" json:"description"`
}

type Characteristic struct {
	ID            uint        `gorm:"primaryKey;not null" json:"-"`
	Gender        string      `gorm:"size:1;index:,unique,composite:characteristic" json:"gender"`
	BirthCountry  string      `gorm:"size:64;index:,unique,composite:characteristic" json:"birth_country"`
	AgeGroupRange string      `gorm:"size:256;index:,unique,composite:characteristic" json:"age_group"`
	SocMediaHours string      `gorm:"size:256;index:,unique,composite:characteristic" json:"social_media_hours"`
	Audiences     []*Audience `gorm:"many2many:audience_characteristics;" json:"-"`
}

type Audience struct {
	ID              uint              `gorm:"primaryKey;not null" json:"-"`
	AssetID         uint              `gorm:"unique;not null" json:"-"`
	Asset           Asset             `json:"-"`
	Characteristics []*Characteristic `gorm:"many2many:audience_characteristics;" json:"characteristics,omitempty"`
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
