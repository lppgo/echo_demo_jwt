package main

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func login(ctx echo.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	if username == "Lucas" && password == "123456" {
		// 设置自定义的claims
		claims := &jwtCustomClaims{
			"Lucas Stu",
			true,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}
		// with claims 生成一个token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// 生成编码token
		str, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return ctx.JSON(200, echo.Map{"token": str})
	}
	return echo.ErrUnauthorized
}
func accessiable(ctx echo.Context) error {
	return ctx.String(200, "有权限的")
}
func restricted(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	fmt.Println("000000")
	fmt.Println(name)
	return ctx.String(200, "Hello:"+name)
}
func main() {
	e := echo.New()
	//Middle
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)
	//
	e.GET("/", accessiable)

	//
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":1322"))
}
