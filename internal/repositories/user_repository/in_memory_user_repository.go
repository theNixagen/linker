package user_repository

import "context"

type InMemoryUserRepository struct {
	Users []User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		Users: []User{},
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user User) (int32, error) {
	for _, u := range r.Users {
		if u.Username == user.Username {
			return 0, ErrDuplicatedUsername
		}
		if u.Email == user.Email {
			return 0, ErrDuplicatedEmail
		}
	}

	r.Users = append(r.Users, user)
	return user.ID, nil
}

func (r *InMemoryUserRepository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	for _, u := range r.Users {
		if u.Username == username {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (r *InMemoryUserRepository) UpdateBio(ctx context.Context, username, bio string) error {
	for i, u := range r.Users {
		if u.Username == username {
			r.Users[i].Bio = bio
			return nil
		}
	}
	return ErrUserNotFound
}

func (r *InMemoryUserRepository) UpdateProfilePhoto(ctx context.Context, username, objectName string) error {
	for i, u := range r.Users {
		if u.Username == username {
			r.Users[i].ProfilePicture = objectName
			return nil
		}
	}
	return ErrUserNotFound
}

func (r *InMemoryUserRepository) UpdateBannerPhoto(ctx context.Context, username, objectName string) error {
	for i, u := range r.Users {
		if u.Username == username {
			r.Users[i].BannerPicture = objectName
			return nil
		}
	}
	return ErrUserNotFound
}
