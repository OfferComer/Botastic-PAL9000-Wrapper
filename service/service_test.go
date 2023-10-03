package service

import "testing"

func TestFormatLink(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "Pando的Web应用程序是Pando Proto，可在ht