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

func (repo *DBRepo) LoginUser(username, password string) (models.User ,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	query := `
		Select id, username, password From users Where username = $1
	`

	var user models.User
	var id int
	var hashedPassword string

	row := repo.DB.QueryRowContext(ctx,query,username)
	err := row.Scan(&id,&username, &hashedPassword)
	if err != nil{
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return user ,errors.New("Password does not match")
	} else if err != nil{
		return user, err
	}

	user, err = repo.GetUserFromID(id)
	if err != nil {
		return user , err
	}
	

	return user, nil
}

func (repo *DBRepo) GetUserFromID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	var user models.User
	
	query := `
		SELECT id, username, email, phone, name, password, coalesce(profile_image,'nothing'), created_at, updated_at
		FROM users WHERE id = $1
	`

	row := repo.DB.QueryRowContext(ctx,query,id)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Name,
		&user.Password,
		&user.ProfileImage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user , err
	}

	return user , nil
}

func (repo *DBRepo) GetAllBlocks() ([]models.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	var blocks []models.Block

	query := `
		SELECT id, users_id, content, created_at, updated_at FROM blocks ORDER BY id DESC
	`

	rows, err :=  repo.DB.QueryContext(ctx,query)
	if err != nil {
		return nil, err
	}


	defer rows.Close()

	for rows.Next(){
		var block models.Block
		err = rows.Scan(
			&block.ID,
			&block.UserID,
			&block.Content,
			&block.CreatedAt,
			&block.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		user, err := repo.GetUserFromID(block.UserID)
		if err != nil{
			return nil, err
		}

		block.User = user

		blocks = append(blocks, block)
	}

	if err = rows.Err() ; err != nil{
		return nil, err
	}

	return blocks, nil
}

func (repo *DBRepo) InsertNewBlock(userID int, content string) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	query := `
		INSERT INTO blocks (users_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := repo.DB.ExecContext(ctx,query,userID,content,time.Now(),time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (repo *DBRepo) GetBlockByUserID(id int) ([]models.Block, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	var blocks []models.Block

	query := `
		SELECT id, users_id, content, created_at, updated_at FROM blocks
		WHERE users_id = $1
	`

	rows, err :=  repo.DB.QueryContext(ctx,query,id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		var block models.Block
		err = rows.Scan(
			&block.ID,
			&block.UserID,
			&block.Content,
			&block.CreatedAt,
			&block.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		user, err := repo.GetUserFromID(block.UserID)
		if err != nil{
			return nil, err
		}

		block.User = user

		blocks = append(blocks, block)
	}

	if err = rows.Err() ; err != nil{
		return nil, err
	}

	return blocks, nil
}

func (repo *DBRepo) UpdateProfile(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	query := `
		UPDATE users SET name = $1 , email = $2 , phone = $3 WHERE id = $4
	`

	_, err := repo.DB.ExecContext(ctx,query,u.Name,u.Email,u.Phone,u.ID)
	if err != nil {
		return err
	}

	return  nil
}

func (repo *DBRepo) UpdateUserPassword(id int , oldPassword, newPassword string) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	query := `
		SELECT password FROM users WHERE id = $1
	`

	var oldHashedPassword string

	row := repo.DB.QueryRowContext(ctx,query,id)
	err := row.Scan(&oldHashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(oldHashedPassword),[]byte(oldPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword{
		return bcrypt.ErrMismatchedHashAndPassword
	}else if err != nil{
		return err
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword),bcrypt.DefaultCost)
	if err != nil{
		return err
	}

	query = `UPDATE users SET password = $1 WHERE id = $2`
	_, err = repo.DB.ExecContext(ctx,query,newHashedPassword,id)
	if err != nil{
		return err
	}

	return nil
}
