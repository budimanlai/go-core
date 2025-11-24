package dto_test

import (
	"testing"
	"time"

	"github.com/budimanlai/go-core/account/dto"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRequest_Validation(t *testing.T) {
	validate := validator.New()

	t.Run("Valid Request", func(t *testing.T) {
		dob := time.Now().AddDate(-20, 0, 0)
		addr := "123 Main St"
		req := dto.RegisterRequest{
			Username:      "johndoe",
			Email:         "john@example.com",
			Password:      "secret123",
			Fullname:      "John Doe",
			Handphone:     "081234567890",
			Dob:           &dob,
			Gender:        "M",
			Address:       &addr,
			Zipcode:       "12345",
			DistrictID:    1,
			SubdistrictID: 1,
			CityID:        1,
			ProvinceID:    1,
			CountryID:     "ID",
		}

		err := validate.Struct(req)
		assert.NoError(t, err)
	})

	t.Run("Invalid Request - Missing Required Fields", func(t *testing.T) {
		req := dto.RegisterRequest{}
		err := validate.Struct(req)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Contains(t, validationErrors.Error(), "Username")
		assert.Contains(t, validationErrors.Error(), "Email")
		assert.Contains(t, validationErrors.Error(), "Password")
		assert.Contains(t, validationErrors.Error(), "Fullname")
		assert.Contains(t, validationErrors.Error(), "Handphone")
		assert.Contains(t, validationErrors.Error(), "CountryID")
	})

	t.Run("Invalid Request - Invalid Email", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username:  "johndoe",
			Email:     "invalid-email",
			Password:  "secret123",
			Fullname:  "John Doe",
			Handphone: "081234567890",
			CountryID: "ID",
		}
		err := validate.Struct(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
	})

	t.Run("Invalid Request - Short Password", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username:  "johndoe",
			Email:     "john@example.com",
			Password:  "123",
			Fullname:  "John Doe",
			Handphone: "081234567890",
			CountryID: "ID",
		}
		err := validate.Struct(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Password")
	})

	t.Run("Invalid Request - Invalid Gender", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username:  "johndoe",
			Email:     "john@example.com",
			Password:  "secret123",
			Fullname:  "John Doe",
			Handphone: "081234567890",
			CountryID: "ID",
			Gender:    "X",
		}
		err := validate.Struct(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Gender")
	})
}

func TestLoginRequest_Validation(t *testing.T) {
	validate := validator.New()

	t.Run("Valid Request", func(t *testing.T) {
		req := dto.LoginRequest{
			Username: "johndoe",
			Password: "secret123",
		}
		err := validate.Struct(req)
		assert.NoError(t, err)
	})

	t.Run("Invalid Request - Missing Fields", func(t *testing.T) {
		req := dto.LoginRequest{}
		err := validate.Struct(req)
		assert.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)
		assert.Contains(t, validationErrors.Error(), "Username")
		assert.Contains(t, validationErrors.Error(), "Password")
	})
}

func TestUpdateUserRequest_Validation(t *testing.T) {
	validate := validator.New()

	t.Run("Valid Request - Partial Update", func(t *testing.T) {
		gender := "F"
		req := dto.UpdateUserRequest{
			Gender: &gender,
		}
		err := validate.Struct(req)
		assert.NoError(t, err)
	})

	t.Run("Invalid Request - Invalid Gender", func(t *testing.T) {
		gender := "X"
		req := dto.UpdateUserRequest{
			Gender: &gender,
		}
		err := validate.Struct(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Gender")
	})

	t.Run("Invalid Request - Invalid CountryID", func(t *testing.T) {
		countryID := "IND" // Too long
		req := dto.UpdateUserRequest{
			CountryID: &countryID,
		}
		err := validate.Struct(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CountryID")
	})
}
