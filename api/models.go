package api

import "github.com/GlobalWebIndex/platform2.0-go-challenge/db"

// This file has API models on which requests are validateds

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) getDBUser() (db.User, error) {
	user := db.User{
		Username: u.Username,
	}
	err := user.SetPassword(string(u.Password))
	return user, err
}

type Favourite struct {
	ID uint `json:"id"`
}

func (u *Favourite) getDBAsset() (db.Asset, error) {
	asset := db.Asset{ID: u.ID}
	return asset, nil
}

type Chart struct {
	Title  string `json:"title"`
	TitleX string `json:"title_x"`
	TitleY string `json:"title_y"`
	Data   string `binding:"base64" json:"data"` // base64
}

type Insight struct {
	Description string `json:"description"`
}

type Characteristic struct {
	Gender        string `json:"gender" binding:"exists,alphanum,min=1,max=1"`
	BirthCountry  string `json:"birth_country"`
	AgeGroupRange string `json:"age_group"`
	SocMediaHours string `json:"social_media_hours"`
}

type Audience struct {
	Characteristics []Characteristic `json:"characteristics"`
}

type Asset struct {
	Chart    *Chart    `json:"chart"`
	Insight  *Insight  `json:"insight"`
	Audience *Audience `json:"audience"`
}

func (a *Asset) getDBAsset() (db.Asset, error) {
	var chart *db.Chart
	if a.Chart != nil {
		chart = &db.Chart{
			Title:  a.Chart.Title,
			TitleX: a.Chart.TitleX,
			TitleY: a.Chart.TitleY,
			Data:   a.Chart.Data,
		}
	}
	var insight *db.Insight
	if a.Insight != nil {
		insight = &db.Insight{
			Description: a.Insight.Description,
		}
	}
	var audience *db.Audience
	if a.Audience != nil {
		var characteristics []*db.Characteristic
		if len(a.Audience.Characteristics) > 0 {
			for _, c := range a.Audience.Characteristics {
				characteristics = append(characteristics, &db.Characteristic{
					Gender:        c.Gender,
					BirthCountry:  c.BirthCountry,
					AgeGroupRange: c.AgeGroupRange,
					SocMediaHours: c.SocMediaHours,
				})
			}
		}
		audience = &db.Audience{
			Characteristics: characteristics,
		}
	}

	asset := db.Asset{
		Chart:    chart,
		Insight:  insight,
		Audience: audience,
	}
	return asset, nil
}
