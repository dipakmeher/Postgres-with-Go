package main

// Methods in this project are created as per Fiber Framework

import (
	"log"
	"net/http"

	// All these packages are available in github
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
	r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
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

	// Creating a repository with database
	r := Repository{
		DB: db,
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
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
