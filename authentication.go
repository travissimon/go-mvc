package mvc

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"fmt"
)

var UnrecognisedIP error = errors.New("Unrecognised IP Address")
var UnrecognisedSessionId = errors.New("Unrecognised Session Id")

type User struct {
	Id                   int64
	Username             string
	Password             string
	RecoveryEmailAddress string
}

type Authentication struct {
	SessionId string
	UserId    int64
	IpAddress string
}

type Authenticator struct {
	db *AuthenticationDatabase
}

func NewAuthenticator() *Authenticator {
	auth := new(Authenticator)
	auth.db = NewAuthenticationDatabase()
	return auth
}

func (auth *Authenticator) CreateUser(username, emailAddress, password string, sessionId, ipAddress string) (user *User, err error) {
	encrypted, err := encryptPassword(password)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	userId, err := auth.db.CreateUser(sessionId, ipAddress, username, emailAddress, encrypted)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	user = new(User)
	user.Id = userId
	user.Username = username
	user.Password = encrypted
	user.RecoveryEmailAddress = emailAddress

	return user, err
}

var ErrInvalidUsernamePassword error = errors.New("Invalid username and password combination")

func (auth *Authenticator) Login(username, password, ipAddress, sessionId string) (*User, error) {
	user, err := auth.db.GetUserByUsername(username)
	if err != nil {
		return nil, ErrInvalidUsernamePassword
	}

	match := comparePasswords(user.Password, password)
	if !match {
		return nil, ErrInvalidUsernamePassword
	}
	return user, nil
}

func (auth *Authenticator) GetAuthentication(sessionId, ipAddress string) (authentication *Authentication, user *User, err error) {
	authentication, user, err = auth.db.GetAuth(sessionId)

	if err == sql.ErrNoRows {
		return nil, nil, UnrecognisedSessionId
	} else if err != nil {
		return nil, nil, err
	}

	if authentication.IpAddress != ipAddress {
		return nil, nil, UnrecognisedIP
	}

	return authentication, user, nil
}

func comparePasswords(hashedPassword, textPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(textPassword))
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		fmt.Printf("compare Passwords error: %s\n", err)
	}
	return err == nil
}

func encryptPassword(password string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Printf("encrypt Password, err: %v\n", err)
	return string(enc), err
}
