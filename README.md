# GlobalWebIndex Engineering Challenge

This project is an implementation of an api that manages users, assets and relationships between them.


# Installation

## Local

Make sure you have `git` installed or download this repository from github manually.

```sh
git clone https://github.com/pavel-rezabek/platform2.0-go-challenge.git
```

This project requires golang version >=1.18 to run. Install it [here](https://go.dev/dl/)

Open terminal in the repository root and run:
```sh
go mod download
```

## Docker

Make sure docker is installed and your user can execute the `docker` command, [info](https://docs.docker.com/engine/install/linux-postinstall/).

```sh
docker build -f Dockerfile -t go_challenge:latest .
docker run --rm -p8080:8080 go_challenge:latest
```

Feel free to replace the first number in the `-p<host_port>:<container_port>` argument, which represents the port at which the app will be accessible on your host.

In this case that would be `localhost:8080`


# Usage

```sh
go run cmd/go_challenge/main.go
```

Host and port of the server can be changed by setting the `HOST` and `PORT` environment variables respectively.


##  Customising the database

This api is designed to run on any database supported by [gorm](https://gorm.io/docs/write_driver.html).

The example runs on sqlite, but same wrapper can be written for any database.

```golang
package main

import (
	"github.com/GlobalWebIndex/platform2.0-go-challenge/api"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/db"
	"gorm.io/gorm"
)

func main(){
    database, _ := gorm.Open(<put your dialector here>, &gorm.Config{})
    db.Migrate(database)
    // Optionally add test data
    db.FillDB(database)
    engine := api.CreateEngine(database)
    engine.Run(":8080")
}

```

The `engine` returned can be built-upon to extend this api.

## Example requests

All endpoint paths are defined in [api/engine.go](api/engine.go). 

### User creation

```sh
curl -X POST localhost:8080/api/v1/users -d '{"username": "test", "password": "testpass"}'
```

### Authentication + Authorization

Assuming you have `jq` command-line tool installed
```sh
AUTH_TOKEN=$(curl -X POST localhost:8080/api/v1/token -d '{"username":"test","password":"testpass"}' | jq -r '.token')
```
the `-r` parameter results in output without string quotes

Otherwise copy the `token` field from response manually
```sh
curl -X POST localhost:8080/api/v1/token -d '{"username":"test","password":"testpass"}'
# {"expires_in":3600,"id":1,"token":"eyJhbGciOiJIUz<redacted>Grdm8eOQ","token_type":"Bearer"}
AUTH_TOKEN="eyJhbGciOiJIUz<redacted>Grdm8eOQ"
```

Then use your token in the `Authorization` header
```sh
curl -X GET localhost:8080/api/v1/users -H "Authorization: Bearer ${AUTH_TOKEN}"
# [{"id":1,"username":"test"}]
curl -X GET localhost:8080/api/v1/users/1 -H "Authorization: Bearer ${AUTH_TOKEN}"
# {"id":1,"username":"test"}
curl -X GET localhost:8080/api/v1/users/2/favourites -H "Authorization: Bearer ${AUTH_TOKEN}"
# {"error":"Forbidden","message":"You do not have access to this resource."}
```

### Assets and favourites

Assets have a complex structure that is defined [here](api/models.go#L27). 
One asset can have one or more of the subassets (Chart, Insight, Audience). The assets can then be favourited by the user based on the asset id:

```sh
curl -X POST localhost:8080/api/v1/assets -H "Authorization: Bearer ${AUTH_TOKEN}" -d '{"insight": {"description": "A very great description"}}'
# {"id":1,"insight": {"description": "A very great description"}}
curl -X POST localhost:8080/api/v1/users/1/favourites -H "Authorization: Bearer ${AUTH_TOKEN}" -d '{"id": 1}'
# {"id":1}
```

The path to the favourite can be found in the `Location` header or by adding the returned id to the path:

```sh
curl -X GET localhost:8080/api/v1/users/1/favourites/1 -H "Authorization: Bearer ${AUTH_TOKEN}"
# {"id":1,"insight":{"description":"A very great description"}}
# Or among all of the favourites
curl -X GET localhost:8080/api/v1/users/1/favourites -H "Authorization: Bearer ${AUTH_TOKEN}"
# [{"id":1,"insight":{"description":"A very great description"}}]
```


## Further ideas

- Swagger ui for more user-friendly api documentation and invocation
- Different levels of access for users
    - Currently everyone has full control over assets
    - Only maintainers/owners should have the power to manage
- Use a certificate file or more complex secret for JWT signing
- Token refresh endpoint and refresh_token
    - User would not have to re-login, only refresh current token
- Deduplicate characteristics
    - Existing characteristics are ignored for newly created assets
- Remove orphaned characteristics from db on asset deletion
