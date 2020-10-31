package main

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	// .envの読み取り
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("not found .env file: %v\n", err)
	}

	apiKey := os.Getenv("APIKEY")
	fmt.Println(apiKey)
	authDomain := os.Getenv("AUTHDOMAIN")
	databaseUrl := os.Getenv("DATABASEURL")
	projectId := os.Getenv("PROJECTID")
	storageBucket := os.Getenv("STORAGEBUCKET")
	messagingSenderId := os.Getenv("MESSAGINGSENDERID")
	appId := os.Getenv("APPID")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{
			"apiKey":            apiKey,
			"authDomain":        authDomain,
			"databaseUrl":       databaseUrl,
			"projectId":         projectId,
			"storageBucket":     storageBucket,
			"messagingSenderId": messagingSenderId,
			"appId":             appId,
		})
	})

	router.POST("/auth", authByFirebase)

	router.Run()
}

func authByFirebase(ctx *gin.Context) {
	opt := option.WithCredentialsFile("secretkey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Panic(fmt.Errorf("error initializing app: %v", err))
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	// トークンを確認
	idToken := ctx.PostForm("id_token")
	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	email := token.Claims["email"]
	log.Printf("email: %v\n", email)

	ctx.HTML(200, "result.html", gin.H{
		"idToken": idToken,
		"email":   email,
	})
}
