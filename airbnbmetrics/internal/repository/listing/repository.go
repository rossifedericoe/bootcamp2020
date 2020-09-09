package listing

import (
	"airbnbmetrics/internal/domain/listing"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ListingRepositoryImpl struct {
	client *mongo.Client
}

func NewListingRepository(client *mongo.Client) *ListingRepositoryImpl {
	return &ListingRepositoryImpl{client:client}
}

func (repository *ListingRepositoryImpl) GetByID(id string) (*listing.Listing, error) {
	var listingFound listing.Listing

	filter := bson.M{"_id": id}

	repository.client.Database("sample_airbnb").
		Collection("listingsAndReviews").
		FindOne(context.TODO(), filter).
		Decode(&listingFound)
	return &listingFound, nil
}

func (repository *ListingRepositoryImpl) GetByMinPrice(minPrice primitive.Decimal128) ([]listing.Listing, error) {
	filter := bson.M{"price": bson.M{"$gt": minPrice}}

	cursor, err := repository.client.Database("sample_airbnb").
		Collection("listingsAndReviews").
		Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var listings []listing.Listing
	for cursor.Next(context.TODO()) {
		var listingFound listing.Listing
		cursor.Decode(&listingFound)
		listings = append(listings, listingFound)
	}
	return listings, nil
}

func (repository *ListingRepositoryImpl) GetAll() ([]listing.Listing, error){
	filter := bson.M{}

	cursor, err := repository.client.Database("sample_airbnb").
		Collection("listingsAndReviews").
		Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var listings []listing.Listing
	for cursor.Next(context.TODO()) {
		var listingFound listing.Listing
		cursor.Decode(&listingFound)
		listings = append(listings, listingFound)
	}
	return listings, nil
}
