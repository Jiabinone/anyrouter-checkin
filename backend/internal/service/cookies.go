package service

import "strings"

func extractSessionValue(raw string) string {
	cookies := parseCookieHeader(raw)
	if value, ok := cookies["session"]; ok {
		return value
	}
	return strings.TrimSpace(raw)
}

func parseCookieHeader(raw string) map[string]string {
	cookies := make(map[string]string)
	if raw == "" {
		return cookies
	}
	if !strings.Contains(raw, "session=") && !strings.Contains(raw, ";") {
		return cookies
	}
	for _, part := range strings.Split(raw, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		name := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		if name == "" {
			continue
		}
		cookies[name] = value
	}
	return cookies
}
