package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryEntity struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Products  []string  `bson:"products"`
	UpdatedAt time.Time `bson:"updatedat"`
	CreatedAt time.Time `bson:"createdat"`
}

type CategoriesEntity struct {
	Categories []CategoryEntity `bson:"categories"`
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

func (repository *Repository) CategoryAdd(categoryEntity CategoryEntity) (*CategoryEntity, error) {
	collection := repository.client.Database("category").Collection("categories")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	_, err := collection.InsertOne(ctx, categoryEntity)

	if err != nil {
		return nil, err
	}

	return repository.CategoryGet(categoryEntity.ID)
}

func (repository *Repository) CategoryGet(categoryId string) (*CategoryEntity, error) {
	collection := repository.client.Database("category").Collection("categories")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	categoryEntity := CategoryEntity{}
	err := collection.FindOne(ctx, bson.M{"_id": categoryId}).Decode(&categoryEntity)

	if err != nil {
		fmt.Println("InsertONE Error:", err)
		return nil, err
	}
	return &categoryEntity, nil
}

func (repository *Repository) ProductCategoryGet(productId string) (*CategoriesEntity, error) {
	collection := repository.client.Database("category").Collection("categories")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"products": productId})
	if err != nil {
		fmt.Println("InsertONE Error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	categoriesEntity := CategoriesEntity{}
	for cursor.Next(ctx) {
		categoryEntity := CategoryEntity{}
		cursor.Decode(&categoryEntity)
		categoriesEntity.Categories = append(categoriesEntity.Categories, categoryEntity)
	}

	return &categoriesEntity, nil
}
