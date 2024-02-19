package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"strings"
	"time"
	"url_shortner/db"
	"url_shortner/models"
)

func shortenURL(c *fiber.Ctx) error {
	body := new(models.Request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	rdb := db.CreateClient()
	defer rdb.Close()

	sqlStmt := `select url, shortened, expires_at from urls where url = ?;`
	row := rdb.QueryRow(sqlStmt, body.URL)

	var url string
	var shortened string
	var expiresAt time.Time

	err := row.Scan(&url, &shortened, &expiresAt)

	if url != "" {
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"original":   body.URL,
			"shortened":  shortened,
			"expires_at": expiresAt,
		})
	}

	host := c.Hostname()

	if strings.Contains(host, "localhost") {
		host = "http://" + host + "/"
	} else {
		host = "https://" + host + "/"
	}

	uid := uuid.New().String()[:8]
	shortenedUrl := host + uid
	expiresAt = time.Now().Add(time.Hour * 24 * 3)

	sqlStmt = `insert into urls (url, shortened, expires_at) values (?, ?, ?);`
	_, err = rdb.Exec(sqlStmt, body.URL, uid, expiresAt)

	if err != nil {
		log.Fatalf("cannot insert into database: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot insert into database",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"original":   body.URL,
		"shortened":  shortenedUrl,
		"expires_at": expiresAt,
	})
}

func resolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	rdb := db.CreateClient()
	defer rdb.Close()

	sqlStmt := `select url from urls where shortened = ? and expires_at > ?;`
	row := rdb.QueryRow(sqlStmt, url, time.Now())

	var original string
	err := row.Scan(&original)

	if err != nil {
		return err
	}

	return c.Redirect(original, fiber.StatusMovedPermanently)
}
