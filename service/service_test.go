package service

import "testing"

func TestFormatLink(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{