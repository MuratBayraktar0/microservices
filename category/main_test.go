package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"example.com/category/api"
	"example.com/category/repository"
	"example.com/category/service"
	"github.com/gofiber/fiber"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_CategoryAdd(t *testing.T) {
	Convey("Given category data", t, func() {
		SetupAppConfig()

		testRepository, err := repository.NewRepository(appConfig.DBUrl)
		if err != nil {
			fmt.Println("Repository Error: ", err)
		}
		testService := service.NewService(testRepository, appConfig.Port)
		testApi := api.NewAPI(testService)
		app := SetupService(testApi)

		categoryDTO := api.CategoryDTO{
			Name:     "Samsung Serisi",
			Products: []api.ProductDTO{{ID: "a726a88c-4368-4bf8-85ed-75d70ef5b9f7"}, {ID: "0688f385-d486-490e-9fda-cb683eedc951"}},
		}

		Convey("When category post request", func() {
			reqBody, _ := json.Marshal(categoryDTO)

			request, _ := http.NewRequest(http.MethodPost, "/category", bytes.NewReader(reqBody))

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Content-Length", strconv.Itoa(len(reqBody)))
			response, err := app.Test(request, 30000)
			So(err, ShouldBeNil)

			Convey("Then Status Code Should be 201", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusCreated)
			})

			Convey("Then Category Should be returned", func() {
				responseBody, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)

				categoryDTOr := api.CategoryDTO{}

				err = json.Unmarshal(responseBody, &categoryDTOr)
				So(err, ShouldBeNil)

				resp, err := http.Get(fmt.Sprint("http://product_service:", appConfig.Port, "/product/a726a88c-4368-4bf8-85ed-75d70ef5b9f7"))
				So(err, ShouldBeNil)

				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				productDTO := api.ProductDTO{}
				err = json.Unmarshal(body, &productDTO)
				So(err, ShouldBeNil)

				So(categoryDTOr.ID, ShouldNotBeEmpty)
				So(categoryDTOr.Name, ShouldEqual, categoryDTO.Name)
				So(len(categoryDTOr.Products), ShouldEqual, len(categoryDTO.Products))
				So(categoryDTOr.Products[0].ID, ShouldEqual, categoryDTO.Products[0].ID)
				So(categoryDTOr.Products[1].ID, ShouldEqual, categoryDTO.Products[1].ID)
				So(categoryDTOr.Products[0].Name, ShouldEqual, productDTO.Name)
			})
		})
	})
}

// func Test_CategoryGet(t *testing.T) {
// 	Convey("Given category data in database", t, func() {
// 		categoryId := uuid.New().String()

// 		SetupAppConfig()

// 		testRepository, err := repository.NewRepository(appConfig.DBUrl)
// 		if err != nil {
// 			fmt.Println("Repository Error: ", err)
// 		}
// 		testService := service.NewService(testRepository)
// 		testApi := api.NewAPI(testService)
// 		app := SetupService(testApi)

// 		categoryEntity := repository.CategoryEntity{
// 			ID:       categoryId,
// 			Name:     "Samsung Serisi",
// 			Products: []string{"a726a88c-4368-4bf8-85ed-75d70ef5b9f7", "0688f385-d486-490e-9fda-cb683eedc951"},
// 		}

// 		testRepository.CategoryAdd(categoryEntity)

// 		Convey("When category get request", func() {

// 			request, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/category/", categoryId), nil)

// 			request.Header.Add("Content-Type", "application/json")
// 			response, err := app.Test(request, 30000)
// 			So(err, ShouldBeNil)

// 			Convey("Then Status Code Should be 200", func() {
// 				So(response.StatusCode, ShouldEqual, fiber.StatusOK)
// 			})

// 			Convey("Then Category Should be returned", func() {
// 				responseBody, err := ioutil.ReadAll(response.Body)
// 				So(err, ShouldBeNil)

// 				categoryDTO := api.CategoryDTO{}

// 				err = json.Unmarshal(responseBody, &categoryDTO)
// 				So(err, ShouldBeNil)

// 				So(categoryDTO.ID, ShouldEqual, categoryId)
// 				So(categoryDTO.Name, ShouldEqual, categoryEntity.Name)
// 				So(len(categoryDTO.Products), ShouldEqual, len(categoryEntity.Products))
// 				So(categoryDTO.Products[0].ID, ShouldEqual, categoryEntity.Products[0])
// 				So(categoryDTO.Products[1].ID, ShouldEqual, categoryEntity.Products[1])
// 			})
// 		})
// 	})
// }

// func Test_ProductCategoryGet(t *testing.T) {
// 	Convey("Given category data in database", t, func() {
// 		categoryId := uuid.New().String()
// 		productId := "a726a88c-4368-4bf8-85ed-75d70ef5b9f7"

// 		SetupAppConfig()

// 		testRepository, err := repository.NewRepository(appConfig.DBUrl)
// 		if err != nil {
// 			fmt.Println("Repository Error: ", err)
// 		}
// 		testService := service.NewService(testRepository)
// 		testApi := api.NewAPI(testService)
// 		app := SetupService(testApi)

// 		categoryEntity := repository.CategoryEntity{
// 			ID:       categoryId,
// 			Name:     "Samsung Serisi",
// 			Products: []string{productId, "0688f385-d486-490e-9fda-cb683eedc951"},
// 		}

// 		testRepository.CategoryAdd(categoryEntity)

// 		Convey("When category with productid get request", func() {

// 			request, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/category/product/", productId), nil)

// 			request.Header.Add("Content-Type", "application/json")
// 			response, err := app.Test(request, 30000)
// 			So(err, ShouldBeNil)

// 			Convey("Then Status Code Should be 200", func() {
// 				So(response.StatusCode, ShouldEqual, fiber.StatusOK)
// 			})

// 			Convey("Then Category Should be returned", func() {
// 				responseBody, err := ioutil.ReadAll(response.Body)
// 				So(err, ShouldBeNil)

// 				categoriesDTO := api.CategoriesDTO{}

// 				err = json.Unmarshal(responseBody, &categoriesDTO)
// 				So(err, ShouldBeNil)

// 				lastIndex := len(categoriesDTO.Categories) - 1
// 				So(categoriesDTO.Categories[lastIndex].ID, ShouldEqual, categoryId)
// 				So(categoriesDTO.Categories[lastIndex].Name, ShouldEqual, categoryEntity.Name)
// 				So(len(categoriesDTO.Categories[0].Products), ShouldEqual, len(categoryEntity.Products))
// 				So(categoriesDTO.Categories[lastIndex].Products[0], ShouldEqual, categoryEntity.Products[0])
// 				So(categoriesDTO.Categories[lastIndex].Products[1], ShouldEqual, categoryEntity.Products[1])
// 			})
// 		})
// 	})
// }
