package main

import "testing"

func TestVideoIDFromLink(t *testing.T) {
	cases := map[string]string{
		"https://music.youtube.com/watch?v=dN3y8kVNplA&si=abc": "dN3y8kVNplA",
		"https://www.youtube.com/watch?v=dN3y8kVNplA":          "dN3y8kVNplA",
		"https://m.youtube.com/watch?list=x&v=dN3y8kVNplA":     "dN3y8kVNplA",
		"https://youtu.be/dN3y8kVNplA?t=10":                    "dN3y8kVNplA",
		"https://youtube.com/shorts/dN3y8kVNplA":               "dN3y8kVNplA",
		"  https://music.youtube.com/watch?v=dN3y8kVNplA  ":    "dN3y8kVNplA",
		"rosalia despecha":                                                   "",
		"https://example.com/watch?v=dN3y8kVNplA":                            "",
		"https://music.youtube.com/playlist?list=PL123":                      "",
		"https://music.youtube.com/watch?v=short":                            "",
		"javascript:alert(1)//https://music.youtube.com/watch?v=aaaaaaaaaaa": "",
	}
	for in, want := range cases {
		got, ok := videoIDFromLink(in)
		if ok != (want != "") || ok && got != want {
			t.Errorf("%q: got %q, %v; want %q", in, got, ok, want)
		}
	}
}
