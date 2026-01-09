package handler

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/useragent"
)

func Handler(c *fiber.Ctx) string {
	uaHeader := c.Get("User-Agent")
	referer := c.Get("Referer")

	// Parse user agent
	ua := useragent.New(uaHeader)
	name, version := ua.Browser()
	os := ua.OS()
	platform := ua.Platform()
	engine := DetectEngine(referer)

	userAgent := fmt.Sprintf(`%s %s (%s %s) %s`, name, version, engine, os, platform)

	return userAgent
}

func DetectEngine(ref string) string {
	if ref == "" {
		return ""
	}

	url, err := url.Parse(ref)
	if err != nil {
		return ""
	}

	host := url.Hostname()
	var engine string

	switch {
	case strings.Contains(host, "google."):
		engine = "google"
	case strings.Contains(host, "mozilla."):
		engine = "mozilla"
	case strings.Contains(host, "bing."):
		engine = "bing"
	default:
		engine = "host"
	}

	return engine
}
