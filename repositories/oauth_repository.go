package repositories

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"diandi-backend/domains"
	"diandi-backend/services"
)

const (
	oauthTokensCollection   = "oauth_tokens"
	oauthProfilesCollection = "oauth_profiles"
)

type mongoOAuthRepository struct {
	db *mongo.Database
}

// NewMongoOAuthRepository creates a new MongoDB repository for OAuth data
func NewMongoOAuthRepository(db *mongo.Database) services.OAuthRepository {
	return &mongoOAuthRepository{
		db: db,
	}
}

func (r *mongoOAuthRepository) SaveToken(ctx context.Context, token *domains.OAuthToken) error {
	collection := r.db.Collection(oauthTokensCollection)

	filter := bson.M{
		"userId":   token.UserID,
		"provider": token.Provider,
	}

	update := bson.M{
		"$set": token,
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

func (r *mongoOAuthRepository) GetToken(ctx context.Context, userID string, provider domains.OAuthProvider) (*domains.OAuthToken, error) {
	collection := r.db.Collection(oauthTokensCollection)

	filter := bson.M{
		"userId":   userID,
		"provider": provider,
	}

	var token domains.OAuthToken
	err := collection.FindOne(ctx, filter).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return &token, nil
}

func (r *mongoOAuthRepository) DeleteToken(ctx context.Context, userID string, provider domains.OAuthProvider) error {
	collection := r.db.Collection(oauthTokensCollection)

	filter := bson.M{
		"userId":   userID,
		"provider": provider,
	}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

func (r *mongoOAuthRepository) SaveProfile(ctx context.Context, profile *domains.OAuthProfile) error {
	collection := r.db.Collection(oauthProfilesCollection)

	filter := bson.M{
		"providerId": profile.ProviderID,
		"provider":   profile.Provider,
	}

	update := bson.M{
		"$set": profile,
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	return nil
}

func (r *mongoOAuthRepository) GetProfile(ctx context.Context, providerID string, provider domains.OAuthProvider) (*domains.OAuthProfile, error) {
	collection := r.db.Collection(oauthProfilesCollection)

	filter := bson.M{
		"providerId": providerID,
		"provider":   provider,
	}

	var profile domains.OAuthProfile
	err := collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &profile, nil
}
