package model

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/jinzhu/gorm"
	"github.com/sendgrid/sendgrid-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	DeletedAt      *core.NullableTimestamp `json:"deleted_at,omitempty" settable:"false"`
	LastLogin      *core.NullableTimestamp `json:"last_login,omitempty"`
	Name           string                  `json:"name" sql:"not null" valid:"length(3|255),required"`
	Username       string                  `json:"username" sql:"not null" valid:"length(3|15),matches(^[a-zA-Z0-9][a-zA-Z0-9-_]+$),required"`
	Email          string                  `json:"email" sql:"not null" valid:"email,required"`
	Password       string                  `json:"password,omitempty" sql:"-" valid:"length(5|64)"`
	HashedPassword []byte                  `json:"-" sql:"hashed_password;not null" gorm:"size:32"`
	Image          string                  `json:"image"`
	Phone          string                  `json:"phone"`
	Balance        float64                 `json:"balance" sql:"default:0"`
	Admin          bool                    `json:"admin"`
	Ban            bool                    `json:"ban"`
}

func (u User) ValidateForCreate() core.DefaultError {
	err := u.ValidateField("name")
	if err != nil {
		return err
	}
	err = u.ValidateField("username")
	if err != nil {
		return err
	}
	err = u.ValidateField("email")
	if err != nil {
		return err
	}
	err = u.ValidateField("password")
	if err != nil {
		return err
	}
	err = u.ValidateField("phone")
	if err != nil {
		return err
	}

	return nil
}

func (u User) ValidateForUpdate() core.DefaultError {
	err := u.ValidateField("name")
	if err != nil {
		return err
	}
	err = u.ValidateField("email")
	if err != nil {
		return err
	}
	err = u.ValidateField("phone")
	if err != nil {
		return err
	}
	return nil
}

func (u User) ValidateForDelete(ctx *ModelCtx) core.DefaultError {
	return nil
}

func (u User) ValidateField(f string) core.DefaultError {
	data := map[string]interface{}{"field": f}
	err := ValidateStructField(u, f)

	if err != nil {
		switch f {
		case "email":
			if strings.Contains(err.Error(), "validate as email") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_EMAIL_FORMAT, data)
			}
		case "name":
			if strings.Contains(err.Error(), "validate as matches") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_NAME_FORMAT, data)
			} else if strings.Contains(err.Error(), "validate as length") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_NAME_LENGTH, data)
			}
		case "username":
			if strings.Contains(err.Error(), "validate as matches") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_USERNAME_FORMAT, data)
			} else if strings.Contains(err.Error(), "validate as length") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_USERNAME_LENGTH, data)
			}
		case "phone":
			if strings.Contains(err.Error(), "validate as matches") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_PHONE_TAKEN, data)
			} else if strings.Contains(err.Error(), "validate as length") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_PHONE_LENGTH, data)
			}
		case "password":
			if strings.Contains(err.Error(), "validate as password") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_PASSWORD_FORMAT, data)
			} else if strings.Contains(err.Error(), "validate as length") {
				err = core.NewBusinessError(err.Error(), core.ERROR_SUBCODE_PASSWORD_LENGTH, data)
			}
		}
	}

	return err
}

// Restrictor

func (u User) UserCanView(ctx *ModelCtx, viewer User) (bool, core.DefaultError) {
	return true, nil
}

func (u User) UserCanCreate(ctx *ModelCtx, creator User) (bool, core.DefaultError) {
	return true, nil
}

func (u User) UserCanUpdate(ctx *ModelCtx, updater User, fields []string) (bool, core.DefaultError) {
	if u.ID == updater.ID {
		return true, nil
	}
	return false, nil
}

func (u User) UserCanDelete(ctx *ModelCtx, deleter User) (bool, core.DefaultError) {
	return deleter.IsAdmin(), nil
}

// Creator Interface
func (u *User) Create(ctx *ModelCtx, creator User) core.DefaultError {
	db := ctx.Database

	err := u.ValidateForCreate()
	if err != nil {
		return err
	}

	var existing User
	dberr := db.Scopes(ByUserEmail(u.Email)).First(&existing).Error
	if dberr == nil {
		data := map[string]interface{}{
			"creator_id": creator.ID,
			"name":       u.Name,
			"email":      u.Email,
		}
		return core.NewBusinessError(
			"email: already used;",
			core.ERROR_SUBCODE_EMAIL_TAKEN,
			data,
		)
	}

	dberr = db.Scopes(ByUserUsername(u.Username)).First(&existing).Error
	if dberr == nil {
		data := map[string]interface{}{
			"creator_id": creator.ID,
			"name":       u.Name,
			"email":      u.Email,
		}
		return core.NewBusinessError(
			"username: already used;",
			core.ERROR_SUBCODE_USERNAME_TAKEN,
			data,
		)
	}

	// User bcrypt to generate HashedPassword
	u.HashedPassword, dberr = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if dberr != nil {
		data := map[string]interface{}{
			"creator_id": creator.ID,
			"name":       u.Name,
			"email":      u.Email,
		}
		return core.NewServerError(dberr.Error(), data)
	}

	u.Password = ""

	dberr = db.Set("gorm:save_associations", false).Create(&u).Error
	if dberr != nil {
		data := map[string]interface{}{
			"creator_id": creator.ID,
			"name":       u.Name,
			"email":      u.Email,
		}
		return core.NewServerError(dberr.Error(), data)
	}

	return nil
}

func (u *User) Update(ctx *ModelCtx, creator User) core.DefaultError {
	db := ctx.Database

	if u.Email != "" {
		var existing User
		dberr := db.Scopes(ByUserEmail(u.Email)).First(&existing).Error
		if dberr == nil {
			data := map[string]interface{}{
				"creator_id": creator.ID,
				"name":       u.Name,
				"email":      u.Email,
			}
			return core.NewBusinessError(
				"email: already used;",
				core.ERROR_SUBCODE_EMAIL_TAKEN,
				data,
			)
		}
	}

	dberr := db.Set("gorm:save_associations", false).Save(&u).Error
	if dberr != nil {
		data := map[string]interface{}{
			"creator_id": creator.ID,
			"name":       u.Name,
			"email":      u.Email,
		}
		return core.NewServerError(dberr.Error(), data)
	}

	return nil
}

// Business methods

func (u User) IsUser() bool {
	return !u.Admin
}

func (u User) IsAdmin() bool {
	return u.Admin
}

func (u User) VerifyPassword(password string) (bool, error) {
	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password)) == nil, nil
}

// Scopes
func ByUserEmail(email string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(users.email) = LOWER(?)", email)
	}
}

func ByUserUsername(username string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(users.username) = LOWER(?)", username)
	}
}

// email
func (u User) ResetPasswordEmail(token string) {
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = []byte(`{
    "personalizations": [
			{
				"to": [
						{
								"email": "` + u.Email + `"
						}
				],
				"subject": "Forgot your password, ` + u.Name + `?",
				"dynamic_template_data": {
					"email":"` + u.Email + `",
					"name":"` + u.Name + `",
					"token":"` + token + `",
				}
			}
    ],
    "from": {
			"email": "` + os.Getenv("SENDGRID_USER") + `",
			"name": "Server"
    },
    "template_id" : "d-aedb512d0ca7460cadc0027d561d9c25"
	}`)

	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
