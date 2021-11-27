package repository

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"haliscicek.com/gin-api/model"
	"log"
)

type UserRepository interface {
	Save(user model.User) model.User
	Update(user model.User) model.User
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) model.User
	ProfileUser(userID string) model.User
}
type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) Save(user model.User) model.User {
	user.Password = hashAndSalt([]byte(user.Password))
	db.connection.Save(&user)
	return user
}

func (db *userConnection) Update(user model.User) model.User {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser model.User
		db.connection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}
	db.connection.Save(&user)
	return user
}

func (db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user model.User
	res := db.connection.Where("email = ?", email).Take(&user)
	
	if res.Error != nil {
		return nil
	}
	return user
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user model.User
	return db.connection.Where("email= ?", email).Take(&user)
}

func (db *userConnection) FindByEmail(email string) model.User {
	var user model.User
	db.connection.Where("email= ?", email).Take(&user)
	return user
}

func (db *userConnection) ProfileUser(userID string) model.User {
	var user model.User
	db.connection.Find(&user, userID)
	return user
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}
