package service

import "testing"

func TestFormatLink(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "Pando的Web应用程序是Pando Proto，可在https://app.pando.im上获得。这个Web应用程序的目标是为用户提供一个统一的界面，用于访问所有Pando协议和产品。目前，已经将4swap协议集成到Web应用程序中，并且将弃用4swap的旧Web应用程序（https://app.4swap.org）。将来，Leaf协议和Rings协议也将融合到新的Pando Web应用程序中。",
			want:  "Pando的Web应用程序是Pando Proto，可在 https://app.pando.im 上获得。这个Web应用程序的目标