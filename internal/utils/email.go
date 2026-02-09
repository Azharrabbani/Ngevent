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

// Organizer profile register
func OrganizerProfileVerificationEmail(email string, organizerName string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Profile is being verified")

	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Profil Organizer Sedang Diverifikasi</title>
	</head>
	<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
			<tr>
				<td align="center">
					<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
						
						<tr>
							<td style="background:#0EA5E9; padding:20px; text-align:center;">
								<h1 style="color:#ffffff; margin:0;">Ngevent 🎫</h1>
							</td>
						</tr>

						<tr>
							<td style="padding:30px; color:#333333;">
								<p style="font-size:16px; line-height:1.6;">
									Hello <strong>%s</strong> 👋,
								</p>

								<p style="font-size:16px; line-height:1.6;">
									Thank you for completing your <strong>Event Organizer profile</strong> on the <strong>Ngevent</strong> platform.
								</p>

								<p style="font-size:16px; line-height:1.6;">
									Currently, your profile is being <strong>verified by Ngevent Admin</strong>.
                                    This process aims to ensure that the organizer's data is valid and reliable.
								</p>

								<p style="font-size:16px; line-height:1.6;">
									After the verification process is complete, you will receive a follow-up notification email.
								</p>

								<p style="font-size:14px; color:#555555;">
									If you feel that you have never created or completed this organizer profile, please ignore this email.
								</p>

								<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

								<p style="font-size:14px; color:#777777;">
									Best regards,<br>
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
	`, organizerName))

	// SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}

// Organizer profile verified
func OrganizerProfileVerifiedEmail(email string, organizerName string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Profile Has Been Verified 🎉")

	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Profile Organizer Verified</title>
	</head>
	<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
			<tr>
				<td align="center">
					<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
						
						<tr>
							<td style="background:#22C55E; padding:20px; text-align:center;">
								<h1 style="color:#ffffff; margin:0;">Profil Verified ✅</h1>
							</td>
						</tr>

						<tr>
							<td style="padding:30px; color:#333333;">
								<p style="font-size:16px; line-height:1.6;">
									Hello <strong>%s</strong> 👋,
								</p>

								<p style="font-size:16px; line-height:1.6;">
									Good news! 🎉  
									Your <strong>Event Organizer</strong> profile on <strong>Ngevent</strong> has been successfully
                                    <strong>verified by the Admin</strong>.
								</p>

								<p style="font-size:16px; line-height:1.6;">
									Now you can:
								</p>

								<ul style="font-size:16px; line-height:1.8; padding-left:20px;">
									<li>Creating and managing events</li>
									<li>Publish event to public</li>
									<li>Accepting participants and transactions</li>
								</ul>

								<p style="font-size:16px; line-height:1.6;">
									Please log in to your organizer dashboard and start creating your first event 🚀
								</p>

								<div style="text-align:center; margin:30px 0;">
									<a href="https://ngevent.id/dashboard"
										style="
											background:#22C55E;
											color:#ffffff;
											text-decoration:none;
											padding:12px 24px;
											border-radius:6px;
											font-size:16px;
											display:inline-block;
										">
										Log in to the Dashboard
									</a>
								</div>

								<p style="font-size:14px; color:#555555;">
									If you have questions, please contact our support team.
								</p>

								<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

								<p style="font-size:14px; color:#777777;">
									See you at your event!<br>
									<strong>The Ngevent Team</strong>
								</p>
							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, organizerName))

	// SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}

// Organizer profile denied
func OrganizerProfileRejectedEmail(email string, organizerName string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Organizer Profile Rejected")

	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Organizer Profile Rejected</title>
	</head>
	<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
			<tr>
				<td align="center">
					<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
						
						<tr>
							<td style="background:#EF4444; padding:20px; text-align:center;">
								<h1 style="color:#ffffff; margin:0;">Profile Not Approved ❌</h1>
							</td>
						</tr>

						<tr>
							<td style="padding:30px; color:#333333;">
								<p style="font-size:16px; line-height:1.6;">
									Hello <strong>%s</strong> 👋,
								</p>

								<p style="font-size:16px; line-height:1.6;">
									Thank you for registering as an <strong>Event Organizer</strong> on <strong>Ngevent</strong>.
								</p>

								<p style="font-size:16px; line-height:1.6;">
									After reviewing your application, we would like to inform you that
									<strong>we are unable to verify your organizer profile</strong>.
								</p>

								<p style="font-size:16px; line-height:1.6;">
									This is usually caused by incomplete, incorrect data,
                                    or the data is not meet the organizer's verification requirements.
								</p>

								<p style="font-size:16px; line-height:1.6;">
									To continue, please <strong>re-register</strong> and ensure that all
                                    the data you have entered is correct and valid.
								</p>

								<div style="text-align:center; margin:30px 0;">
									<a href="https://ngevent.id/register-organizer"
										style="
											background:#EF4444;
											color:#ffffff;
											text-decoration:none;
											padding:12px 24px;
											border-radius:6px;
											font-size:16px;
											display:inline-block;
										">
										Re-Register Organizer
									</a>
								</div>

								<p style="font-size:14px; color:#555555;">
									If you need assistance or further information,
                                    please contact our support team.
								</p>

								<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

								<p style="font-size:14px; color:#777777;">
									Thank you for your understanding,<br>
									<strong>Ngevent team</strong>
								</p>
							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, organizerName))

	// SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}
