package main

import (
	"log"
	"net/http"

	"example/web-service-gin/database"

	"github.com/gin-gonic/gin"
)

func main() {
    // Initialize the database
    database.InitDB()
    defer database.CloseDB()

    router := gin.Default()
    router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
    router.POST("/albums", postAlbums)

    log.Println("Starting server on localhost:8181")
    router.Run("localhost:8181")
}

func getAlbums(c *gin.Context) {
    albums, err := database.GetAllAlbums()
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
    var newAlbum database.Album

    if err := c.BindJSON(&newAlbum); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := database.CreateAlbum(newAlbum); err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
    id := c.Param("id")

    album, err := database.GetAlbumByID(id)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if album == nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
        return
    }

    c.IndentedJSON(http.StatusOK, album)
}
