package controllers

import (
	"errors"
	"strings"
)

func splitIDs(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func errNoID() error { return errors.New("no id provided") }

func errForbidden() error { return errors.New("forbidden") }
