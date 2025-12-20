package main

import (
	"fmt"
	"log"

	"github.com/GAKiknadze/one_time_secret_service/internal/cypher"
	"github.com/GAKiknadze/one_time_secret_service/internal/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	storage, err := storage.NewStorageDatabase(db)
	if err != nil {
		log.Fatalf("Failed to initialize storage database: %v", err)
	}

	cypher := &cypher.CypherBase{}
	app := fiber.New()

	app.Use("/static", static.New("./static"))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})

	app.Get("/s/:id", func(c fiber.Ctx) error {
		return c.SendFile("./templates/get_secret.html")
	})

	app.Post("/api/create", func(c fiber.Ctx) error {
		cypherKey, err := cypher.GenerateKey(128)
		if err != nil {
			log.Printf("Failed to generate cypher key: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate encryption key",
			})
		}

		data := c.FormValue("secret")

		encryptedData := cypher.Encrypt(cypherKey, data)

		secretId := uuid.New().ID()

		storageErr := storage.Save(secretId, []byte(encryptedData))
		if storageErr != nil {
			log.Printf("Failed to save secret: %v", storageErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save secret",
			})
		}

		preparedId := fmt.Sprintf("%x-%s", secretId, cypherKey)

		return c.JSON(fiber.Map{
			"id": preparedId,
		})
	})

	app.Post("/api/get", func(c fiber.Ctx) error {
		preparedId := c.FormValue("id")

		var secretId uint32
		var cypherKey string
		_, scanErr := fmt.Sscanf(preparedId, "%x-%s", &secretId, &cypherKey)
		if scanErr != nil {
			log.Printf("Failed to parse prepared ID: %v", scanErr)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		encryptedData, getErr := storage.Get(secretId)
		if getErr != nil {
			log.Printf("Failed to get secret: %v", getErr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Secret not found or already retrieved",
			})
		}

		decryptedData, decryptErr := cypher.Decrypt(cypherKey, encryptedData)
		if decryptErr != nil {
			log.Printf("Failed to decrypt secret: %v", decryptErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decrypt secret",
			})
		}

		deleteErr := storage.DeleteById(secretId)
		if deleteErr != nil {
			log.Printf("Failed to delete secret after retrieval: %v", deleteErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete secret after retrieval",
			})
		}

		return c.JSON(fiber.Map{
			"secret": decryptedData,
		})
	})

	log.Fatal(app.Listen(":8000"))
}
