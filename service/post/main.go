package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID      string `json:"id" bson:"id"`
	Title   string `json:"title" bson:"title"`
	Author  string `json:"author" bson:"author"`
	Pages   string `json:"pages,omitempty" bson:"pages,omitempty"`
	Edition string `json:"edition,omitempty" bson:"edition,omitempty"`
	Year    string `json:"year,omitempty" bson:"year,omitempty"`
}

var collection *mongo.Collection

func seedData() {
	ctx := context.Background()
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count > 0 {
		log.Println("✅ Seed already exists")
		return
	}

	var books []interface{}
	for i := 1; i <= 30; i++ {
		book := Book{
			ID:      "book" + strconv.Itoa(i),
			Title:   "Seed Book " + strconv.Itoa(i),
			Author:  "Author " + strconv.Itoa(i),
			Pages:   "100",
			Edition: "1st Edition",
			Year:    "2020",
		}
		books = append(books, book)
	}
	_, err := collection.InsertMany(ctx, books)
	if err != nil {
		log.Println("❌ Failed to insert seed data:", err)
		return
	}
	log.Println("✅ Seed data inserted")
}

func createBook(c echo.Context) error {
	var book Book
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusOK, echo.Map{"error": "Invalid input"})
	}
	if book.ID == "" || book.Title == "" || book.Author == "" {
		return c.JSON(http.StatusOK, echo.Map{"error": "Missing required fields"})
	}
	ctx := context.Background()
	res, err := collection.InsertOne(ctx, book)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{"error": err.Error()})
	}
	log.Println("✅ Book inserted with _id:", res.InsertedID)
	return c.JSON(http.StatusCreated, book)
}

func main() {
	e := echo.New()
	e.POST("/api/books", createBook)

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://mongodb:testdatabase@mongo:27017"
		log.Println("⚠️  Using fallback Mongo URI")
	} else {
		log.Println("✅ Using provided Mongo URI:", uri)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}
	collection = client.Database("testdatabase").Collection("books")
	seedData()

	e.Logger.Fatal(e.Start(":3030"))
}
