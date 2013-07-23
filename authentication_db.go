package mvc

import (
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
)

type User struct {
	Id                   int
	Username             string
	Password             string
	RecoveryEmailAddress string
}

type Authentication struct {
	SessionId string
	UserId    int
	IpAddress string
}

type AuthenticationDatabase struct {
	db                           *sql.DB
	insertUser                   *sql.Stmt
	insertAuthentication         *sql.Stmt
	getUserById                  *sql.Stmt
	getUserByUsername            *sql.Stmt
	getUserByUsernameAndPassword *sql.Stmt
	getAuthenticationBySessionId *sql.Stmt
}

func (auth *AuthenticationDatabase) panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewAuthenticationDatabase() *AuthenticationDatabase {
	auth := new(AuthenticationDatabase)
	db, err := sql.Open("mymysql", "gomvc/gomvc/gomvc")
	auth.panicOnError(err)
	auth.db = db

	insUser, err := db.Prepare(insertUserSQL)
	auth.panicOnError(err)
	auth.insertUser = insUser

	insAuth, err := db.Prepare(insertAuthenticationSQL)
	auth.panicOnError(err)
	auth.insertAuthentication = insAuth

	usrById, err := db.Prepare(getUserByIdSQL)
	auth.panicOnError(err)
	auth.getUserById = usrById

	usrByUsername, err := db.Prepare(getUserByUsernameSQL)
	auth.panicOnError(err)
	auth.getUserByUsername = usrByUsername

	authBySession, err := db.Prepare(getAuthenticationBySessionId)
	auth.panicOnError(err)
	auth.getAuthenticationBySessionId = authBySession

	return auth
}

func (auth *AuthenticationDatabase) CreateUser(sessionId, ipAddress, username, emailAddress, encryptedPassword string) (userId int, err error) {

	res, err := auth.insertUser.Exec(username, encryptedPassword, emailAddress)

	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	id64, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	userId = int(id64)

	res, err = auth.insertAuthentication.Exec(sessionId, userId, ipAddress)

	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	return userId, nil
}

func (auth *AuthenticationDatabase) GetUserById(id int) (u *User, err error) {
	res := auth.getUserById.QueryRow(id)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	usr := new(User)
	err = res.Scan(
		&usr.Id,
		&usr.Username,
		&usr.Password,
		&usr.RecoveryEmailAddress,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return usr, nil
}

func (auth *AuthenticationDatabase) GetUserByUsername(username string) (user *User, err error) {
	res := auth.getUserByUsername.QueryRow(username)

	usr := new(User)
	err = res.Scan(
		&usr.Id,
		&usr.Username,
		&usr.Password,
		&usr.RecoveryEmailAddress,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return usr, nil
}

func (auth *AuthenticationDatabase) GetUserByUsernameAndPassword(username, password string) (user *User, err error) {
	res := auth.getUserByUsername.QueryRow(username, password)

	usr := new(User)
	err = res.Scan(
		&usr.Id,
		&usr.Username,
		&usr.Password,
		&usr.RecoveryEmailAddress,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return usr, nil
}

func (auth *AuthenticationDatabase) GetAuth(sessionId string) (authentication *Authentication, user *User, err error) {
	res := auth.getAuthenticationBySessionId.QueryRow(sessionId)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	a := new(Authentication)
	err = res.Scan(
		&a.SessionId,
		&a.UserId,
		&a.IpAddress,
	)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	usr, err := auth.GetUserById(a.UserId)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return a, usr, nil
}

var insertUserSQL string = `
INSERT INTO gomvc.User (
  Username,
  Password,
  RecoveryEmailAddress
) VALUES (
  ?,
  ?,
  ?
);`

var insertAuthenticationSQL string = `
INSERT INTO gomvc.Authentication (
  SessionId,
  UserId,
  IPAddress
) VALUES (
  ?,
  ?,
  ?
);`

var getUserByIdSQL string = `
SELECT Id, Username, Password, RecoveryEmailAddress
FROM gomvc.User
WHERE Id = ?
;`

var getUserByUsernameSQL string = `
SELECT Id, Username, Password, RecoveryEmailAddress
FROM gomvc.User
WHERE Username = ?
;`

var getUserByUsernameAndPasswordSQL string = `
SELECT Id, Username, Password, RecoveryEmailAddress
FROM gomvc.User
WHERE Username = ?
AND Password = ?
;`

var getAuthenticationBySessionId string = `
SELECT SessionId, UserId, IPAddress
FROM gomvc.Authentication
WHERE SessionId = ?
`
