package main

import "testing"

func TestRegion(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"", defaultRegion, false},
		{"es", "ES", false},
		{" US ", "US", false},
		{"ESP", "", true},
		{"e1", "", true},
	}
	for _, c := range cases {
		got, err := configuration{Region: c.in}.region()
		if (err != nil) != c.wantErr || got != c.want {
			t.Errorf("%q: got %q, %v", c.in, got, err)
		}
	}
}

func TestPostMessage(t *testing.T) {
	tr := track{VideoID: "dN3y8kVNplA", Title: "Song [Live]", Artists: "A & B"}
	want := "🎵 [Song \\[Live\\] — A & B](https://music.youtube.com/watch?v=dN3y8kVNplA)"
	if got := postMessage("  ", tr); got != want {
		t.Errorf("got %q", got)
	}
	if got := postMessage("listen", tr); got != "listen\n\n"+want {
		t.Errorf("got %q", got)
	}
}
