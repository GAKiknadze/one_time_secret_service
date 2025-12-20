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
	// Initialize SQLite in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize storage layer
	storage, err := storage.NewStorageDatabase(db)
	if err != nil {
		log.Fatalf("Failed to initialize storage database: %v", err)
	}

	// Create cipher instance for encryption/decryption
	cypher := &cypher.CypherBase{}
	app := fiber.New()

	// Serve static files (CSS, JavaScript)
	app.Use("/static", static.New("./static"))

	// Route: Serve home page for creating secrets
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})

	// Route: Serve secret retrieval page
	app.Get("/s/:id", func(c fiber.Ctx) error {
		return c.SendFile("./templates/get_secret.html")
	})

	// API Endpoint: Create a new secret
	app.Post("/api/create", func(c fiber.Ctx) error {
		// Generate random encryption key
		cypherKey, err := cypher.GenerateKey(128)
		if err != nil {
			log.Printf("Failed to generate cypher key: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate encryption key",
			})
		}

		// Extract secret from form data
		data := c.FormValue("secret")

		// Encrypt the secret data
		encryptedData := cypher.Encrypt(cypherKey, data)

		// Generate unique ID for storing the secret
		secretId := uuid.New().ID()

		// Save encrypted secret to storage
		storageErr := storage.Save(secretId, []byte(encryptedData))
		if storageErr != nil {
			log.Printf("Failed to save secret: %v", storageErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save secret",
			})

		}

		// Combine ID and key: the key is embedded in the URL (never stored on server)
		preparedId := fmt.Sprintf("%x-%s", secretId, cypherKey)

		// Return the secret link identifier to the client
		return c.JSON(fiber.Map{
			"id": preparedId,
		})
	})

	// API Endpoint: Retrieve and delete a secret
	app.Post("/api/get", func(c fiber.Ctx) error {
		// Get the prepared ID from request (contains both ID and encryption key)
		preparedId := c.FormValue("id")

		// Parse the prepared ID to extract secret ID and encryption key
		var secretId uint32
		var cypherKey string
		_, scanErr := fmt.Sscanf(preparedId, "%x-%s", &secretId, &cypherKey)
		if scanErr != nil {
			log.Printf("Failed to parse prepared ID: %v", scanErr)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		// Retrieve encrypted data from storage
		encryptedData, getErr := storage.Get(secretId)
		if getErr != nil {
			log.Printf("Failed to get secret: %v", getErr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Secret not found or already retrieved",
			})
		}

		// Decrypt the secret using the key from the URL
		decryptedData, decryptErr := cypher.Decrypt(cypherKey, encryptedData)
		if decryptErr != nil {
			log.Printf("Failed to decrypt secret: %v", decryptErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decrypt secret",
			})
		}

		// Delete the secret from storage (single-use access guarantee)
		deleteErr := storage.DeleteById(secretId)
		if deleteErr != nil {
			log.Printf("Failed to delete secret after retrieval: %v", deleteErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete secret after retrieval",
			})
		}

		// Return the decrypted secret to the client
		return c.JSON(fiber.Map{
			"secret": decryptedData,
		})
	})

	// Start the server on port 8000
	log.Fatal(app.Listen(":8000"))
}
