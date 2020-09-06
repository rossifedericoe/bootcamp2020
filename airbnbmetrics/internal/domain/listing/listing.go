package listing

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Listing struct {
	ID string `bson:"_id"`
	Name string `bson:"name"`
	Price primitive.Decimal128 `bson:"price"`
	Beds int `bson:"beds"`
}

type ListingService interface {
	GetByID(id string) (*Listing, error)
	GetByMinPrice(minPrice primitive.Decimal128) ([]Listing, error)
}

type ListingRepository interface {
	GetByID(id string) (*Listing, error)
	GetByMinPrice(minPrice primitive.Decimal128) ([]Listing, error)
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