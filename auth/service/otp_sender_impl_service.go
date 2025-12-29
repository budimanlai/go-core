package service

import (
	"github.com/budimanlai/go-core/service"
)

type OtpSenderServiceImpl struct {
	EmailService service.SMTPMailService
	PhoneService service.WhatsAppService
}

func NewOtpSenderServiceImpl(mailService service.SMTPMailService, whatsappService service.WhatsAppService) OtpSenderService {
	return &OtpSenderServiceImpl{
		EmailService: mailService,
		PhoneService: whatsappService,
	}
}

// Send sends OTP via the specified channel (email or phone)
// to the given recipient with the provided pin code.
// It uses the EmailService for email channel and PhoneService for phone channel.
// If the channel is not recognized, it does nothing.
// Returns an error if any occurs during the sending process.
// Send OTP in background job
func (s *OtpSenderServiceImpl) Send(channel, to, pin_code string) error {
	if channel == "email" && s.EmailService != nil {
		return s.EmailService.SendWithTemplate(to, "otp_notification", map[string]interface{}{
			"pin_code": pin_code,
		})
	} else if channel == "phone" && s.PhoneService != nil {
		return s.PhoneService.SendWithTemplate(to, "otp_notification", map[string]interface{}{
			"pin_code": pin_code,
		})
	}

	return nil
}
