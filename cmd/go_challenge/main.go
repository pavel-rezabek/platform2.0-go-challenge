package main

import (
	"fmt"
	"os"

	"github.com/GlobalWebIndex/platform2.0-go-challenge/api"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getSqliteDB() *gorm.DB {
	database, err := gorm.Open(sqlite.Open("gorm.sqlite"), &gorm.Config{})
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		panic("Failed to connect to database")
	}
	db.Migrate(database)
	return database
}

func resolveAddress(host, port string) string {
	if port != "" {
		return host + ":" + port
	}
	return host + ":8080"
}

func main() {

	database := getSqliteDB()
 	// db.FillDB(database)
	engine := api.CreateEngine(database)

	addr := resolveAddress(os.Getenv("HOST"), os.Getenv("PORT"))
	engine.Run(addr)
}

/*

		>Auth:
		POST	/login
				authorization token
				refresh token
		POST	/refresh

	Security (user authorization so they cannot access other user's data)
	Tests
	db.Offset(offset_int).Limit(limit_int) for size and offset
*/
