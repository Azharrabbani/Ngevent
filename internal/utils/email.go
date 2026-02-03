package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// Forgot password email
func ForgotPasswordMail(email, otpID string) error {

	urlHost := os.Getenv("APP_HOST")
	urlPort := os.Getenv("APP_PORT")

	// Send to email
	m := gomail.NewMessage()
	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password")

	resetLink := fmt.Sprintf(
		"%s:%s/api/v1/reset-password/%s",
		urlHost,
		urlPort,
		otpID,
	)

	m.SetBody("text/html", fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Reset Password</title>
		</head>
		<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
			<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
				<tr>
					<td align="center">
						<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
							
							<tr>
								<td style="background:#00D9FF; padding:20px; text-align:center;">
									<h1 style="color:#ffffff; margin:0;">Reset Password</h1>
								</td>
							</tr>

							<tr>
								<td style="padding:30px; color:#333333;">
									<p style="font-size:16px; line-height:1.6;">
										Hello👋,
									</p>
									<p style="font-size:16px; line-height:1.6;">
										We have received your request to reset your account password.
										Please click the button below to continue.
									</p>

									<div style="text-align:center; margin:30px 0;">
										<a href="%s"
										style="background:#00D9FF; color:#ffffff; text-decoration:none;
												padding:14px 24px; border-radius:6px; font-size:16px;
												display:inline-block;">
											Reset Password
										</a>
									</div>

									<p style="font-size:14px; color:#555555;">
										If you feel that you are not making this request, please ignore this email.
									</p>

									<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

									<p style="font-size:14px; color:#777777;">
										Best regards👏,<br>
										<strong>Ngevent Team</strong>
									</p>
								</td>
							</tr>

						</table>
					</td>
				</tr>
			</table>
		</body>
		</html>
`, resetLink))

	// SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, smtpPassword)
	return d.DialAndSend(m)
}

// Verify email
func VerifyEmailMail(otp, email, otpID string) error {
	urlHost := os.Getenv("APP_HOST")
	urlPort := os.Getenv("APP_PORT")

	// Send to email
	m := gomail.NewMessage()
	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verifify Email")

	verifyLink := fmt.Sprintf(
		"%s:%s/api/v1/verify-email/%s",
		urlHost,
		urlPort,
		otpID,
	)

	m.SetBody("text/html", fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Verify Your Email</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
	<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
		<tr>
			<td align="center">
				<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
					
					<tr>
						<td style="background:#00D9FF; padding:20px; text-align:center;">
							<h1 style="color:#ffffff; margin:0;">Verify Your Email</h1>
						</td>
					</tr>

					<tr>
						<td style="padding:30px; color:#333333;">
							<p style="font-size:16px; line-height:1.6;">
								Hello 👋,
							</p>

							<p style="font-size:16px; line-height:1.6;">
								Thank you for registering. Please use the verification code below
								to confirm your email address.
							</p>

							<!-- OTP CODE -->
							<div style="text-align:center; margin:24px 0;">
								<div style="
									display:inline-block;
									padding:14px 32px;
									font-size:24px;
									letter-spacing:6px;
									font-weight:bold;
									background:#f4f8fb;
									border:1px dashed #0084ffff;
									border-radius:8px;
									color:#333333;">
									%s
								</div>
							</div>

							<p style="font-size:14px; color:#555555; text-align:center;">
								This verification code will expire shortly.
							</p>

							<!-- BUTTON -->
							<div style="text-align:center; margin:30px 0;">
								<a href="%s"
								   style="background:#00D9FF; color:#ffffff; text-decoration:none;
										  padding:14px 24px; border-radius:6px; font-size:16px;
										  display:inline-block;">
									Verify Email
								</a>
							</div>

							<p style="font-size:14px; color:#555555;">
								If you did not create an account, you can safely ignore this email.
							</p>

							<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

							<p style="font-size:14px; color:#777777;">
								Best regards 👏,<br>
								<strong>Ngevent Team</strong>
							</p>
						</td>
					</tr>

				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, otp, verifyLink))

	// SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, smtpPassword)

	return d.DialAndSend(m)

}
