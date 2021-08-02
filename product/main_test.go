package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"example.com/product/api"
	"example.com/product/repository"
	"example.com/product/service"
	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_ProductAdd(t *testing.T) {
	Convey("Given product data", t, func() {
		SetupAppConfig()

		testRepository, err := repository.NewRepository(appConfig.DBUrl)
		if err != nil {
			fmt.Println("Repository Error: ", err)
		}
		testService := service.NewService(testRepository)
		testApi := api.NewAPI(testService)
		app := SetupService(testApi)

		images := []string{
			"https://st2.myideasoft.com/idea/fh/66/myassets/products/091/a16ig-ldyql-sy445.jpg?revision=1574530751",
			"https://yorumbudur.com/content/icerik/big/samsung-galaxy-note-4_yorumbudurcom.jpg",
		}

		properties := []api.ProductPropertyDTO{
			{PropertyName: "Pil Omru", PropertyContent: "2 Yil"},
			{PropertyName: "Cozunurluk", PropertyContent: "1080p full hd"},
		}

		productDTO := api.ProductDTO{
			Name:        "Cep telefonu",
			Description: "Samsung galaxy Note 4 Akilli Telefon",
			Price:       2999.99,
			Images:      images,
			Properties:  properties,
		}

		Convey("When product post request", func() {
			reqBody, _ := json.Marshal(productDTO)

			request, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewReader(reqBody))

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Content-Length", strconv.Itoa(len(reqBody)))
			response, err := app.Test(request, 30000)
			So(err, ShouldBeNil)

			Convey("Then Status Code Should be 201", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusCreated)
			})

			Convey("Then Product Should be returned", func() {
				responseBody, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)

				productDTOr := api.ProductDTO{}

				err = json.Unmarshal(responseBody, &productDTOr)
				So(err, ShouldBeNil)

				So(productDTOr.ID, ShouldNotBeEmpty)
				So(productDTOr.Name, ShouldEqual, productDTO.Name)
				So(productDTOr.Description, ShouldEqual, productDTO.Description)
				So(productDTOr.Price, ShouldEqual, productDTO.Price)
				So(len(productDTOr.Images), ShouldEqual, len(productDTO.Images))
				So(len(productDTOr.Properties), ShouldEqual, len(productDTO.Properties))
				So(productDTOr.Properties[0].PropertyName, ShouldEqual, productDTO.Properties[0].PropertyName)
				So(productDTOr.Properties[1].PropertyName, ShouldEqual, productDTO.Properties[1].PropertyName)
				So(productDTOr.Properties[0].PropertyContent, ShouldEqual, productDTO.Properties[0].PropertyContent)
				So(productDTOr.Properties[1].PropertyContent, ShouldEqual, productDTO.Properties[1].PropertyContent)
			})

		})
	})
}

func Test_ProductGet(t *testing.T) {
	Convey("Given product data in database", t, func() {
		productId := uuid.New().String()
		SetupAppConfig()

		testRepository, err := repository.NewRepository(appConfig.DBUrl)
		if err != nil {
			fmt.Println("Repository Error: ", err)
		}
		testService := service.NewService(testRepository)
		testApi := api.NewAPI(testService)
		app := SetupService(testApi)

		images := []string{
			"https://st2.myideasoft.com/idea/fh/66/myassets/products/091/a16ig-ldyql-sy445.jpg?revision=1574530751",
			"https://yorumbudur.com/content/icerik/big/samsung-galaxy-note-4_yorumbudurcom.jpg",
		}

		properties := []repository.ProductPropertyEntity{
			{PropertyName: "Pil Omru", PropertyContent: "2 Yil"},
			{PropertyName: "Cozunurluk", PropertyContent: "1080p full hd"},
		}

		productEntity := repository.ProductEntity{
			ID:          productId,
			Name:        "Cep telefonu",
			Description: "Samsung galaxy Note 4 Akilli Telefon",
			Price:       2999.99,
			Images:      images,
			Properties:  properties,
		}

		testRepository.ProductAdd(productEntity)

		Convey("When product get request", func() {

			request, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/product/", productId), nil)

			request.Header.Add("Content-Type", "application/json")
			response, err := app.Test(request, 30000)
			So(err, ShouldBeNil)

			Convey("Then Status Code Should be 200", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then Product Should be returned", func() {
				responseBody, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)

				productDTO := api.ProductDTO{}

				err = json.Unmarshal(responseBody, &productDTO)
				So(err, ShouldBeNil)

				So(productDTO.ID, ShouldEqual, productId)
				So(productDTO.Name, ShouldEqual, productEntity.Name)
				So(productDTO.Description, ShouldEqual, productEntity.Description)
				So(productDTO.Price, ShouldEqual, productEntity.Price)
				So(len(productDTO.Images), ShouldEqual, len(productEntity.Images))
				So(len(productDTO.Properties), ShouldEqual, len(productEntity.Properties))
				So(productDTO.Properties[0].PropertyName, ShouldEqual, productEntity.Properties[0].PropertyName)
				So(productDTO.Properties[1].PropertyName, ShouldEqual, productEntity.Properties[1].PropertyName)
				So(productDTO.Properties[0].PropertyContent, ShouldEqual, productEntity.Properties[0].PropertyContent)
				So(productDTO.Properties[1].PropertyContent, ShouldEqual, productEntity.Properties[1].PropertyContent)
			})
		})
	})
}
