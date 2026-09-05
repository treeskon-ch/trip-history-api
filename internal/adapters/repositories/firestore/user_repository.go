package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/iterator"
	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type userRepository struct {
	firestoreClient *firestore.Client
	authClient      *auth.Client
}

func NewUserRepository(firestoreClient *firestore.Client, authClient *auth.Client) ports.UserRepository {
	return &userRepository{
		firestoreClient: firestoreClient,
		authClient:      authClient,
	}
}

func (r *userRepository) CreateUserInAuth(ctx context.Context, email, password, name string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(name)

	u, err := r.authClient.CreateUser(ctx, params)
	if err != nil {
		return "", err
	}
	return u.UID, nil
}

func (r *userRepository) SaveUser(ctx context.Context, user domain.User) error {
	_, err := r.firestoreClient.Collection("users").Doc(user.UserID).Set(ctx, user)
	return err
}

func (r *userRepository) GetUser(ctx context.Context, uid string) (domain.User, error) {
	doc, err := r.firestoreClient.Collection("users").Doc(uid).Get(ctx)
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = doc.DataTo(&user)
	return user, err
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	iter := r.firestoreClient.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var u domain.User
		if err := doc.DataTo(&u); err == nil {
			users = append(users, u)
		}
	}
	return users, nil
}

func (r *userRepository) DeleteUserInAuth(ctx context.Context, uid string) error {
	return r.authClient.DeleteUser(ctx, uid)
}

func (r *userRepository) UpdateFCMToken(ctx context.Context, uid string, token string) error {
	_, err := r.firestoreClient.Collection("users").Doc(uid).Update(ctx, []firestore.Update{
		{Path: "fcmToken", Value: token},
	})
	return err
}

func (r *userRepository) DeleteUserInFirestore(ctx context.Context, uid string) error {
	_, err := r.firestoreClient.Collection("users").Doc(uid).Delete(ctx)
	return err
}

func (r *userRepository) ClearAllUsers(ctx context.Context) error {
	iter := r.firestoreClient.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		doc.Ref.Delete(ctx)
	}

	pager := r.authClient.Users(ctx, "")
	for {
		user, err := pager.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		r.authClient.DeleteUser(ctx, user.UID)
	}
	return nil
}

func (r *userRepository) GenerateEmailVerificationLink(ctx context.Context, email string) (string, error) {
	return r.authClient.EmailVerificationLink(ctx, email)
}
