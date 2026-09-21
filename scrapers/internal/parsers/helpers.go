package parsers

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// closedListingKeywords are phrases that mean "don't bother the user with this one".
// If any of these appear literally on the card page, the parser emits nothing for it.
var closedListingKeywords = []string{
	"no longer accepting applications",
	"applications closed",
	"application closed",
	"this job is no longer accepting applications",
	"job closed",
	"vacancy closed",
	"no longer accepting",
	"closed to applicants",
	"we are no longer hiring for this role",
	"hiring paused",
	"not accepting applications",
	"role filled",
	"this role has been filled",
	"position closed",
}

func getIDFromURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Sprint(time.Now().UnixNano())
	}

	url = strings.TrimSuffix(url, "/")
	url = strings.Split(url, "?")[0]
	url = strings.Split(url, "#")[0]
	parts := strings.FieldsFunc(url, func(r rune) bool {
		return r == '/' || r == '-' || r == '_' || r == '='
	})
	if len(parts) == 0 {
		return fmt.Sprint(time.Now().UnixNano())
	}
	return parts[len(parts)-1]
}

// IsClosedListing returns true when the visible page/card text contains any
// "closed / no longer accepting" phrase. This lets each parser drop dead jobs
// before they ever enter the DB, so the UI never shows stale "closed" cards.
func IsClosedListing(pageText string) bool {
	lower := strings.ToLower(pageText)
	for _, kw := range closedListingKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

var relativeNumRe = regexp.MustCompile(`(\d+)\s*(hour|hours|hr|hrs|day|days|d|week|weeks|wk|wks|month|months|min|mins|minute|minutes|second|seconds)`)

// ParseRelativeDate converts strings like "about 19 hours ago",
// "3 days ago", "Active 28 days ago", "just posted", "today" into
// an absolute time.Time value. Falls back to time.Now() on empty input.
func ParseRelativeDate(text string) time.Time {
	text = strings.ToLower(strings.TrimSpace(text))
	now := time.Now()

	if text == "" {
		return now
	}

	if strings.Contains(text, "today") ||
		strings.Contains(text, "just posted") ||
		strings.Contains(text, "hiring ongoing") ||
		strings.Contains(text, "posted today") {
		return now
	}

	match := relativeNumRe.FindStringSubmatch(text)
	if match == nil {
		// fallback: if "ago" is present without a number, guess 1 day
		if strings.Contains(text, "ago") {
			return now.AddDate(0, 0, -1)
		}
		return now
	}

	n := 0
	if _, err := fmt.Sscanf(match[1], "%d", &n); err != nil || n <= 0 {
		return now
	}
	unit := match[2]

	switch unit {
	case "hour", "hours", "hr", "hrs":
		return now.Add(-time.Duration(n) * time.Hour)
	case "day", "days", "d":
		return now.AddDate(0, 0, -n)
	case "week", "weeks", "wk", "wks":
		return now.AddDate(0, 0, -n*7)
	case "month", "months":
		if n > 12 {
			n = 12
		}
		return now.AddDate(0, -n, 0)
	case "min", "mins", "minute", "minutes":
		return now.Add(-time.Duration(n) * time.Minute)
	case "second", "seconds":
		return now.Add(-time.Duration(n) * time.Second)
	}

	return now
}
