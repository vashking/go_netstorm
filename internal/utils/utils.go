package utils

import (
	"strings"
	"net/url"
)

func ParseHeaders(headerStr string) map[string]string {
	headers := make(map[string]string)
	parts := strings.Split(headerStr, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return headers
}

func RandomizeHeaders(base map[string]string) map[string]string {
	headers := make(map[string]string)
	for k, v := range base {
		headers[k] = v
	}
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	headers["Referer"] = "https://" + randomString(8) + ".com/" + randomString(5)
	headers["Cookie"] = "sessionid=" + randomString(16)
	return headers
}

func RandomizeURL(base string) string {
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	q := u.Query()
	q.Set(randomString(5), randomString(8))
	u.RawQuery = q.Encode()
	return u.String()
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
