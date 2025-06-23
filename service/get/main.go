package main

import (
    "context"
    "net/http"
    "os"
    "time"

    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    e := echo.New()

    // 環境變數取得 URI
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        panic("MONGODB_URI is not set")
    }

    // 連線到 MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }
    defer client.Disconnect(ctx)

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "MongoDB connected ✅")
    })

    e.Logger.Fatal(e.Start(":3030"))
}
