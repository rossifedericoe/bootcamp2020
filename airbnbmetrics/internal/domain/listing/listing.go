package listing

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"sync"
	"time"
)

type Listing struct {
	ID string `bson:"_id"`
	Name string `bson:"name"`
	Price primitive.Decimal128 `bson:"price"`
	Beds int `bson:"beds"`
	Score float64
	Reviews []Review `bson:"reviews"`
}

type Review struct {
	ID string `bson:"_id"`
}

type ListingService interface {
	GetByID(id string) (*Listing, error)
	GetByMinPrice(minPrice primitive.Decimal128) ([]Listing, error)
	GetAllScored() ([]Listing, error)
}

type ListingRepository interface {
	GetByID(id string) (*Listing, error)
	GetByMinPrice(minPrice primitive.Decimal128) ([]Listing, error)
	GetAll() ([]Listing, error)
}

type ListingServiceImpl struct {
	repository ListingRepository
}

func NewListingService(repository ListingRepository) *ListingServiceImpl {
	return &ListingServiceImpl{repository:repository}
}

func (service *ListingServiceImpl) GetByID(id string) (*Listing, error) {
	return service.repository.GetByID(id)
}

func (service *ListingServiceImpl) GetByMinPrice(minPrice primitive.Decimal128) ([]Listing, error) {
	return service.repository.GetByMinPrice(minPrice)
}

func (service *ListingServiceImpl) GetAllScored() ([]Listing, error) {

	// Me traigo todos los listings
	allListings, allListingsErr := service.repository.GetAll()
	if allListingsErr != nil {
		return nil, allListingsErr
	}

	// Calculo concurrentemente los scores: por price/beds y con el servicio "externo"
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(allListings))

	for i, _ := range allListings {
		go service.calculateBothScores(&allListings[i], &waitGroup)
	}

	waitGroup.Wait()

	// Calcular la métrica de confiables y no confiables
	trustedChannel := make(chan Listing, len(allListings))
	untrustedChannel := make(chan Listing, len(allListings))

	for _, currentlisting := range allListings {
		go service.classifyTrusteds(currentlisting, trustedChannel, untrustedChannel)
	}

	trustedListings := []Listing{}
	untrustedListings := []Listing{}
	for i:= 0; i < len(allListings); i++ {
		select {
			case trustedListing := <- trustedChannel:
				trustedListings = append(trustedListings, trustedListing)
			case untrustedListing := <- untrustedChannel:
				untrustedListings = append(untrustedListings, untrustedListing)
		}
	}

	fmt.Println("La cantidad de listings confiables es :" + fmt.Sprint(len(trustedListings)))
	fmt.Println("La cantidad de listings NO confiables es :" + fmt.Sprint(len(untrustedListings)))

	return allListings, nil
}

func (service *ListingServiceImpl) applyScoreByTitle(listingToScore *Listing) {
	if len(listingToScore.Name) > 10 {
		listingToScore.Score = listingToScore.Score * 2
	}
	time.Sleep(500 * time.Millisecond)
}

func (service *ListingServiceImpl) calculateBothScores(listingItem *Listing, wg *sync.WaitGroup) {
	price, parseErr := strconv.ParseFloat(listingItem.Price.String(), 64)
	if parseErr != nil {
		fmt.Println(parseErr)
	}
	listingItem.Score = price / float64(listingItem.Beds)
	service.applyScoreByTitle(listingItem)

	wg.Done()
}

func (sevice *ListingServiceImpl) classifyTrusteds(listingToClassify Listing, trustedChannel chan Listing, untrustedChannel chan Listing) {
	if len(listingToClassify.Reviews) > 50 {
		trustedChannel <- listingToClassify
	} else {
		untrustedChannel <- listingToClassify
	}
	time.Sleep(500 * time.Millisecond)
}