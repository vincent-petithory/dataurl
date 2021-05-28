package dataurl

import (
	"net/url"
)

// Escape implements URL escaping, as defined in RFC 2397 (http://tools.ietf.org/html/rfc2397).
// It differs a bit from net/url's QueryEscape and QueryUnescape, e.g how spaces are treated (+ instead of %20):
//
// Only ASCII chars are allowed. Reserved chars are escaped to their %xx form.
// Unreserved chars are [a-z], [A-Z], [0-9], and -_.!~*\().
func Escape(data []byte) string {
	return url.PathEscape(string(data))
}

// EscapeString is like Escape, but taking
// a string as argument.
func EscapeString(s string) string {
	return url.PathEscape(s)
}

// Unescape unescapes a character sequence
// escaped with Escape(String?).
func Unescape(s string) ([]byte, error) {
	res, err := url.PathUnescape(s)
	return []byte(res), err
}

// UnescapeToString is like Unescape, but returning
// a string.
func UnescapeToString(s string) (string, error) {
	return url.PathUnescape(s)
}
