package service

import (
	"fmt"
	"time"

	"example.com/product/repository"
	"github.com/google/uuid"
)

type ProductPropertyModel struct {
	PropertyName    string `json:"name"`
	PropertyContent string `json:"content"`
}

type ProductModel struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       float32                `json:"price"`
	Images      []string               `json:"images"`
	Properties  []ProductPropertyModel `json:"properties"`
	CratedAt    time.Time              `json:"createdat"`
	UpdatedAt   time.Time              `json:"updatedat"`
}

type ProductsModel struct {
	Products []ProductModel `json:"products"`
}

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func (service *Service) ProductPostService(productModel ProductModel) (*ProductModel, error) {
	productModel.ID = GenerateUUID()
	productModel.UpdatedAt = time.Now().Round(time.Minute).UTC()
	productModel.CratedAt = time.Now().Round(time.Minute).UTC()
	productEntity := ConvertProductModelToEntity(productModel)

	productEntity, err := service.repository.ProductAdd(*productEntity)
	if err != nil {
		fmt.Println("service.go => Error: ", err)
		return nil, err
	}

	return ConvertProductEntityToModel(*productEntity), nil
}

func (service *Service) ProductGetService(productId string) (*ProductModel, error) {
	productEntity, err := service.repository.ProductGet(productId)
	if err != nil {
		fmt.Println("service.go => Error: ", err)
		return nil, err
	}

	return ConvertProductEntityToModel(*productEntity), nil
}

func (service *Service) ProductsGetService() (*ProductsModel, error) {
	productsEntity, err := service.repository.ProductsGet()
	if err != nil {
		fmt.Println("service.go => Error: ", err)
		return nil, err
	}

	return ConvertProductsEntityToModel(*productsEntity), nil
}

func ConvertProductModelToEntity(productModel ProductModel) *repository.ProductEntity {

	properties := []repository.ProductPropertyEntity{}
	for _, v := range productModel.Properties {
		property := repository.ProductPropertyEntity{PropertyName: v.PropertyName, PropertyContent: v.PropertyContent}
		properties = append(properties, property)
	}

	productEntity := repository.ProductEntity{
		ID:          productModel.ID,
		Name:        productModel.Name,
		Description: productModel.Description,
		Price:       productModel.Price,
		Images:      productModel.Images,
		Properties:  properties,
		CratedAt:    productModel.CratedAt,
		UpdatedAt:   productModel.UpdatedAt,
	}

	return &productEntity
}

func ConvertProductEntityToModel(productEntity repository.ProductEntity) *ProductModel {

	properties := []ProductPropertyModel{}
	for _, v := range productEntity.Properties {
		property := ProductPropertyModel{PropertyName: v.PropertyName, PropertyContent: v.PropertyContent}
		properties = append(properties, property)
	}

	productModel := ProductModel{
		ID:          productEntity.ID,
		Name:        productEntity.Name,
		Description: productEntity.Description,
		Price:       productEntity.Price,
		Images:      productEntity.Images,
		Properties:  properties,
		CratedAt:    productEntity.CratedAt,
		UpdatedAt:   productEntity.UpdatedAt,
	}

	return &productModel
}

func ConvertProductsEntityToModel(productsEntity repository.ProductsEntity) *ProductsModel {
	productsModel := ProductsModel{}

	for _, v := range productsEntity.Products {
		productsModel.Products = append(productsModel.Products, *ConvertProductEntityToModel(v))
	}

	return &productsModel
}
