package main

//服务端（使用 map）
import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func login(ctx echo.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	if username == "Lucas" && password == "123456" {
		token := jwt.New(jwt.SigningMethodHS256)

		//Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Lucas Snow"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		//Generate encode token and send it as reponse
		//生成密码令牌and将它作为响应发送
		str, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]string{"token": str})
	}
	return echo.ErrUnauthorized
}

func accessiable(ctx echo.Context) error {
	return ctx.String(200, "Accessaiable")
}

func restricted(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return ctx.String(200, "欢迎："+name)
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

	e.Logger.Fatal(e.Start(":1323"))
}
