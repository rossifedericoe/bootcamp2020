package internal

import (
	"airbnbmetrics/internal/domain/listing"
	listingRepo "airbnbmetrics/internal/repository/listing"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"context"

)

type App struct {
	listingService listing.ListingService
}

func NewApp() *App {
	listingRepository := listingRepo.NewListingRepository(newMongoConn())
	listingService := listing.NewListingService(listingRepository)
	app := App{}
	app.listingService = listingService
	return &app
}

func newMongoConn() *mongo.Client {
	stringConn :=  "mongodb+srv://[USER]:[PASS]@[DB_HOST]/sample_airbnb?retryWrites=true&w=majority"
	dbConfig := options.Client().ApplyURI(stringConn)
	client, err := mongo.Connect(context.TODO(), dbConfig)

	if err != nil {
		log.Fatal("Cannot connect with mongo")
		panic(err)
	}

	return client
}

func (app *App) Run() {
	id := "10006546"
	list, err := app.listingService.GetByID(id)
	if err != nil{
		fmt.Println("Falló en encontrar el listing " + id)
	} else {
		fmt.Println("El nombre del listing " + id + " es " + list.Name)
		fmt.Println("El precio del listing " + id + " es " + list.Price.String())
	}

	minPriceStr := "50000000"
	minPrice, parseDecimalErr := primitive.ParseDecimal128(minPriceStr)
	if parseDecimalErr != nil {
		fmt.Println("Falló al convertir el decimal")
	}

	filteredListings, filterListingsErr := app.listingService.GetByMinPrice(minPrice)
	if filterListingsErr != nil {
		fmt.Println("Falló al filtrar por precio")
	}
	fmt.Println("La cantidad de listings con precio mayor a " + minPriceStr + " es ", len(filteredListings))
}