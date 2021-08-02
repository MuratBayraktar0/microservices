package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"example.com/category/repository"
	"github.com/google/uuid"
)

type CategoriesModel struct {
	Categories []CategoryModel `json:"categories"`
}

type CategoryModel struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Products  []ProductModel `json:"products"`
	UpdatedAt time.Time      `json:"updatedat"`
	CreatedAt time.Time      `json:"createdat"`
}

type ProductModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	repository *repository.Repository
	port       int
}

func NewService(repository *repository.Repository, port int) *Service {
	return &Service{
		repository: repository,
		port:       port,
	}
}

func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func (service *Service) CategoryPostService(categoryModel CategoryModel) (*CategoryModel, error) {
	categoryModel.ID = GenerateUUID()
	categoryModel.UpdatedAt = time.Now().Round(time.Minute).UTC()
	categoryModel.CreatedAt = time.Now().Round(time.Minute).UTC()
	categoryEntity := ConvertCategoryModelToEntity(categoryModel)

	categoryEntity, err := service.repository.CategoryAdd(*categoryEntity)
	if err != nil {
		fmt.Println("service.go => Error: ", err)
		return nil, err
	}

	return ConvertCategoryEntityToModel(*categoryEntity), nil
}

func (service *Service) CategoryGetService(categoryId string) (*CategoryModel, error) {
	categoryEntity, err := service.repository.CategoryGet(categoryId)
	if err != nil {
		fmt.Println("service.go => Error: ", err)
		return nil, err
	}

	categoryModel := ConvertCategoryEntityToModel(*categoryEntity)
	for i, v := range categoryModel.Products {
		resp, err := http.Get(fmt.Sprint("http://product_service:", service.port, "/product/", v.ID))
		if err != nil {
			fmt.Println("service.go => Error: ", err)
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		productModel := ProductModel{}
		err = json.Unmarshal(body, &productModel)
		categoryModel.Products[i].ID = productModel.ID
		categoryModel.Products[i].Name = productModel.Name
	}

	return categoryModel, nil
}

func (service *Service) ProductCategoryGetService(productId string) (*CategoriesModel, error) {
	categoriesEntity, err := service.repository.ProductCategoryGet(productId)
	if err != nil {
		fmt.Println("service.go => Error: ", err)
		return nil, err
	}

	return ConvertCategoriesEntityToModel(*categoriesEntity), nil
}

func ConvertCategoryModelToEntity(categoryModel CategoryModel) *repository.CategoryEntity {

	categoryEntity := repository.CategoryEntity{
		ID:        categoryModel.ID,
		Name:      categoryModel.Name,
		UpdatedAt: categoryModel.UpdatedAt,
		CreatedAt: categoryModel.CreatedAt,
	}
	for _, v := range categoryModel.Products {
		categoryEntity.Products = append(categoryEntity.Products, v.ID)
	}

	return &categoryEntity
}

func ConvertCategoryEntityToModel(categoryEntity repository.CategoryEntity) *CategoryModel {
	categoryModel := CategoryModel{
		ID:   categoryEntity.ID,
		Name: categoryEntity.Name,

		UpdatedAt: categoryEntity.UpdatedAt,
		CreatedAt: categoryEntity.CreatedAt,
	}
	for _, v := range categoryEntity.Products {
		categoryModel.Products = append(categoryModel.Products, ProductModel{ID: v})
	}

	return &categoryModel
}

func ConvertCategoriesEntityToModel(categoriesEntity repository.CategoriesEntity) *CategoriesModel {

	categoriesModel := CategoriesModel{}
	for _, v := range categoriesEntity.Categories {
		categoriesModel.Categories = append(categoriesModel.Categories, *ConvertCategoryEntityToModel(v))
	}

	return &categoriesModel
}
