package entity

import (
	"strings"
	"time"

	"github.com/arvinpaundra/cent/user/core/trait"
	"github.com/google/uuid"
	"github.com/sqids/sqids-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	trait.Updateable

	ID        int64      `json:"id" redis:"id"`
	Email     string     `json:"email" redis:"email"`
	Password  *string    `json:"password" redis:"-"`
	Fullname  string     `json:"fullname" redis:"-"`
	Key       string     `json:"key" redis:"-"`
	Slug      *string    `json:"slug" redis:"slug"`
	Image     *string    `json:"image" redis:"-"`
	DeletedAt *time.Time `json:"deleted_at" redis:"-"`
}

func (e *User) IsNew() bool {
	return e.ID <= 0
}

func (e *User) GeneratePassword(password string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashed := string(b)

	e.Password = &hashed

	return nil
}

func (e *User) ComparePassword(password string) bool {
	if e.Password != nil {
		err := bcrypt.CompareHashAndPassword([]byte(*e.Password), []byte(password))

		return err == nil
	}

	return false
}

func (e *User) GenerateSlug() error {
	s, err := sqids.New(sqids.Options{
		MinLength: 8,
	})

	if err != nil {
		return err
	}

	slug, err := s.Encode([]uint64{uint64(e.ID)})
	if err != nil {
		return err
	}

	e.Slug = &slug

	return nil
}

func (e *User) GenerateKey() error {
	key, err := uuid.NewV7()
	if err != nil {
		return err
	}

	sanitizedKey := strings.ReplaceAll(key.String(), "-", "")

	e.Key = sanitizedKey

	return nil
}

func (e *User) IsEmpty() bool {
	return e == nil
}
