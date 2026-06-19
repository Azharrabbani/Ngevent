package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"ngevent/internal/dto"
	"ngevent/internal/repository"
	"os"
	"os/exec"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ReportService struct {
	ReportRepo    repository.ReportRepo
	OrganizerRepo repository.OrganizerProfileRepo
}

func NewReportService(
	reportRepo repository.ReportRepo,
	organizerRepo repository.OrganizerProfileRepo,
) *ReportService {
	return &ReportService{ReportRepo: reportRepo, OrganizerRepo: organizerRepo}
}

func (s *ReportService) GenerateEventReportPDF(
	userID string,
	filter *dto.EventReportFilter,
) ([]byte, string, error) {

	profile, err := s.OrganizerRepo.FindByUserID(userID)
	if err != nil {
		log.Printf("[ReportService] FindByUserID error: %v", err)
		return nil, "", err
	}

	summary, err := s.ReportRepo.GetEventReport(profile.ID, filter)
	if err != nil {
		log.Printf("[ReportService] GetEventReport error: %v", err)
		return nil, "", err
	}

	htmlStr, err := renderHTML(summary)
	if err != nil {
		log.Printf("[ReportService] renderHTML error: %v", err)
		return nil, "", err
	}

	pdf, err := generatePDF(htmlStr)
	if err != nil {
		log.Printf("[ReportService] generatePDF error: %v", err)
		return nil, "", err
	}

	filename := fmt.Sprintf("event-report-%s.pdf", summary.Period)
	return pdf, filename, nil
}

const reportTemplateHTML = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

  * { box-sizing: border-box; margin: 0; padding: 0; }

  body {
    font-family: 'Inter', Arial, sans-serif;
    font-size: 12px;
    color: #0f172a;
    background: #fff;
    padding: 32px 40px 48px;
  }

  /* ── Header ───────────────────────────────── */
  .header {
    background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
    border-radius: 10px;
    padding: 28px 32px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 6px;
  }
  .header-left h1 {
    font-size: 24px;
    font-weight: 700;
    color: #fff;
    letter-spacing: 1.5px;
  }
  .header-left p {
    font-size: 11px;
    color: #94a3b8;
    margin-top: 3px;
    font-weight: 400;
  }
  .header-right {
    text-align: right;
  }
  .header-right .period-label {
    font-size: 10px;
    color: #64748b;
    text-transform: uppercase;
    letter-spacing: 0.8px;
    font-weight: 600;
  }
  .header-right .period-value {
    font-size: 15px;
    font-weight: 700;
    color: #f8fafc;
    margin-top: 2px;
  }

  /* ── Divider ───────────────────────────────── */
  .divider {
    height: 3px;
    background: linear-gradient(90deg, #3b82f6, #6366f1, #8b5cf6);
    border-radius: 0 0 4px 4px;
    margin-bottom: 28px;
  }

  /* ── Summary cards ─────────────────────────── */
  .section-title {
    font-size: 10px;
    font-weight: 700;
    color: #94a3b8;
    text-transform: uppercase;
    letter-spacing: 1px;
    margin-bottom: 10px;
  }

  .summary {
    display: flex;
    gap: 8px;
    margin-bottom: 28px;
  }
  .stat-card {
    flex: 1;
    padding: 16px 12px 14px;
    border-radius: 8px;
    text-align: center;
    border: 1px solid #e2e8f0;
    background: #fff;
    position: relative;
    overflow: hidden;
  }
  .stat-card::before {
    content: '';
    position: absolute;
    top: 0; left: 0; right: 0;
    height: 3px;
  }
  .stat-card.total::before   { background: #1e293b; }
  .stat-card.active::before  { background: #16a34a; }
  .stat-card.pending::before { background: #ea580c; }
  .stat-card.done::before    { background: #3b82f6; }
  .stat-card.rejected::before{ background: #dc2626; }

  .stat-card .num {
    font-size: 28px;
    font-weight: 700;
    line-height: 1;
    margin-top: 4px;
  }
  .stat-card .lbl {
    font-size: 10px;
    font-weight: 500;
    color: #64748b;
    margin-top: 5px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .num-total    { color: #1e293b; }
  .num-active   { color: #16a34a; }
  .num-pending  { color: #ea580c; }
  .num-done     { color: #3b82f6; }
  .num-rejected { color: #dc2626; }

  /* ── Table ─────────────────────────────────── */
  .table-wrap {
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid #e2e8f0;
  }
  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 11px;
  }
  thead tr th {
    background: #1e293b;
    color: #e2e8f0;
    padding: 11px 14px;
    text-align: left;
    font-weight: 600;
    font-size: 10px;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    white-space: nowrap;
  }
  thead tr th:first-child { width: 28px; text-align: center; }

  tbody tr { transition: background 0.1s; }
  tbody tr:nth-child(odd)  { background: #f8fafc; }
  tbody tr:nth-child(even) { background: #ffffff; }

  tbody tr td {
    padding: 11px 14px;
    border-bottom: 1px solid #f1f5f9;
    vertical-align: middle;
    color: #334155;
  }
  tbody tr:last-child td { border-bottom: none; }
  td:first-child { text-align: center; color: #94a3b8; font-weight: 600; font-size: 10px; }

  .event-name { font-weight: 600; color: #0f172a; font-size: 11px; }
  .city       { color: #64748b; }
  .date       { white-space: nowrap; color: #475569; }
  .time-range { white-space: nowrap; color: #64748b; font-variant-numeric: tabular-nums; }

  /* ── Status badge ──────────────────────────── */
  .badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 3px 9px;
    border-radius: 99px;
    font-size: 9px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.4px;
    white-space: nowrap;
  }
  .badge::before {
    content: '';
    width: 5px; height: 5px;
    border-radius: 50%;
    display: inline-block;
  }
  .badge-active   { background: #dcfce7; color: #15803d; }
  .badge-active::before   { background: #16a34a; }
  .badge-pending  { background: #ffedd5; color: #c2410c; }
  .badge-pending::before  { background: #ea580c; }
  .badge-rejected { background: #fee2e2; color: #b91c1c; }
  .badge-rejected::before { background: #dc2626; }
  .badge-cancelled { background: #fee2e2; color: #b91c1c; }
  .badge-cancelled::before { background: #dc2626; }
  .badge-done     { background: #dbeafe; color: #1d4ed8; }
  .badge-done::before     { background: #3b82f6; }
  .badge-default  { background: #475569; color: #f1f5f9; }
  .badge-default::before  { background: #94a3b8; }

  /* ── Rejected reason row ───────────────────── */
  .reason-row td {
    background: #fff5f5;
    color: #b91c1c;
    font-size: 10px;
    padding: 7px 14px 9px 28px;
    border-bottom: 1px solid #fee2e2;
    font-style: italic;
  }

  /* ── Footer ────────────────────────────────── */
  .footer {
    margin-top: 28px;
    text-align: center;
    font-size: 10px;
    color: #94a3b8;
  }
</style>
</head>
<body>

  <!-- Header -->
  <div class="header">
    <div class="header-left">
      <h1>NGEVENT</h1>
      <p>Event Report</p>
    </div>
    <div class="header-right">
      <div class="period-label">Period</div>
      <div class="period-value">{{.Period}}</div>
    </div>
  </div>
  <div class="divider"></div>

  <!-- Summary -->
  <div class="section-title">Overview</div>
  <div class="summary">
    <div class="stat-card total">
      <div class="num num-total">{{.Total}}</div>
      <div class="lbl">Total</div>
    </div>
    <div class="stat-card active">
      <div class="num num-active">{{.Active}}</div>
      <div class="lbl">Active</div>
    </div>
    <div class="stat-card pending">
      <div class="num num-pending">{{.Pending}}</div>
      <div class="lbl">Pending</div>
    </div>
    <div class="stat-card done">
      <div class="num num-done">{{.Done}}</div>
      <div class="lbl">Done</div>
    </div>
    <div class="stat-card rejected">
      <div class="num num-rejected">{{.Rejected}}</div>
      <div class="lbl">Rejected</div>
    </div>
  </div>

  <!-- Event table -->
  <div class="section-title">Event List</div>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>#</th>
          <th>Event Name</th>
          <th>Status</th>
          <th>City</th>
          <th>Date</th>
          <th>Time (WIB)</th>
        </tr>
      </thead>
      <tbody>
        {{range $i, $r := .Rows}}
        <tr>
          <td>{{inc $i}}</td>
          <td class="event-name">{{$r.Name}}</td>
          <td><span class="badge badge-{{$r.Status}}">{{$r.Status}}</span></td>
          <td class="city">{{$r.City}}</td>
          <td class="date">{{$r.StartDate}} – {{$r.EndDate}}</td>
          <td class="time-range">{{$r.TimeRange}}</td>
        </tr>
        {{if $r.RejectedReason}}
        <tr class="reason-row">
          <td colspan="6">⚠ {{deref $r.RejectedReason}}</td>
        </tr>
        {{end}}
        {{end}}
      </tbody>
    </table>
  </div>

  <div class="footer">Generated by Ngevent &nbsp;·&nbsp; {{.Period}}</div>

</body>
</html>`

var templateFuncs = template.FuncMap{
	"inc": func(i int) int { return i + 1 },
	"deref": func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	},
}

func renderHTML(summary *dto.EventReportSummary) (string, error) {
	tmpl, err := template.New("report").Funcs(templateFuncs).Parse(reportTemplateHTML)
	if err != nil {
		return "", fmt.Errorf("template parse: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, summary); err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return buf.String(), nil
}

func generatePDF(htmlStr string) ([]byte, error) {
	chromePath, err := findChrome()
	if err != nil {
		return nil, fmt.Errorf("chrome not found: %w", err)
	}
	log.Printf("[generatePDF] using chrome: %s", chromePath)

	// Write HTML to a temp file and use file:// URI.
	tmpFile, err := os.CreateTemp("", "ngevent-report-*.html")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(htmlStr); err != nil {
		return nil, fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	fileURI := "file://" + tmpFile.Name()

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("headless", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	var pdfBytes []byte

	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURI),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBytes, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp run: %w", err)
	}

	return pdfBytes, nil
}

func findChrome() (string, error) {
	if bin := os.Getenv("CHROME_BIN"); bin != "" {
		return bin, nil
	}
	candidates := []string{
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
		"/usr/bin/google-chrome",
		"/usr/bin/google-chrome-stable",
		"/snap/bin/chromium",
		"chromium-browser",
		"chromium",
		"google-chrome",
		"google-chrome-stable",
	}
	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("none of %s found in PATH", strings.Join(candidates, ", "))
}
