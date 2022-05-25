package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mythosmystery/typenotes-go-graphql/auth"
	"github.com/mythosmystery/typenotes-go-graphql/graph/generated"
	"github.com/mythosmystery/typenotes-go-graphql/graph/model"
	"github.com/mythosmystery/typenotes-go-graphql/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateNote(ctx context.Context, input model.NewNote) (*model.Note, error) {
	user := middleware.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	result, err := r.DB.Note.InsertOne(ctx, model.Note{
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: time.Now().UnixMilli(),
		CreatedBy: user,
	})
	if err != nil {
		return nil, err
	}
	res := r.DB.Note.FindOne(ctx, bson.M{"_id": result.InsertedID.(primitive.ObjectID)})
	var note model.Note
	err = res.Decode(&note)
	if err != nil {
		return nil, err
	}
	id, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return nil, err
	}
	update, err := r.DB.User.UpdateByID(ctx, id, bson.M{"$push": bson.M{"notes": note}})
	if err != nil {
		return nil, err
	}
	if update.ModifiedCount == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &note, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	result, err := r.DB.User.InsertOne(ctx, model.User{
		Email:     input.Email,
		Password:  hash,
		Name:      input.Name,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, err
	}
	res := r.DB.User.FindOne(ctx, bson.M{"_id": result.InsertedID.(primitive.ObjectID)})
	var user model.User
	err = res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mutationResolver) UpdateNote(ctx context.Context, id string, input model.UpdateNote) (*model.Note, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res, err := r.DB.Note.UpdateByID(ctx, _id, bson.M{"$set": bson.M{
		"title":     input.Title,
		"content":   input.Content,
		"updatedAt": time.Now().UnixMilli(),
	}})
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("note not found")
	}
	result := r.DB.Note.FindOne(ctx, bson.M{"_id": _id})
	var note model.Note
	err = result.Decode(&note)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *mutationResolver) DeleteNote(ctx context.Context, id string) (*bool, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res, err := r.DB.Note.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return nil, err
	}
	if res.DeletedCount == 0 {
		return nil, fmt.Errorf("note not found")
	}
	success := true
	return &success, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*bool, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res, err := r.DB.User.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return nil, err
	}
	if res.DeletedCount == 0 {
		return nil, fmt.Errorf("user not found")
	}
	success := true
	return &success, nil
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.Auth, error) {
	res := r.DB.User.FindOne(ctx, bson.M{"email": email})
	var user model.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	if isValid := user.ValidatePassword(password); !isValid {
		return nil, fmt.Errorf("invalid password")
	}
	token, refreshToken, _, err := auth.CreateTokens(user.ID)
	if err != nil {
		return nil, err
	}
	return &model.Auth{
		Token:        token,
		User:         &user,
		RefreshToken: refreshToken,
	}, nil
}

func (r *mutationResolver) Register(ctx context.Context, input model.NewUser) (*model.Auth, error) {
	user, err := r.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}
	token, refreshToken, _, err := auth.CreateTokens(user.ID)
	if err != nil {
		return nil, err
	}
	return &model.Auth{
		Token:        token,
		User:         user,
		RefreshToken: refreshToken,
	}, nil
}

func (r *mutationResolver) RefreshTokens(ctx context.Context, token string) (*model.Auth, error) {
	userClaim, err := auth.ParseToken(token, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		return nil, err
	}
	user, err := model.GetUserById(userClaim.UserID, r.DB)
	if err != nil {
		return nil, err
	}
	token, refreshToken, _, err := auth.CreateTokens(user.ID)
	if err != nil {
		return nil, err
	}
	return &model.Auth{
		Token:        token,
		User:         user,
		RefreshToken: refreshToken,
	}, nil
}

func (r *queryResolver) Notes(ctx context.Context) ([]*model.Note, error) {
	cursor, err := r.DB.Note.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var notes []*model.Note
	cursor.All(ctx, &notes)
	return notes, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	cursor, err := r.DB.User.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var users []*model.User
	cursor.All(ctx, &users)
	return users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user, err := model.GetUserById(id, r.DB)
	return user, err
}

func (r *queryResolver) Note(ctx context.Context, id string) (*model.Note, error) {
	noteID, _ := primitive.ObjectIDFromHex(id)
	res := r.DB.Note.FindOne(ctx, bson.M{"_id": noteID})
	var note model.Note
	err := res.Decode(&note)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user := middleware.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
