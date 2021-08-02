package api

import (
	"fmt"

	"example.com/category/service"
	"github.com/gofiber/fiber"
)

type CategoriesDTO struct {
	Categories []CategoryDTO `json:"categories"`
}

type CategoryDTO struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Products []ProductDTO `json:"products"`
}

type ProductDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Api struct {
	service *service.Service
}

func NewAPI(service *service.Service) *Api {
	return &Api{
		service: service,
	}
}

func (api *Api) CategoryPostEndpoint(c *fiber.Ctx) {
	categoryDTO := CategoryDTO{}
	err := c.BodyParser(&categoryDTO)
	if err != nil {
		fmt.Println("api.go => Error: ", err)
		c.Status(fiber.StatusBadRequest)
	}
	categoryModel := ConvertCategoryDTOToModel(categoryDTO)
	result, err := api.service.CategoryPostService(*categoryModel)
	if err != nil {
		fmt.Println("api.go => Error: ", err)
		c.Status(fiber.StatusBadRequest)
	}

	switch err {
	case nil:
		c.JSON(ConvertCategoryModelToDTO(*result))
		c.Status(fiber.StatusCreated)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
}

func (api *Api) CategoryGetEndpoint(c *fiber.Ctx) {
	categoryId := c.Params("categoryId")

	result, err := api.service.CategoryGetService(categoryId)

	switch err {
	case nil:
		c.JSON(ConvertCategoryModelToDTO(*result))
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
}

func (api *Api) ProductCategoryGetEndpoint(c *fiber.Ctx) {
	productId := c.Params("productId")

	result, err := api.service.ProductCategoryGetService(productId)

	switch err {
	case nil:
		c.JSON(ConvertCategoriesDTOToModel(*result))
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
}

func ConvertCategoryDTOToModel(categoryDTO CategoryDTO) *service.CategoryModel {

	categoryModel := service.CategoryModel{
		ID:   categoryDTO.ID,
		Name: categoryDTO.Name,
	}

	for _, v := range categoryDTO.Products {
		productModel := service.ProductModel{ID: v.ID}
		categoryModel.Products = append(categoryModel.Products, productModel)
	}

	return &categoryModel
}

func ConvertCategoryModelToDTO(categoryModel service.CategoryModel) *CategoryDTO {

	categoryDTO := CategoryDTO{
		ID:   categoryModel.ID,
		Name: categoryModel.Name,
	}

	for _, v := range categoryModel.Products {
		categoryDTO.Products = append(categoryDTO.Products, ProductDTO{ID: v.ID, Name: v.Name})
	}
	fmt.Println(categoryDTO.Products)

	return &categoryDTO
}

func ConvertCategoriesDTOToModel(categoriesModel service.CategoriesModel) *CategoriesDTO {

	categoriesDTO := CategoriesDTO{}
	for _, v := range categoriesModel.Categories {
		categoriesDTO.Categories = append(categoriesDTO.Categories, *ConvertCategoryModelToDTO(v))
	}

	return &categoriesDTO
}
