package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/auth/domain/repository"
	"github.com/budimanlai/go-core/auth/domain/usecase"
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-core/auth/service"
	"github.com/budimanlai/go-core/base"
	"gorm.io/gorm"

	pkg_helpers "github.com/budimanlai/go-pkg/helpers"
)

type OtpConfig struct {
	UserInitiated      bool
	BotPhoneNumber     string
	CommandPrefix      string
	MaxPendingRequests int
	ExpiredDuration    time.Duration
}

type OtpUsecaseImpl struct {
	base.BaseUsecase[entity.Otp]

	config OtpConfig
	sender service.OtpSenderService
}

func NewOtpUsecaseImpl(db *gorm.DB, repo repository.OtpRepository, config OtpConfig) usecase.OtpUsecase {
	return &OtpUsecaseImpl{
		BaseUsecase: base.NewBaseUsecase(repo, db),
		config:      config,
	}
}

func (uc *OtpUsecaseImpl) SetSender(sender service.OtpSenderService) {
	uc.sender = sender
}

func (uc *OtpUsecaseImpl) IsUserInitiated() bool {
	return uc.config.UserInitiated
}

func (uc *OtpUsecaseImpl) SetUserInitiated(userInitiated bool) {
	uc.config.UserInitiated = userInitiated
}

// GenerateCommandCode generates the command code for WhatsApp bot based on transaction ID and pin code.
// The format of the command code is: trx_id_pin_code
//
// Parameters:
//   - trx_id: The transaction ID associated with the OTP
//   - pin_code: The generated pin code for the OTP
//
// Returns:
//   - string: The generated command code in the format trx_id_pin_code
func (uc *OtpUsecaseImpl) GenerateCommandCode(trx_id, pin_code string) string {
	return fmt.Sprintf("%s_%s", trx_id, pin_code)
}

// ParseCommandCode parses the command code received from WhatsApp bot to extract transaction ID and pin code.
// The expected format of the command code is: trx_id_pin_code
//
// Parameters:
//   - commandCode: The command code string to be parsed
//
// Returns:
//   - string: The extracted transaction ID
//   - string: The extracted pin code
//   - error: An error object if the parsing fails, otherwise nil
func (uc *OtpUsecaseImpl) ParseCommandCode(commandCode string) (string, string, error) {
	var trx_id, pin_code string
	n, err := fmt.Sscanf(commandCode, "%[^_]_%s", &trx_id, &pin_code)
	if err != nil {
		return "", "", err
	}
	if n != 2 {
		return "", "", errors.New("invalid command code format")
	}
	return trx_id, pin_code, nil
}

// GenerateOTP generates a WhatsApp OTP for the given phone number and transaction ID.
// It creates a new WAOtp entity with a randomly generated 6-digit pin code and saves it to the repository.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and deadlines
//   - phoneNumber: The phone number to which the OTP will be sent
//   - trx_id: The transaction ID associated with the OTP request
//
// Returns:
//   - string: The generated 6-digit pin code
//   - error: An error object if the operation fails, otherwise nil
func (uc *OtpUsecaseImpl) GenerateOTP(ctx context.Context, request dto.OtpRequest) (*dto.OtpResponse, error) {
	// 1. validate request.Identifier
	if request.Channel == "email" {
		if !pkg_helpers.IsValidEmail(request.Identifier) {
			return nil, errors.New("invalid email format")
		}
	} else {
		if !pkg_helpers.IsValidPhoneNumber(request.Identifier) {
			return nil, errors.New("invalid phone number format")
		}
	}

	// 2. validate trx_id uniqueness
	existingOtp, err := uc.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("handphone = ? and trx_id = ?", request.Identifier, request.TrxID)
	})
	if err != nil {
		return nil, err
	}
	if existingOtp != nil {
		return nil, errors.New("OTP with the given trx_id already exists")
	}

	// 3. check max pending requests in same day
	if uc.config.MaxPendingRequests > 0 {
		pendingCount, err := uc.Count(ctx, func(d *gorm.DB) *gorm.DB {
			return d.Where("handphone = ? and status = ? and created_at >= ?", request.Identifier, "pending", time.Now().Truncate(24*time.Hour))
		})
		if err != nil {
			return nil, err
		}
		if pendingCount >= int64(uc.config.MaxPendingRequests) {
			return nil, errors.New("maximum pending OTP requests reached")
		}
	}

	// 4. generate OTP
	pin_code := pkg_helpers.GenerateRandomNumberString(6)

	et := entity.Otp{
		Handphone: request.Identifier,
		TrxID:     request.TrxID,
		Status:    "pending",
		PinCode:   pin_code,
		CreatedAt: time.Now(),
	}

	err = uc.Create(ctx, &et)
	if err != nil {
		return nil, err
	}

	var out dto.OtpResponse = dto.OtpResponse{
		Identifier: request.Identifier,
		TrxID:      request.TrxID,
	}

	// 5. generate WhatsApp URL if user initiated
	// user initiated meaning is the user will send the OTP request via WhatsApp bot
	if uc.IsUserInitiated() {
		// example: OTP#trxid_pincode
		out.WaUrl = "https://wa.me/" + uc.config.BotPhoneNumber + "?text=" + uc.GenerateCommandCode(request.TrxID, pin_code)
	} else {
		out.WaUrl = ""

		if uc.sender != nil {
			// send OTP in background job
			err = uc.sender.Send(request.Channel, request.Identifier, pin_code)
			if err != nil {
				return nil, err
			}
		}
	}

	return &out, nil
}

// IsValid checks if the provided pin code is valid for the given phone number and transaction ID.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and deadlines
//   - phoneNumber: The phone number associated with the OTP
//   - trx_id: The transaction ID associated with the OTP
//   - pin_code: The pin code to be validated
//
// Returns:
//   - bool: True if the pin code is valid, false otherwise
//   - error: An error object if the operation fails, otherwise nil
func (uc *OtpUsecaseImpl) Status(ctx context.Context, phoneNumber string, trx_id string) (bool, error) {
	otp, err := uc.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("handphone = ? and trx_id = ?", phoneNumber, trx_id)
	})
	if err != nil {
		return false, err
	}

	// if otp not found, return invalid
	if otp == nil {
		return false, nil
	}

	var valid bool = false
	if otp.Status == "verified" {
		valid = true
	}

	return valid, nil
}

// VerifyOtp verifies the OTP for the given phone number, transaction ID, and pin code.
// It updates the status of the OTP to "verified" if the provided details are valid.
//
// Parameters:
//   - ctx: Context for managing request-scoped values and deadlines
//   - phoneNumber: The phone number associated with the OTP
//   - trx_id: The transaction ID associated with the OTP
//   - pin_code: The pin code to be verified
//
// Returns:
//   - error: An error object if the operation fails, otherwise nil
func (uc *OtpUsecaseImpl) VerifyOtp(ctx context.Context, phoneNumber, trx_id, pin_code string) error {
	// 1. find OTP by phone number, trx_id, pin_code, status = pending and not expired
	otp, err := uc.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("handphone = ? and trx_id = ? and pin_code = ? and created_at >= ?",
			phoneNumber, trx_id, pin_code, time.Now().Add(-uc.config.ExpiredDuration))
	})
	if err != nil {
		return err
	}

	// if otp not found, return invalid
	if otp == nil {
		return errors.New("Invalid OTP or expired")
	}

	if otp.Status == "verified" {
		return nil
	}

	// update status to verified
	otp.Status = "verified"
	err = uc.Update(ctx, otp)
	if err != nil {
		return err
	}

	return nil
}
