package user

import (
	"gorm.io/gorm"
)

func FetchUserByID(tx *gorm.DB, userID int64) (User, error) {
	user := User{}

	err := tx.Table("user_t").
		Where("id = ?", userID).
		Scan(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FetchUserByUsername(tx *gorm.DB, username string) (User, error) {
	user := User{}

	err := tx.Table("user_t").
		Where("username = ?", username).
		Scan(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}
