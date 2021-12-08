package database

import (
	"fmt"
	"time"

	"github.com/daniilsolovey/graphql-service/internal/tools"
	"github.com/go-pg/pg"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

type Repository interface {
	Close() error
	GetAllProducts() ([]Product, error)
	InsertSMSCode(phoneNumber, code string, expirationTimer time.Duration) error
	RegisterNewUser(user User) error
	GetUserByPhoneNumber(string) (*User, error)
	GetSMSCode(phoneNumber string) (*SMSCode, error)
}

type User struct {
	ID    int
	Name  string
	Phone string
}

type SMSCode struct {
	ID        int
	Phone     string
	Code      string
	ExpiredAt time.Time
}

type Product struct {
	ID   int
	Name string
}

type Database struct {
	name     string
	host     string
	port     string
	user     string
	password string
	client   *pg.DB
}

func NewDatabase(
	name, host, port, user, password string,
) *Database {
	database := &Database{
		name:     name,
		host:     host,
		user:     user,
		password: password,
		port:     port,
	}

	connection, err := database.connect()
	if err != nil {
		log.Fatal(err)
	}

	database.client = connection

	return database
}

func (database *Database) connect() (*pg.DB, error) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s%s",
		database.user,
		database.password,
		database.host,
		database.port,
		database.name,
		"?sslmode=disable",
	)

	databaseOptions, err := pg.ParseURL(databaseURL)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to connect to the database: %s",
			database.name,
		)
	}

	connection := pg.Connect(databaseOptions)

	return connection, nil
}

func (database *Database) Close() error {
	database.client.Close()
	return nil
}

func (database *Database) GetAllProducts() ([]Product, error) {
	log.Info("receiving all products from database")
	var products []Product
	err := database.client.Model(&products).Select()
	if err != nil {
		return nil, karma.Format(
			err,
			"error during receiving all products",
		)
	}

	log.Info("products successfully received")

	return products, nil
}

func (database *Database) InsertSMSCode(
	phoneNumber string,
	code string,
	expirationTimer time.Duration,
) error {
	log.Infof(
		nil, "writting sms code to database, phone_number: %s, code: %s",
		phoneNumber, code,
	)
	timeNow, err := tools.GetCurrentMoscowTime()
	if err != nil {
		return karma.Format(
			err,
			"unable to get current moscow time",
		)
	}

	exipiredAt := timeNow.Add(expirationTimer * time.Minute)

	smsCode := &SMSCode{
		Phone:     phoneNumber,
		Code:      code,
		ExpiredAt: exipiredAt,
	}

	_, err = database.client.Model(smsCode).
		OnConflict("(phone) DO UPDATE").
		Set("code = ?", code).
		Set("expired_at = ?", exipiredAt).
		Insert()

	if err != nil {
		return karma.Format(
			err,
			"unable to insert sms code to database",
		)
	}

	log.Info("sms code successfully written to the database")
	return nil
}

func (database *Database) RegisterNewUser(user User) error {
	log.Info("register new user")
	_, err := database.client.Model(&user).Insert()
	if err != nil {
		return karma.Format(
			err,
			"unable to insert user to database",
		)
	}

	log.Info("user successfully inserted to database")
	return nil
}

func (database *Database) GetUserByPhoneNumber(phoneNumber string) (*User, error) {
	log.Info("receiving user by phone number from database")
	var user User
	err := database.client.Model(&user).Where("phone = ?", phoneNumber).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}

		return nil, karma.Format(
			err,
			"error during receiving user by phone number",
		)
	}

	log.Info("user successfully received from database")

	return &user, nil

}

func (database *Database) GetSMSCode(phoneNumber string) (*SMSCode, error) {
	log.Info("receiving expiration date of sms code from database")
	var smsCode []SMSCode
	err := database.client.Model(&smsCode).Where("phone = ?", phoneNumber).Select()
	if err != nil {
		return nil, karma.Format(
			err,
			"error during receiving expiration date of sms code",
		)
	}

	if len(smsCode) != 0 {
		log.Info("expiration date of sms code successfully received")
		return &smsCode[0], nil
	}

	log.Info("sms code not found in database")
	return nil, nil
}
