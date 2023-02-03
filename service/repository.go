package service

import (
	validator2 "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ivandi1980/my-gofiber/models"
	"gorm.io/gorm"
	"net/http"
)

type Book struct {
	Author    string `json:"author" validate:"required"`
	Title     string `json:"title" validate:"required"`
	Publisher string `json:"publisher" validate:"required"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("book/create", r.CreateBook)
	api.Get("books", r.GetBooks)
	api.Get("book/:id", r.GetBook)
	api.Put("book/:id", r.UpdateBook)
	api.Delete("book/:id", r.DeleteBook)
}

// GetBooks GetBook Get All Books
func (r *Repository) GetBooks(context *fiber.Ctx) error {

	bookModels := &[]models.Book{}

	err := r.DB.Find(&bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books retrieved successfully", "data": bookModels})
	return nil
}

// GetBook Get Book by Id
func (r *Repository) GetBook(context *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Book{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Id is required"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book retrieved successfully", "data": bookModel})
	return nil
}

// CreateBook Creating Book
func (r *Repository) CreateBook(context *fiber.Ctx) error {

	book := Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Invalid Request"})
		return err
	}

	validator := validator2.New()
	err = validator.Struct(Book{})
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book created successfully"})
	return nil
}

// UpdateBook Updating Book
func (r *Repository) UpdateBook(context *fiber.Ctx) error {

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Invalid Request"})
		return nil
	}

	bookModel := &models.Book{}
	book := Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Invalid Request"})
		return err
	}

	err = r.DB.Model(bookModel).Where("id = ?", id).Updates(book).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Could not update book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book updated successfully"})
	return nil
}

// DeleteBook Deleting Book
func (r *Repository) DeleteBook(context *fiber.Ctx) error {

	bookModel := &models.Book{}

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Id Cannot be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Could not delete book"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book deleted successfully"})
	return nil
}
