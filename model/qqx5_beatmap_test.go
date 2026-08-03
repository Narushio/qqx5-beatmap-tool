package model

import "testing"

func TestCalculateBarAndPosSupportsOddPositions(t *testing.T) {
	barMS, posMS := calculateBarAndPosMS(120)

	tests := []struct {
		name    string
		ms      float64
		wantBar int
		wantPos int
	}{
		{
			name:    "first odd position",
			ms:      posMS,
			wantBar: 1,
			wantPos: 1,
		},
		{
			name:    "last odd position in first bar",
			ms:      63 * posMS,
			wantBar: 1,
			wantPos: 63,
		},
		{
			name:    "position 65 carries into next bar",
			ms:      65 * posMS,
			wantBar: 2,
			wantPos: 1,
		},
		{
			name:    "position 129 carries across two bars",
			ms:      129 * posMS,
			wantBar: 3,
			wantPos: 1,
		},
		{
			name:    "rounded position 64 carries to position zero",
			ms:      63.6 * posMS,
			wantBar: 2,
			wantPos: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBar, gotPos := calculateBarAndPos(tt.ms, barMS, posMS)
			if gotBar != tt.wantBar || gotPos != tt.wantPos {
				t.Fatalf("calculateBarAndPos(%v) = (%d, %d), want (%d, %d)", tt.ms, gotBar, gotPos, tt.wantBar, tt.wantPos)
			}
		})
	}
}
