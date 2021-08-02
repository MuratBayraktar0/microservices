package api

import (
	"fmt"

	"example.com/product/service"
	"github.com/gofiber/fiber"
)

type ProductDTO struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float32              `json:"price"`
	Images      []string             `json:"images"`
	Properties  []ProductPropertyDTO `json:"properties"`
}

type ProductPropertyDTO struct {
	PropertyName    string `json:"name"`
	PropertyContent string `json:"content"`
}

type Api struct {
	service *service.Service
}

func NewAPI(service *service.Service) *Api {
	return &Api{
		service: service,
	}
}

func (api *Api) ProductPostEndpoint(c *fiber.Ctx) {
	productDTO := ProductDTO{}
	err := c.BodyParser(&productDTO)
	if err != nil {
		fmt.Println("api.go => Error: ", err)
		c.Status(fiber.StatusBadRequest)
	}
	productModel := ConvertDTOToModel(productDTO)
	result, err := api.service.ProductPostService(productModel)
	if err != nil {
		fmt.Println("api.go => Error: ", err)
		c.Status(fiber.StatusBadRequest)
	}

	switch err {
	case nil:
		c.JSON(result)
		c.Status(fiber.StatusCreated)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
}

func (api *Api) ProductGetEndpoint(c *fiber.Ctx) {
	productId := c.Params("id")

	result, err := api.service.ProductGetService(productId)

	switch err {
	case nil:
		c.JSON(result)
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
}

func (api *Api) ProductsGetEndpoint(c *fiber.Ctx) {
	result, err := api.service.ProductsGetService()

	switch err {
	case nil:
		c.JSON(result)
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
}

func ConvertDTOToModel(productDTO ProductDTO) service.ProductModel {

	properties := []service.ProductPropertyModel{}
	for _, v := range productDTO.Properties {
		property := service.ProductPropertyModel{PropertyName: v.PropertyName, PropertyContent: v.PropertyContent}
		properties = append(properties, property)
	}

	productModel := service.ProductModel{
		ID:          productDTO.ID,
		Name:        productDTO.Name,
		Description: productDTO.Description,
		Price:       productDTO.Price,
		Images:      productDTO.Images,
		Properties:  properties,
	}

	return productModel
}
