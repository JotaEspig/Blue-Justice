package controllers

import (
	"net/http"
	"sigma/config"
	"sigma/models/student"
	"sigma/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddStudent changes the user type to "student" and adds a student row to the
// database if it doesn't already exists
func AddStudent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := getUsername(ctx)
		u := user.User{}
		err := config.DB.Select("id").Where("username = ?", username).First(&u).Error
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		err = config.DB.Transaction(func(tx *gorm.DB) error {
			u.Type = "student"
			err := tx.Updates(&u).Error
			if err != nil {
				return err
			}

			// Checks if a student already exists with this ID.
			// If it exists, it doesn't create another one
			// because the admin may have committed an error
			// when changing the type of the user
			s := &student.Student{}
			tx.First(s, u.ID)
			if s.UID != 0 {
				return nil
			}

			s, err = student.InitStudent(&u)
			if err != nil {
				return err
			}

			return tx.Create(s).Error
		})
	}
}

// GetStudentInfo gets student info according to the username
func GetStudentInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := getUsername(ctx)
		u := user.User{}
		s := student.Student{}
		err := config.DB.Select("id").Where("username = ?", username).First(&u).Error
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		err = config.DB.Preload("User").Where("id = ?", u.ID).First(&s).Error
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.JSON(
			http.StatusOK,
			gin.H{
				"student": s.ToMap(),
			},
		)
	}
}

// AutoMigrate the student table
func init() {
	config.DB.AutoMigrate(&student.Student{})
}
