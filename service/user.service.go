package service

import (
	"github.com/mashingan/smapping"
	"haliscicek.com/gin-api/dto"
	"haliscicek.com/gin-api/model"
	"haliscicek.com/gin-api/repository"
	"log"
)

type UserService interface {
	Update(user dto.UserUpdateDto) model.User
	Profile(userId string) model.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) Update(user dto.UserUpdateDto) model.User {
	userToUpdate := model.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	updatedUser := service.userRepository.Update(userToUpdate)
	return updatedUser
}
func (service *userService) Profile(userId string) model.User {
	return service.userRepository.ProfileUser(userId)
}
