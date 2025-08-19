package main

import (
	"testing"
	"time"

	"snippetbox.glebich/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 5, 25, 15, 30, 0, 0, time.UTC),
			want: "25 May 2025 at 15:30",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2025, 6, 25, 15, 30, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "25 Jun 2025 at 14:30",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}

/*
func TestHumanDate(t *testing.T) {
	tm := time.Date(2025, 5, 25, 15, 30, 0, 0, time.UTC)
	hd := humanDate(tm)
	if hd != "25 May 2025 at 15:30" {
		t.Errorf("got %q; want %q", hd, "25 May 2025 at 15:30")
	}
}
*/
