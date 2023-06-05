package main

// Methods in this project are created as per Fiber Framework
// Fiber is a layer over the HTTP package

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// All these packages are available in github
	"PostgresWithGo/Postgres-with-Go/models"
	"PostgresWithGo/Postgres-with-Go/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

/*
* Book is custom datatype to define book
* Golang doesn't understand JSON directly
* So, we've to specify the json format
* which tells the Golang how the json will look like
 */
type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

// Custom datatype refering DB
type Repository struct {
	DB *gorm.DB
}

/*
* CreateBook function will be called when create_book api is hit
* context is a abstract layer over API in fiber and
* BodyParser method will help converting JSON in API into book format of GO
 */
func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"},
		)
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book deleted successfully",
	})
	return nil
}

/*
* GetBooks called when /books hits
* models.Books{} is in models/ folder
 */
func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

// Creating a repository function SetupRoutes
// Basically it is creating an endpoint for an API
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// storage is a package created in storage folder
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	/*
	* This code will create a connection to the database
	 */
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}

	// If db doesn't exist in Postgres, MigrateBook create the db
	err = models.MigrateBook(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	// Creating a repository with database
	r := Repository{
		DB: db,
	}
	/*
	* Fiber: Fiber provides a minimalistic and flexible API for creating web servers and handling HTTP requests and responses.
	* It focuses on high performance and aims to be one of the fastest web frameworks available in Go.
	* fiber.New: create a new instance of the Fiber web framework.
	* app.listen: used to start the Fiber application and make it listen for incoming HTTP requests on a specific address.
	 */

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
