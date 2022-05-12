package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/mythosmystery/typenotes-go-graphql/graph/generated"
	"github.com/mythosmystery/typenotes-go-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateNote(ctx context.Context, input model.NewNote) (*model.Note, error) {
	result, err := r.DB.Note.InsertOne(ctx, model.Note{
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: time.Now().UnixMilli(),
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
	return &note, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateNote(ctx context.Context, id string, input model.UpdateNote) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteNote(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.Auth, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Register(ctx context.Context, input model.NewUser) (*model.Auth, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*model.Auth, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Notes(ctx context.Context) ([]*model.Note, error) {
	cursor, err := r.DB.Note.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	var notes []*model.Note
	cursor.All(ctx, &notes)
	return notes, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	cursor, err := r.DB.User.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	var users []*model.User
	cursor.All(ctx, &users)
	return users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Note(ctx context.Context, id string) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
