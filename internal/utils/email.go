package utils

import (
	"fmt"
	"html"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// Forgot password email
func ForgotPasswordMail(email, otpID string) error {

	urlHost := os.Getenv("APP_HOST")
	urlPort := "5173"

	// Send to email
	m := gomail.NewMessage()
	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password")

	resetLink := fmt.Sprintf(
		"http://%s:%s/reset-password?token=%s",
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
func VerifyEmailMail(otp, email string) error {
	// Send to email
	m := gomail.NewMessage()
	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verifify Email")

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
`, otp))

	// SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, smtpPassword)

	return d.DialAndSend(m)

}

// Email to admin for verify the organizer's profile
func OrganizerProfileAdminNotificationEmail(
	adminEmail string,
	organizerName string,
	userEmail string,
	actionType string, // "registered" or "updated"
) error {

	m := gomail.NewMessage()

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", adminEmail)
	m.SetHeader("Subject", "Organizer Profile Requires Approval")

	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Organizer Profile Approval Needed</title>
	</head>
	<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
			<tr>
				<td align="center">
					<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
						
						<tr>
							<td style="background:#F59E0B; padding:20px; text-align:center;">
								<h1 style="color:#ffffff; margin:0;">Approval Required ⚠️</h1>
							</td>
						</tr>

						<tr>
							<td style="padding:30px; color:#333333;">
								<p style="font-size:16px;">
									Hello Admin,
								</p>

								<p style="font-size:16px; line-height:1.6;">
									An organizer has <strong>%s</strong> their profile and is waiting for approval.
								</p>

								<div style="background:#FEF3C7; padding:15px; border-radius:6px; margin:20px 0;">
									<p style="margin:0; font-size:14px;">
										<strong>Organizer Name:</strong> %s<br>
										<strong>User Email:</strong> %s
									</p>
								</div>

								<div style="text-align:center; margin:30px 0;">
									<a href="https://ngevent.id/admin/organizers"
										style="
											background:#F59E0B;
											color:#ffffff;
											text-decoration:none;
											padding:12px 24px;
											border-radius:6px;
											font-size:16px;
											display:inline-block;
										">
										Review Organizer Profile
									</a>
								</div>

								<p style="font-size:14px; color:#777777;">
									Please review and take the appropriate action.
								</p>

								<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

								<p style="font-size:14px; color:#777777;">
									Ngevent System Notification
								</p>
							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, actionType, organizerName, userEmail))

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

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

// Organizer profile rejected
func OrganizerProfileRejectedEmail(email, organizerName, rejectedReason string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Organizer Profile Needs Revision")

	// Escape reason to prevent HTML injection
	safeReason := html.EscapeString(rejectedReason)

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
								<h1 style="color:#ffffff; margin:0;">Profile Requires Revision</h1>
							</td>
						</tr>

						<tr>
							<td style="padding:30px; color:#333333;">
								<p style="font-size:16px; line-height:1.6;">
									Hello <strong>%s</strong>,
								</p>

								<p style="font-size:16px; line-height:1.6;">
									Thank you for registering as an <strong>Event Organizer</strong> on <strong>Ngevent</strong>.
								</p>

								<p style="font-size:16px; line-height:1.6;">
									After reviewing your application, we found that your organizer profile
									<strong>cannot be approved at this time</strong>.
								</p>

								<div style="background:#FEE2E2; padding:15px; border-radius:6px; margin:20px 0;">
									<p style="margin:0; font-size:14px; color:#991B1B;">
										<strong>Reason:</strong><br>
										%s
									</p>
								</div>

								<p style="font-size:16px; line-height:1.6;">
									Please update your organizer profile and ensure that all submitted
									information is accurate and complete.
								</p>

								<div style="text-align:center; margin:30px 0;">
									<a href="https://ngevent.id/organizer/profile"
										style="
											background:#EF4444;
											color:#ffffff;
											text-decoration:none;
											padding:12px 24px;
											border-radius:6px;
											font-size:16px;
											display:inline-block;
										">
										Update Organizer Profile
									</a>
								</div>

								<p style="font-size:14px; color:#555555;">
									If you need assistance, please contact our support team.
								</p>

								<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

								<p style="font-size:14px; color:#777777;">
									Thank you for your understanding,<br>
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
	`, organizerName, safeReason))

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}

// ========== Event email ===========
// Admin notification
func AdminEventNotification(email, organizerName, eventName, EOEmail, status string) error {
	m := gomail.NewMessage()

	var subject string
	var title string
	var message string
	var color string
	var boxColor string

	switch status {
	case "create":
		subject = "New Event Requires Review"
		title = "New Event Submission"
		color = "#2563EB"
		boxColor = "#EFF6FF"
		message = "A new event has been submitted and is currently waiting for your review."

	case "update":
		subject = "Event Update Requires Review"
		title = "Event Updated"
		color = "#F59E0B"
		boxColor = "#FEF3C7"
		message = "An existing event has been updated by the organizer and requires your review."

	default:
		subject = "Event Notification"
		title = "Event Notification"
		color = "#6B7280"
		boxColor = "#F3F4F6"
		message = "There is an update regarding an event."
	}

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>%s</title>
	</head>
	<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
			<tr>
				<td align="center">
					<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">

						<tr>
							<td style="background:%s; padding:20px; text-align:center;">
								<h1 style="color:#ffffff; margin:0;">%s</h1>
							</td>
						</tr>

						<tr>
							<td style="padding:30px; color:#333333;">

								<p style="font-size:16px; line-height:1.6;">
									Hello Admin,
								</p>

								<p style="font-size:16px; line-height:1.6;">
									%s
								</p>

								<div style="background:%s; padding:15px; border-radius:6px; margin:20px 0;">
									<p style="margin:5px 0; font-size:14px;">
										<strong>Event Name:</strong> %s
									</p>
									<p style="margin:5px 0; font-size:14px;">
										<strong>Organizer Name:</strong> %s
									</p>
									<p style="margin:5px 0; font-size:14px;">
										<strong>Organizer Email:</strong> %s
									</p>
								</div>

								<p style="font-size:16px; line-height:1.6;">
									Please review the event details and approve or reject it from the admin dashboard.
								</p>

								<div style="text-align:center; margin:30px 0;">
									<a href="https://ngevent.id/admin/events"
										style="
											background:%s;
											color:#ffffff;
											text-decoration:none;
											padding:12px 24px;
											border-radius:6px;
											font-size:16px;
											display:inline-block;
										">
										Review Event
									</a>
								</div>

								<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

								<p style="font-size:14px; color:#777777;">
									This notification was sent automatically by the Ngevent system.
								</p>

							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, title, color, title, message, boxColor, eventName, organizerName, EOEmail, color))

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}

// Admin req update event notification
func AdminUpdatedEventNotification(email, organizerName, eventName, EOEmail string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "New Event Requires Review")

	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>New Event Review</title>
	</head>
	<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
			<tr>
				<td align="center">
					<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">

						<tr>
							<td style="background:#2563EB; padding:20px; text-align:center;">
								<h1 style="color:#ffffff; margin:0;">New Event Submission</h1>
							</td>
						</tr>

						<tr>
							<td style="padding:30px; color:#333333;">

								<p style="font-size:16px; line-height:1.6;">
									Hello Admin,
								</p>

								<p style="font-size:16px; line-height:1.6;">
									A new event has been submitted and is currently waiting for your review.
								</p>

								<div style="background:#EFF6FF; padding:15px; border-radius:6px; margin:20px 0;">
									<p style="margin:5px 0; font-size:14px;">
										<strong>Event Name:</strong> %s
									</p>
									<p style="margin:5px 0; font-size:14px;">
										<strong>Organizer Name:</strong> %s
									</p>
									<p style="margin:5px 0; font-size:14px;">
										<strong>Organizer Email:</strong> %s
									</p>
								</div>

								<p style="font-size:16px; line-height:1.6;">
									Please review the event details and approve or reject it from the admin dashboard.
								</p>

								<div style="text-align:center; margin:30px 0;">
									<a href="https://ngevent.id/admin/events"
										style="
											background:#2563EB;
											color:#ffffff;
											text-decoration:none;
											padding:12px 24px;
											border-radius:6px;
											font-size:16px;
											display:inline-block;
										">
										Review Event
									</a>
								</div>

								<hr style="border:none; border-top:1px solid #eeeeee; margin:30px 0;">

								<p style="font-size:14px; color:#777777;">
									This notification was sent automatically by the Ngevent system.
								</p>

							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, eventName, organizerName, EOEmail))

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}

// EO notification
func OrganizerEventNotification(email, organizerName, eventName string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Event Has Been Submitted")

	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Event Submitted</title>
	</head>
	<body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
			<tr>
				<td align="center">
					<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
						
						<tr>
							<td style="background:#22C55E; padding:20px; text-align:center;">
								<h1 style="color:#ffffff; margin:0;">Event Successfully Submitted</h1>
							</td>
						</tr>

						<tr>
							<td style="padding:30px; color:#333333;">
								<p style="font-size:16px; line-height:1.6;">
									Hello <strong>%s</strong>,
								</p>

								<p style="font-size:16px; line-height:1.6;">
									Your event has been successfully created on <strong>Ngevent</strong>.
								</p>

								<div style="background:#ECFDF5; padding:15px; border-radius:6px; margin:20px 0;">
									<p style="margin:0; font-size:15px;">
										<strong>Event Name:</strong> %s
									</p>
								</div>

								<p style="font-size:16px; line-height:1.6;">
									Your event is currently <strong>pending review</strong> by our admin team.
									Once the review process is completed, you will receive another notification
									regarding the approval status of your event.
								</p>

								<div style="text-align:center; margin:30px 0;">
									<a href="https://ngevent.id/organizer/events"
										style="
											background:#22C55E;
											color:#ffffff;
											text-decoration:none;
											padding:12px 24px;
											border-radius:6px;
											font-size:16px;
											display:inline-block;
										">
										View My Events
									</a>
								</div>

								<p style="font-size:14px; color:#555555;">
									Thank you for using <strong>Ngevent</strong> to manage your events.
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
	`, organizerName, eventName))

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}

// EO verification notification
func OrganizerEventVerification(email, organizerName, eventName, status, reason string) error {
	m := gomail.NewMessage()

	var subject string
	var title string
	var message string
	var color string
	var boxColor string

	switch status {
	case "active":
		subject = "Your Event Has Been Approved"
		title = "Event Approved"
		color = "#22C55E"
		boxColor = "#ECFDF5"
		message = "Congratulations! Your event has been approved by our admin team and is now live on Ngevent."

	case "rejected":
		subject = "Your Event Has Been Rejected"
		title = "Event Rejected"
		color = "#EF4444"
		boxColor = "#FEE2E2"
		message = "Unfortunately, your event submission could not be approved by our admin team. Please review your event details and make the necessary updates."

	default:
		subject = "Event Status Update"
		title = "Event Status Updated"
		color = "#6B7280"
		boxColor = "#F3F4F6"
		message = "Your event status has been updated."
	}

	reasonBlock := ""
	if status == "rejected" && reason != "" {
		reasonBlock = fmt.Sprintf(`
            <div style="background:#FEF2F2; border-left:4px solid #EF4444; padding:15px; border-radius:0 6px 6px 0; margin:20px 0;">
                <p style="margin:0 0 6px 0; font-size:13px; font-weight:bold; color:#991B1B; text-transform:uppercase; letter-spacing:0.05em;">
                    Reason for Rejection
                </p>
                <p style="margin:0; font-size:15px; color:#7F1D1D; line-height:1.6;">
                    %s
                </p>
            </div>
        `, reason)
	}

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>%s</title>
    </head>
    <body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
        <table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
            <tr>
                <td align="center">
                    <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
                        
                        <tr>
                            <td style="background:%s; padding:20px; text-align:center;">
                                <h1 style="color:#ffffff; margin:0;">%s</h1>
                            </td>
                        </tr>

                        <tr>
                            <td style="padding:30px; color:#333333;">
                                <p style="font-size:16px; line-height:1.6;">
                                    Hello <strong>%s</strong>,
                                </p>

                                <p style="font-size:16px; line-height:1.6;">
                                    %s
                                </p>

                                <div style="background:%s; padding:15px; border-radius:6px; margin:20px 0;">
                                    <p style="margin:0; font-size:15px;">
                                        <strong>Event Name:</strong> %s
                                    </p>
                                </div>

                                %s

                                <div style="text-align:center; margin:30px 0;">
                                    <a href="https://ngevent.id/organizer/events"
                                        style="
                                            background:%s;
                                            color:#ffffff;
                                            text-decoration:none;
                                            padding:12px 24px;
                                            border-radius:6px;
                                            font-size:16px;
                                            display:inline-block;
                                        ">
                                        View My Events
                                    </a>
                                </div>

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
	`, title, color, title, organizerName, message, boxColor, eventName, reasonBlock, color))

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}

func OrganizerUpdatedEventNotif(email, organizerName, eventName, status, reason string) error {
	m := gomail.NewMessage()

	var subject string
	var title string
	var message string
	var color string
	var boxColor string

	switch status {

	case "approved":
		subject = "Your Event Update Has Been Approved"
		title = "Event Update Approved"
		color = "#22C55E"
		boxColor = "#ECFDF5"
		message = "Your event update has been successfully reviewed and approved. The latest changes are now live on Ngevent."

	case "rejected":
		subject = "Your Event Update Has Been Rejected"
		title = "Event Update Rejected"
		color = "#EF4444"
		boxColor = "#FEE2E2"
		message = "Your recent event update could not be approved. Please update your event accordingly."
	default:
		subject = "Event Status Update"
		title = "Event Status Updated"
		color = "#6B7280"
		boxColor = "#F3F4F6"
		message = "Your event status has been updated."
	}

	reasonBlock := ""
	if status == "rejected" && reason != "" {
		reasonBlock = fmt.Sprintf(`
            <div style="background:#FEF2F2; border-left:4px solid #EF4444; padding:15px; border-radius:0 6px 6px 0; margin:20px 0;">
                <p style="margin:0 0 6px 0; font-size:13px; font-weight:bold; color:#991B1B; text-transform:uppercase; letter-spacing:0.05em;">
                    Reason for Rejection
                </p>
                <p style="margin:0; font-size:15px; color:#7F1D1D; line-height:1.6;">
                    %s
                </p>
            </div>
        `, reason)
	}

	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>%s</title>
    </head>
    <body style="margin:0; padding:0; background-color:#f4f4f4; font-family:Arial, Helvetica, sans-serif;">
        <table width="100%%" cellpadding="0" cellspacing="0" style="padding:20px;">
            <tr>
                <td align="center">
                    <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; overflow:hidden;">
                        
                        <tr>
                            <td style="background:%s; padding:20px; text-align:center;">
                                <h1 style="color:#ffffff; margin:0;">%s</h1>
                            </td>
                        </tr>

                        <tr>
                            <td style="padding:30px; color:#333333;">
                                <p style="font-size:16px;">Hello <strong>%s</strong>,</p>

                                <p style="font-size:16px; line-height:1.6;">
                                    %s
                                </p>

                                <div style="background:%s; padding:15px; border-radius:6px; margin:20px 0;">
                                    <p style="margin:0;">
                                        <strong>Event Name:</strong> %s
                                    </p>
                                </div>

                                %s

                                <div style="text-align:center; margin:30px 0;">
                                    <a href="https://ngevent.id/organizer/events"
                                        style="
                                            background:%s;
                                            color:#ffffff;
                                            text-decoration:none;
                                            padding:12px 24px;
                                            border-radius:6px;
                                            font-size:16px;
                                            display:inline-block;
                                        ">
                                        View My Events
                                    </a>
                                </div>

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
    `, title, color, title, organizerName, message, boxColor, eventName, reasonBlock, color))

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}
