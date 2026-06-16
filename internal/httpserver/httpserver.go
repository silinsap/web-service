package httpserver

import (
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"web-service/internal/memory"
	"web-service/internal/models"
	"web-service/internal/storage/sqlite"

	"github.com/gin-gonic/gin"
)

// GetLinks возвращает список всех ссылок с поддержкой пагинации через query параметры limit и offset.
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
		// Получаем ссылки из БД с учетом пагинации
		links, err := db.GetLinks(limit, offset)
		if err != nil {
			log.Error("Error fetching links", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error fetching links"})
			return
		}
		c.IndentedJSON(http.StatusOK, links)
	}
}

// PostLinks принимает JSON с полем original_url, генерирует уникальный short_code и сохраняет пару в БД.
func PostLinks(log *slog.Logger, db *sqlite.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		var newLink models.Link

		if err := c.BindJSON(&newLink); err != nil {
			return
		}

		if newLink.Original_url == "" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "EmptyUrl"})
			return
		}
		// Генерируем уникальный short_code и сохраняем пару в БД
		newLink.Short_code = NewRandomString(6)
		_, err := db.SaveLink(newLink.Short_code, newLink.Original_url)
		if err != nil {
			log.Error("Error saving link", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error saving link"})
			return
		}

		c.IndentedJSON(http.StatusCreated, models.Link{Short_code: newLink.Short_code})
	}
}

// GetLinkByShortCode принимает short_code в URL, ищет соответствующую original_url и возвращает ее. Если short_code не найден, возвращает 404.
func GetLinkByShortCode(log *slog.Logger, db *sqlite.Storage, memory *memory.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		short_code := c.Param("short_code")
		// Сначала пытаемся найти ссылку в in-memory кеше, если не находим, то обращаемся к БД
		link, exist := memory.GetByShortCode(short_code, true)

		if !exist {

			linkDB, err := db.GetLink(short_code, true)
			if err != nil {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Link not found"})
				return
			}
			// Сохраняем найденную ссылку в in-memory кеш для ускорения последующих запросов
			memory.Links[short_code] = &linkDB
			c.IndentedJSON(http.StatusOK, models.Link{Original_url: linkDB.Original_url, Visits: linkDB.Visits})
		} else {
			// Если ссылка найдена в кеше, то возвращаем ее и увеличиваем счетчик посещений
			db.AddVisit(short_code, *link)
			c.IndentedJSON(http.StatusOK, models.Link{Original_url: link.Original_url, Visits: link.Visits})
		}

	}
}

// GetStatsByShortCode принимает short_code в URL и возвращает статистику по ссылке (кол-во посещений, дату создания). Если short_code не найден, возвращает 404.
func GetStatsByShortCode(log *slog.Logger, db *sqlite.Storage, memory *memory.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		short_code := c.Param("short_code")
		// Сначала пытаемся найти ссылку в in-memory кеше, если не находим, то обращаемся к БД
		link, exist := memory.GetByShortCode(short_code, false)

		if !exist {
			linkDB, err := db.GetLink(short_code, false)
			if err != nil {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Link not found"})
				return
			}
			c.IndentedJSON(http.StatusOK, linkDB)
		} else {
			c.IndentedJSON(http.StatusOK, link)
		}
	}
}

// DeleteLinkByShortCode принимает short_code в URL, удаляет соответствующую пару из БД и возвращает статус 200. Если short_code не найден, возвращает 404.
func DeleteLinkByShortCode(log *slog.Logger, db *sqlite.Storage, memory *memory.MemoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		short_code := c.Param("short_code")
		// Сначала удаляем ссылку из in-memory кеша, затем из БД
		memory.DeleteByShortCode(short_code)
		err := db.DeleteLink(short_code)
		if err != nil {
			log.Error("Error deleting link", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error deleting link"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Link deleted"})
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
