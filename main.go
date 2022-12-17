package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Entry struct {
	Id      string `json:"id,omitempty" firestore:"id,omitempty`
	Author  string `json:"author" binding:"required" firestore:"author"`
	Title   string `json:"title" binding:"required" firestore:"title"`
	Content string `json:"content" binding:"required" firestore:"content"`
}

func main() {

	router := gin.Default()
	router.POST("/entry", NewEntry)
	router.GET("/entries", GetAllEntries)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // Default port if not specified
	}

	router.Run(":" + port)
}

func GetAllEntries(c *gin.Context) {

	ctx := context.Background()
	sa := option.WithCredentialsFile("ex2-app-firebase-credentials.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	iter := client.Collection("entries").Documents(ctx)
	entries := []Entry{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		var entry Entry
		err = doc.DataTo(&entry)
		entries = append(entries, entry)
	}
	c.JSON(http.StatusOK, entries)

	defer client.Close()
}

func NewEntry(c *gin.Context) {
	var json Entry
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	sa := option.WithCredentialsFile("ex2-app-firebase-credentials.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	docRef := client.Collection("entries").NewDoc()
	json.Id = docRef.ID
	writeRes, err := docRef.Set(ctx, json)

	fmt.Println(docRef)
	fmt.Println(writeRes)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	defer client.Close()

	c.JSON(http.StatusOK, json)
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to Venetasa!")
}
