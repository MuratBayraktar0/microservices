package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductPropertyEntity struct {
	PropertyName    string `bson:"name"`
	PropertyContent string `bson:"content"`
}

type ProductEntity struct {
	ID          string                  `bson:"_id"`
	Name        string                  `bson:"name"`
	Description string                  `bson:"description"`
	Price       float32                 `bson:"price"`
	Images      []string                `bson:"images"`
	Properties  []ProductPropertyEntity `bson:"properties"`
	CratedAt    time.Time               `bson:"createdat"`
	UpdatedAt   time.Time               `bson:"updatedat"`
}

type ProductsEntity struct {
	Products []ProductEntity `bson:"products"`
}

type Repository struct {
	client *mongo.Client
}

func NewRepository(dbUrl string) (*Repository, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &Repository{client}, nil
}

func (repository *Repository) ProductAdd(productEntity ProductEntity) (*ProductEntity, error) {
	collection := repository.client.Database("product").Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	_, err := collection.InsertOne(ctx, productEntity)

	if err != nil {
		return nil, err
	}

	return repository.ProductGet(productEntity.ID)
}

func (repository *Repository) ProductGet(id string) (*ProductEntity, error) {
	collection := repository.client.Database("product").Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	productEntity := ProductEntity{}
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&productEntity)

	if err != nil {
		fmt.Println("InsertONE Error:", err)
		return nil, err
	}
	return &productEntity, nil
}

func (repository *Repository) ProductsGet() (*ProductsEntity, error) {
	collection := repository.client.Database("product").Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		fmt.Println("InsertONE Error:", err)
		return nil, err
	}

	defer cursor.Close(ctx)
	productsEntity := ProductsEntity{}
	for cursor.Next(ctx) {
		productEntity := ProductEntity{}
		cursor.Decode(&productEntity)
		productsEntity.Products = append(productsEntity.Products, productEntity)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &productsEntity, nil
}
