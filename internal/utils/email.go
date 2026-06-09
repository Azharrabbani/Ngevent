package utils

import (
	"fmt"
	"html"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func emailWrapper(accentColor, headerTitle, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<title>%s</title>
<!--[if mso]>
<noscript><xml><o:OfficeDocumentSettings><o:PixelsPerInch>96</o:PixelsPerInch></o:OfficeDocumentSettings></xml></noscript>
<![endif]-->
<style>
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { background-color: #F0F4F8; font-family: 'Inter', Arial, sans-serif; -webkit-font-smoothing: antialiased; }
  a { color: inherit; }
  @media only screen and (max-width: 600px) {
    .email-container { width: 100%% !important; border-radius: 0 !important; }
    .email-body { padding: 24px 20px !important; }
    .btn { display: block !important; text-align: center !important; }
    .info-box { padding: 14px !important; }
  }
</style>
</head>
<body style="margin:0;padding:0;background-color:#F0F4F8;">

<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#F0F4F8;padding:40px 16px;">
  <tr>
    <td align="center">

      <!-- Container -->
      <table class="email-container" width="600" cellpadding="0" cellspacing="0" border="0"
        style="background:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 4px 24px rgba(0,0,0,0.08);max-width:600px;width:100%%;">

        <!-- Header -->
        <tr>
          <td style="background:%s;padding:36px 40px;text-align:center;">
            <h1 style="color:#ffffff;font-size:22px;font-weight:700;letter-spacing:-0.3px;margin:0;line-height:1.3;">
              %s
            </h1>
          </td>
        </tr>

        <!-- Body -->
        <tr>
          <td class="email-body" style="padding:36px 40px;color:#1E293B;">
            %s
          </td>
        </tr>

        <!-- Footer -->
        <tr>
          <td style="background:#F8FAFC;padding:24px 40px;border-top:1px solid #E2E8F0;text-align:center;">
            <p style="font-size:13px;color:#94A3B8;line-height:1.6;margin:0;">
              This email was sent automatically by <strong style="color:#64748B;">Ngevent</strong>.<br>
              If you did not expect this email, you can safely ignore it.
            </p>
            <p style="font-size:12px;color:#CBD5E1;margin-top:12px;">© 2025 Ngevent. All rights reserved.</p>
          </td>
        </tr>

      </table>
      <!-- /Container -->

    </td>
  </tr>
</table>

</body>
</html>`, headerTitle, accentColor, headerTitle, body)
}

func emailButton(color, link, label string) string {
	return fmt.Sprintf(`
<div style="text-align:center;margin:32px 0;">
  <a href="%s" class="btn"
    style="display:inline-block;background:%s;color:#ffffff;text-decoration:none;
           padding:14px 32px;border-radius:10px;font-size:15px;font-weight:600;
           letter-spacing:0.2px;line-height:1;">
    %s
  </a>
</div>`, link, color, label)
}

func emailInfoBox(bgColor, borderColor, content string) string {
	return fmt.Sprintf(`
<div class="info-box"
  style="background:%s;border:1px solid %s;border-radius:10px;padding:18px 20px;margin:24px 0;">
  %s
</div>`, bgColor, borderColor, content)
}

func emailDivider() string {
	return `<hr style="border:none;border-top:1px solid #E2E8F0;margin:28px 0;">`
}

func emailRow(label, value string) string {
	return fmt.Sprintf(`
<p style="margin:6px 0;font-size:14px;color:#475569;">
  <span style="color:#94A3B8;font-size:12px;text-transform:uppercase;letter-spacing:0.05em;font-weight:600;">%s</span><br>
  <span style="color:#1E293B;font-size:15px;font-weight:500;">%s</span>
</p>`, label, value)
}

func sendMail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, username, password)
	return d.DialAndSend(m)
}

func ForgotPasswordMail(email, otpID string) error {
	urlHost := os.Getenv("APP_HOST")
	resetLink := fmt.Sprintf("http://%s:5173/reset-password?token=%s", urlHost, otpID)

	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello 👋,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  We received a request to reset the password for your Ngevent account.
  Click the button below to choose a new password.
</p>

%s

%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  This link will expire in 1 hour. If you did not request a password reset, no action is needed — your account remains secure.
</p>`,
		emailButton("#6366F1", resetLink, "Reset My Password"),
		emailDivider(),
	)

	html := emailWrapper("#6366F1", "Reset Your Password", body)
	return sendMail(email, "Reset Your Password — Ngevent", html)
}

func VerifyEmailMail(otp, email string) error {
	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello 👋,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:24px;">
  Welcome to <strong style="color:#1E293B;">Ngevent</strong>! Use the verification code below to confirm your email address.
</p>

<!-- OTP Box -->
<div style="text-align:center;margin:28px 0;">
  <div style="display:inline-block;background:#F1F5F9;border:2px dashed #6366F1;
              border-radius:12px;padding:20px 40px;">
    <p style="font-size:11px;color:#94A3B8;text-transform:uppercase;letter-spacing:0.12em;font-weight:600;margin:0 0 8px 0;">
      Your verification code
    </p>
    <p style="font-size:36px;font-weight:700;letter-spacing:10px;color:#1E293B;margin:0;font-family:monospace;">
      %s
    </p>
  </div>
</div>

%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;text-align:center;">
  This code expires shortly. If you did not create an account, please ignore this email.
</p>`, otp, emailDivider())

	html := emailWrapper("#6366F1", "Verify Your Email", body)
	return sendMail(email, "Verify Your Email — Ngevent", html)
}

func OrganizerProfileAdminNotificationEmail(adminEmail, organizerName, userEmail, actionType string) error {
	verb := "registered"
	if actionType == "updated" {
		verb = "updated"
	}

	infoContent := emailRow("Organizer Name", organizerName) + emailRow("User Email", userEmail)
	infoBox := emailInfoBox("#FFFBEB", "#FDE68A", infoContent)

	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello Admin,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  An organizer has <strong style="color:#1E293B;">%s</strong> their profile and is awaiting your review and approval.
</p>

%s

%s
%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  Please log in to the admin dashboard to review and take the appropriate action.
</p>`,
		verb,
		infoBox,
		emailButton("#F59E0B", "https://ngevent.id/admin/organizers", "Review Organizer Profile"),
		emailDivider(),
	)

	html := emailWrapper("#F59E0B", "Approval Required", body)
	return sendMail(adminEmail, "Organizer Profile Requires Approval — Ngevent", html)
}

func OrganizerProfileVerificationEmail(email, organizerName string) error {
	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello <strong style="color:#1E293B;">%s</strong> 👋,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Thank you for completing your <strong style="color:#1E293B;">Event Organizer profile</strong> on Ngevent.
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  Your profile is currently being <strong style="color:#1E293B;">reviewed by our admin team</strong>.
  This process ensures that all organizer data is valid and trustworthy. You will receive a follow-up
  notification once the review is complete.
</p>

%s
%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  If you did not create this profile, please ignore this email.
</p>`,
		organizerName,
		emailDivider(),
		emailDivider(),
	)

	html := emailWrapper("#0EA5E9", "Profile Under Review", body)
	return sendMail(email, "Your Organizer Profile is Being Reviewed — Ngevent", html)
}

func OrganizerProfileVerifiedEmail(email, organizerName string) error {
	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello <strong style="color:#1E293B;">%s</strong> 👋,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:24px;">
  Great news! 🎉 Your <strong style="color:#1E293B;">Event Organizer</strong> profile has been
  <strong style="color:#22C55E;">verified</strong> by our admin team.
</p>

%s

<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:8px;">You can now:</p>
<ul style="font-size:15px;line-height:2;color:#475569;padding-left:20px;margin:0 0 24px 0;">
  <li>Create and manage events</li>
  <li>Publish events to the public</li>
  <li>Accept participants and transactions</li>
</ul>

%s
%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  If you have any questions, please reach out to our support team.
</p>`,
		organizerName,
		emailInfoBox("#F0FDF4", "#BBF7D0", `<p style="margin:0;font-size:14px;color:#166534;font-weight:500;">✅ Your profile is now active and visible on Ngevent.</p>`),
		emailButton("#22C55E", "https://ngevent.id/organizer/dashboard", "Go to Dashboard"),
		emailDivider(),
	)

	html := emailWrapper("#22C55E", "Profile Verified!", body)
	return sendMail(email, "Your Organizer Profile Has Been Verified 🎉 — Ngevent", html)
}

func OrganizerProfileRejectedEmail(email, organizerName, rejectedReason string) error {
	safeReason := html.EscapeString(rejectedReason)

	reasonBox := emailInfoBox("#FEF2F2", "#FECACA", fmt.Sprintf(`
<p style="margin:0 0 6px 0;font-size:12px;font-weight:600;color:#991B1B;text-transform:uppercase;letter-spacing:0.05em;">
  Reason for Rejection
</p>
<p style="margin:0;font-size:14px;color:#7F1D1D;line-height:1.6;">%s</p>`, safeReason))

	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello <strong style="color:#1E293B;">%s</strong>,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  Thank you for registering as an Event Organizer on Ngevent. After reviewing your application,
  we were unable to approve your profile at this time.
</p>

%s

<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  Please update your organizer profile based on the feedback above and resubmit for review.
</p>

%s
%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  Need help? Contact our support team and we'll be happy to assist.
</p>`,
		organizerName,
		reasonBox,
		emailButton("#EF4444", "https://ngevent.id/organizer/profile", "Update My Profile"),
		emailDivider(),
	)

	html := emailWrapper("#EF4444", "Profile Requires Revision", body)
	return sendMail(email, "Your Organizer Profile Needs Revision — Ngevent", html)
}

func AdminEventNotification(email, organizerName, eventName, EOEmail, status string) error {
	var subject, title, color string
	var message string

	switch status {
	case "create":
		subject = "New Event Submission Requires Review"
		title = "New Event Submitted"
		color = "#2563EB"
		message = "A new event has been submitted by an organizer and is awaiting your review."
	case "update":
		subject = "Event Update Requires Review"
		title = "Event Updated"
		color = "#F59E0B"
		message = "An existing event has been updated by the organizer and requires your review."
	default:
		subject = "Event Notification"
		title = "Event Notification"
		color = "#6B7280"
		message = "There is an update regarding an event on Ngevent."
	}

	infoContent := emailRow("Event Name", eventName) +
		emailRow("Organizer Name", organizerName) +
		emailRow("Organizer Email", EOEmail)
	infoBox := emailInfoBox("#EFF6FF", "#BFDBFE", infoContent)

	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello Admin,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  %s
</p>

%s

%s
%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  Please log in to the admin dashboard to review and approve or reject this event.
</p>`,
		message, infoBox,
		emailButton(color, "https://ngevent.id/admin/events", "Review Event"),
		emailDivider(),
	)

	htmlBody := emailWrapper(color, title, body)
	return sendMail(email, subject+" — Ngevent", htmlBody)
}

func AdminUpdatedEventNotification(email, organizerName, eventName, EOEmail string) error {
	return AdminEventNotification(email, organizerName, eventName, EOEmail, "update")
}

func OrganizerEventNotification(email, organizerName, eventName string) error {
	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello <strong style="color:#1E293B;">%s</strong>,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  Your event has been successfully submitted to <strong style="color:#1E293B;">Ngevent</strong>
  and is now pending review by our admin team.
</p>

%s

<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  Once the review is complete, you will receive a notification with the approval status.
  In the meantime, you can track your event from the dashboard.
</p>

%s
%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  Thank you for using Ngevent to manage your events.
</p>`,
		organizerName,
		emailInfoBox("#F0FDF4", "#BBF7D0", fmt.Sprintf(`
%s
<p style="margin:8px 0 0 0;font-size:13px;color:#166534;">Status: <strong>Pending Review</strong></p>`,
			emailRow("Event Name", eventName),
		)),
		emailButton("#22C55E", "https://ngevent.id/organizer/events", "View My Events"),
		emailDivider(),
	)

	htmlBody := emailWrapper("#22C55E", "Event Submitted!", body)
	return sendMail(email, "Your Event Has Been Submitted — Ngevent", htmlBody)
}

func OrganizerEventVerification(email, organizerName, eventName, status, reason string) error {
	var subject, title, color, message string

	switch status {
	case "active":
		subject = "Your Event Has Been Approved"
		title = "Event Approved!"
		color = "#22C55E"
		message = "Congratulations! Your event has been <strong style=\"color:#22C55E;\">approved</strong> by our admin team and is now live on Ngevent."
	case "rejected":
		subject = "Your Event Needs Revision"
		title = "Event Rejected"
		color = "#EF4444"
		message = "Unfortunately, your event submission could not be approved at this time. Please review the feedback below and update your event accordingly."
	default:
		subject = "Event Status Update"
		title = "Event Status Updated"
		color = "#6B7280"
		message = "Your event status has been updated."
	}

	reasonBlock := ""
	if status == "rejected" && reason != "" {
		reasonBlock = emailInfoBox("#FEF2F2", "#FECACA", fmt.Sprintf(`
<p style="margin:0 0 6px 0;font-size:12px;font-weight:600;color:#991B1B;text-transform:uppercase;letter-spacing:0.05em;">Reason for Rejection</p>
<p style="margin:0;font-size:14px;color:#7F1D1D;line-height:1.6;">%s</p>`, html.EscapeString(reason)))
	}

	eventBox := emailInfoBox("#F8FAFC", "#E2E8F0", emailRow("Event Name", eventName))

	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello <strong style="color:#1E293B;">%s</strong>,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  %s
</p>

%s
%s

%s
%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  If you have questions, please contact our support team.
</p>`,
		organizerName, message, eventBox, reasonBlock,
		emailButton(color, "https://ngevent.id/organizer/events", "View My Events"),
		emailDivider(),
	)

	htmlBody := emailWrapper(color, title, body)
	return sendMail(email, subject+" — Ngevent", htmlBody)
}

func OrganizerUpdatedEventNotif(email, organizerName, eventName, status, reason string) error {
	var subject, title, color, message string

	switch status {
	case "approved":
		subject = "Your Event Update Has Been Approved"
		title = "Update Approved!"
		color = "#22C55E"
		message = "Your event update has been <strong style=\"color:#22C55E;\">approved</strong> and the latest changes are now live on Ngevent."
	case "rejected":
		subject = "Your Event Update Has Been Rejected"
		title = "Update Rejected"
		color = "#EF4444"
		message = "Your recent event update could not be approved. Please review the feedback below and make the necessary changes."
	default:
		subject = "Event Update Status"
		title = "Update Status Changed"
		color = "#6B7280"
		message = "Your event update status has changed."
	}

	reasonBlock := ""
	if status == "rejected" && reason != "" {
		reasonBlock = emailInfoBox("#FEF2F2", "#FECACA", fmt.Sprintf(`
<p style="margin:0 0 6px 0;font-size:12px;font-weight:600;color:#991B1B;text-transform:uppercase;letter-spacing:0.05em;">Reason for Rejection</p>
<p style="margin:0;font-size:14px;color:#7F1D1D;line-height:1.6;">%s</p>`, html.EscapeString(reason)))
	}

	eventBox := emailInfoBox("#F8FAFC", "#E2E8F0", emailRow("Event Name", eventName))

	body := fmt.Sprintf(`
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:16px;">
  Hello <strong style="color:#1E293B;">%s</strong>,
</p>
<p style="font-size:15px;line-height:1.7;color:#475569;margin-bottom:0;">
  %s
</p>

%s
%s

%s
%s
<p style="font-size:13px;color:#94A3B8;line-height:1.6;">
  If you have any questions, please reach out to our support team.
</p>`,
		organizerName, message, eventBox, reasonBlock,
		emailButton(color, "https://ngevent.id/organizer/events", "View My Events"),
		emailDivider(),
	)

	htmlBody := emailWrapper(color, title, body)
	return sendMail(email, subject+" — Ngevent", htmlBody)
}
