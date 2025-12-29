package service

type OtpSenderService interface {
	// Send sends an OTP message to the specified recipient.
	Send(channel, to, pin_code string) error
}
