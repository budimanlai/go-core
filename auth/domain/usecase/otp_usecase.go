package usecase

import (
	"context"

	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-core/auth/service"
	"github.com/budimanlai/go-core/base"
)

type OtpUsecase interface {
	base.BaseUsecase[entity.Otp]

	// GenerateOTP generates a WhatsApp OTP for the given phone number and transaction ID.
	GenerateOTP(ctx context.Context, request dto.OtpRequest) (*dto.OtpResponse, error)

	// VerifyOTP verifies the OTP pin code for the given phone number and transaction ID.
	Status(ctx context.Context, phoneNumber string, trx_id string) (bool, error)

	// VerifyOtp verifies the OTP pin code for the given phone number and transaction ID.
	VerifyOtp(ctx context.Context, phoneNumber, trx_id, pin_code string) error

	// SetSender sets the OTP sender service.
	SetSender(sender service.OtpSenderService)
}
