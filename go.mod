// +heroku goVersion go1.15
// +heroku install ./cmd/heroku

module github.com/NicolasDutronc/shoppinglist-be

go 1.15

require (
	github.com/NicolasDutronc/autokey v0.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/gorilla/websocket v1.4.2
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	go.mongodb.org/mongo-driver v1.4.6
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
)
