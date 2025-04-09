package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sangketkit01/personal-block/internal/config"
	"github.com/sangketkit01/personal-block/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var Repo *DBRepo

type DBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewDBRepo(app *config.AppConfig, db *sql.DB) *DBRepo {
	return &DBRepo{
		App: app,
		DB:  db,
	}
}

func CreateRepo(repo *DBRepo){
	Repo = repo
}

func (repo *DBRepo) AllUsers() error {
	return nil
}

func (repo *DBRepo) InsertUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users ( username, email, phone, name, password, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7)
	`

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = repo.DB.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Phone,
		user.Name,
		hashPassword,
		time.Now(),
		time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (repo *DBRepo) LoginUser(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	query := `
		Select username, password From users Where username = $1
	`

	var hashedPassword string

	row := repo.DB.QueryRowContext(ctx,query,username)
	err := row.Scan(&username, &hashedPassword)
	if err != nil{
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("Password does not match")
	} else if err != nil{
		return err
	}
	

	return nil
}
