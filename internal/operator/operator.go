package operator

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/daniilsolovey/graphql-service/graph/model"
	"github.com/daniilsolovey/graphql-service/internal/config"
	"github.com/daniilsolovey/graphql-service/internal/database"
	"github.com/daniilsolovey/graphql-service/internal/tools"
	"github.com/dgrijalva/jwt-go"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

const (
	ERR_SMS_CODE_EXPIRED = "sms code has expired"
	ERR_INVALID_SMS_CODE = "invalid sms code"
)

type Operator interface {
	GetAllProducts() ([]*model.Product, error)
	RequestSignInCode(phoneNumber string) error
	SignInByCode(phoneNumber, code string) (model.SignInOrErrorPayload, error)
	Viewer(token string) (*model.Viewer, error)
}

type OperatorData struct {
	config   *config.Config
	database *database.Database
}

type Claims struct {
	UserPhone string `json:"userPhone"`
	jwt.StandardClaims
}

func NewOperator(
	config *config.Config,
	database *database.Database,
) *OperatorData {
	return &OperatorData{
		config:   config,
		database: database,
	}
}

func (operator *OperatorData) RequestSignInCode(phoneNumber string) error {
	randomNumber := getRandomNumber()
	log.Infof(nil, "sms sent to phone number: %s | code: %d", phoneNumber, randomNumber)
	// if not error after sending real message, then write code to database
	err := operator.database.InsertSMSCode(
		phoneNumber,
		strconv.Itoa(randomNumber),
		time.Duration(operator.config.SMS.ExpirationTimer),
	)
	if err != nil {
		return karma.Format(
			err,
			"unable to insert sms code to database",
		)
	}

	return nil
}

func (operator *OperatorData) checkSMSCode(phoneNumber, code string) error {
	smsCode, err := operator.database.GetSMSCode(phoneNumber)
	if err != nil {
		return karma.Format(
			err,
			"unable to get experation time for sms code by phone: %s",
			phoneNumber,
		)
	}

	timeNow, err := tools.GetCurrentMoscowTime()
	if err != nil {
		return karma.Format(
			err,
			"unable to get current moscow time",
		)
	}

	if smsCode.Code != code {
		return errors.New(ERR_INVALID_SMS_CODE)

	}

	if timeNow.After(smsCode.ExpiredAt) {
		return errors.New(ERR_SMS_CODE_EXPIRED)
	}

	return nil
}

func (operator *OperatorData) SignInByCode(phoneNumber, code string) (model.SignInOrErrorPayload, error) {
	err := operator.checkSMSCode(phoneNumber, code)
	if err != nil {
		switch err.Error() {
		case ERR_SMS_CODE_EXPIRED:
			log.Error(err)
			return &model.ErrorPayload{Message: ERR_SMS_CODE_EXPIRED}, nil
		case ERR_INVALID_SMS_CODE:
			log.Error(err)
			return &model.ErrorPayload{Message: ERR_INVALID_SMS_CODE}, nil
		default:
			return nil, karma.Format(
				err,
				"unable to check sms code",
			)

		}
	}

	user, err := operator.database.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to get user by phone from database",
		)
	}

	if user == nil {
		log.Info("registering a user")
		var newUser database.User
		newUser.Phone = phoneNumber
		err := operator.registerUser(newUser, code)
		if err != nil {
			return nil, karma.Format(
				err,
				"unable to register new user",
			)
		}

		user = &newUser
		log.Info("user registered successfully")
	}

	token, err := operator.createToken(
		phoneNumber,
		operator.config.Token.SecretKey,
	)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to create token",
		)
	}

	var signInPayload model.SignInPayload
	var signInViewer model.Viewer
	modelUser := &model.User{Phone: user.Phone}

	signInPayload.Token = token
	signInViewer.User = modelUser
	signInPayload.Viewer = &signInViewer

	return signInPayload, nil
}

func (operator *OperatorData) Viewer(token string) (*model.Viewer, error) {
	tokenClaims, err := operator.checkToken(token)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to check token",
		)
	}

	modelUser := model.User{Phone: tokenClaims.UserPhone}
	result := model.Viewer{User: &modelUser}
	return &result, nil
}

func (operator *OperatorData) createToken(userPhone, secretKey string) (string, error) {
	log.Info("creating token")
	timeNow, err := tools.GetCurrentMoscowTime()
	if err != nil {
		return "", karma.Format(
			err,
			"unable to get current moscow time",
		)
	}

	expirationTime := timeNow.Add(
		time.Duration(operator.config.Token.ExpirationTimer) * time.Minute,
	)
	claims := &Claims{
		UserPhone: userPhone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Local().Unix(),
		},
	}

	jwtKey := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", karma.Format(
			err,
			"unable to create jwt token",
		)
	}

	log.Info("token successfully created")
	return tokenString, nil
}

func (operator *OperatorData) checkToken(token string) (*Claims, error) {
	log.Info("checking token")
	if token == "" {
		return nil, errors.New("token is empty")
	}

	jwtKey := []byte(operator.config.Token.SecretKey)
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, karma.Format(
			err,
			"error with parsing token",
		)
	}

	if !tkn.Valid {
		return nil, errors.New("token is not valid")
	}

	log.Info("token is valid")
	return claims, nil
}

func (operator *OperatorData) registerUser(user database.User, code string) error {
	err := operator.database.RegisterNewUser(
		user,
	)

	if err != nil {
		return karma.Format(
			err,
			"unable to write new user to database",
		)
	}

	return nil
}

func (operator *OperatorData) GetAllProducts() ([]*model.Product, error) {
	products, err := operator.database.GetAllProducts()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to get all products from database",
		)
	}

	result := handleProducts(products)
	return result, nil
}

func getRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	result := rand.Intn(max-min+1) + min
	return result
}

func handleProducts(products []database.Product) []*model.Product {
	var result []*model.Product
	for _, item := range products {
		var product model.Product
		product.ID = item.ID
		product.Name = item.Name
		result = append(result, &product)
	}

	return result
}
