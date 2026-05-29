package httpserver

import (
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"web-service/internal/storage"
	"web-service/internal/storage/sqlite"

	"github.com/gin-gonic/gin"
)

func GetLinks(log *slog.Logger, db *sqlite.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 10
		offset := 0
		if offsetStr := c.Query("offset"); offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil {
				offset = o
			}
		}
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}

		links, err := db.GetLinks(limit, offset)
		if err != nil {
			log.Error("Error fetching links", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error fetching links"})
			return
		}
		c.IndentedJSON(http.StatusOK, links)
	}
}

func PostLinks(log *slog.Logger, db *sqlite.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		var newLink storage.Link

		if err := c.BindJSON(&newLink); err != nil {
			return
		}

		if newLink.Original_url == "" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "EmptyUrl"})
			return
		}

		newLink.Short_code = NewRandomString(6)
		_, err := db.SaveLink(newLink.Short_code, newLink.Original_url)
		if err != nil {
			log.Error("Error saving link", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error saving link"})
			return
		}

		c.IndentedJSON(http.StatusCreated, storage.Link{Short_code: newLink.Short_code})
	}
}

func GetLinkByShortCode(log *slog.Logger, db *sqlite.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		short_code := c.Param("short_code")

		link, err := db.GetLink(short_code, true)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Link not found"})
			return
		}

		c.IndentedJSON(http.StatusOK, storage.Link{Original_url: link.Original_url, Visits: link.Visits})
	}
}

func GetStatsByShortCode(log *slog.Logger, db *sqlite.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		short_code := c.Param("short_code")

		link, err := db.GetLink(short_code, false)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Link not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, link)
	}
}

func DeleteLinkByShortCode(log *slog.Logger, db *sqlite.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		short_code := c.Param("short_code")

		err := db.DeleteLink(short_code)
		if err != nil {
			log.Error("Error deleting link", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error deleting link"})
			return
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Link not found"})
	}
}

// NewRandomString generates random string with given size.
func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
