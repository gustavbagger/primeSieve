package main

import (
	"errors"
	"testing"
)

func checkErr(got, want error) bool {
	if want == nil {
		return got == nil
	} else if got == nil {
		return false
	} else {
		return got.Error() == want.Error()
	}
}

func TestCmp192(t *testing.T) {
	tests := []struct {
		name string
		a, b uint192
		want int
	}{
		{
			name: "equal values",
			a:    uint192{Hi: 1, Mid: 2, Lo: 3},
			b:    uint192{Hi: 1, Mid: 2, Lo: 3},
			want: 0,
		},
		{
			name: "a < b by Hi",
			a:    uint192{Hi: 1, Mid: 0, Lo: 0},
			b:    uint192{Hi: 2, Mid: 0, Lo: 0},
			want: -1,
		},
		{
			name: "a > b by Hi",
			a:    uint192{Hi: 3, Mid: 0, Lo: 0},
			b:    uint192{Hi: 2, Mid: 0, Lo: 0},
			want: 1,
		},
		{
			name: "a < b by Mid",
			a:    uint192{Hi: 5, Mid: 1, Lo: 0},
			b:    uint192{Hi: 5, Mid: 2, Lo: 0},
			want: -1,
		},
		{
			name: "a > b by Mid",
			a:    uint192{Hi: 5, Mid: 3, Lo: 0},
			b:    uint192{Hi: 5, Mid: 2, Lo: 0},
			want: 1,
		},
		{
			name: "a < b by Lo",
			a:    uint192{Hi: 9, Mid: 9, Lo: 1},
			b:    uint192{Hi: 9, Mid: 9, Lo: 2},
			want: -1,
		},
		{
			name: "a > b by Lo",
			a:    uint192{Hi: 9, Mid: 9, Lo: 3},
			b:    uint192{Hi: 9, Mid: 9, Lo: 2},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmp192(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("cmp192(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSub192(t *testing.T) {
	tests := []struct {
		name    string
		a, b    uint192
		want    uint192
		wantErr error
	}{
		{
			name:    "equal values",
			a:       uint192{Hi: 1, Mid: 2, Lo: 3},
			b:       uint192{Hi: 1, Mid: 2, Lo: 3},
			want:    uint192{Hi: 0, Mid: 0, Lo: 0},
			wantErr: nil,
		},
		{
			name:    "a < b by Hi",
			a:       uint192{Hi: 1, Mid: 0, Lo: 0},
			b:       uint192{Hi: 2, Mid: 0, Lo: 0},
			want:    uint192{Hi: 0, Mid: 0, Lo: 0},
			wantErr: errors.New("a<b"),
		},
		{
			name: "a > b by Hi",
			a:    uint192{Hi: 3, Mid: 0, Lo: 0},
			b:    uint192{Hi: 2, Mid: 0, Lo: 0},
			want: uint192{Hi: 1, Mid: 0, Lo: 0},
		},
		{
			name:    "a < b by Mid",
			a:       uint192{Hi: 6, Mid: 1, Lo: 0},
			b:       uint192{Hi: 5, Mid: 2, Lo: 0},
			want:    uint192{Hi: 0, Mid: (1 << 64) - 1, Lo: 0},
			wantErr: nil,
		},
		{
			name: "a > b by Mid",
			a:    uint192{Hi: 5, Mid: 3, Lo: 0},
			b:    uint192{Hi: 5, Mid: 2, Lo: 0},
			want: uint192{Hi: 0, Mid: 1, Lo: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := sub192(tt.a, tt.b)
			if got != tt.want || !checkErr(gotErr, tt.wantErr) {
				t.Fatalf("sub192(%v, %v) = %d,%v, want %d with error %v", tt.a, tt.b, got, gotErr, tt.want, tt.wantErr)
			}
		})
	}
}

func TestAdd192(t *testing.T) {
	tests := []struct {
		name string
		a, b uint192
		want uint192
	}{
		{

			name: "zero plus zero",
			a:    uint192{Lo: 0, Mid: 0, Hi: 0},
			b:    uint192{Lo: 0, Mid: 0, Hi: 0},
			want: uint192{Lo: 0, Mid: 0, Hi: 0},
		},
		{
			name: "simple add no carry",
			a:    uint192{Lo: 1, Mid: 0, Hi: 0},
			b:    uint192{Lo: 2, Mid: 0, Hi: 0},
			want: uint192{Lo: 3, Mid: 0, Hi: 0},
		},
		{
			name: "carry from Lo to Mid",
			a:    uint192{Lo: ^uint64(0), Mid: 1, Hi: 0},
			b:    uint192{Lo: 1, Mid: 1, Hi: 0},
			want: uint192{Lo: 0, Mid: 3, Hi: 0},
		},
		{
			name: "carry from Mid to Hi",
			a:    uint192{Lo: 0, Mid: ^uint64(0), Hi: 1},
			b:    uint192{Lo: 0, Mid: 1, Hi: 1},
			want: uint192{Lo: 0, Mid: 0, Hi: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := add192(tt.a, tt.b)
			if cmp192(got, tt.want) != 0 {
				t.Fatalf("add192(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Check all want cases!
func TestMul192(t *testing.T) {
	tests := []struct {
		name string
		a, b uint192
		want [7]uint64
	}{
		{
			name: "max low limbs",
			a:    uint192{Lo: ^uint64(0), Mid: 0, Hi: 0},
			b:    uint192{Lo: ^uint64(0), Mid: 0, Hi: 0},
			want: [7]uint64{
				1,              // low limb
				^uint64(0) - 1, // 0xFFFFFFFFFFFFFFFE
				0,              // carry into limb 2
				0, 0, 0, 0,
			},
		},
		{
			name: "max mid limbs",
			a:    uint192{Lo: 0, Mid: ^uint64(0), Hi: 0},
			b:    uint192{Lo: 0, Mid: ^uint64(0), Hi: 0},
			want: [7]uint64{
				0,
				0,
				1,
				^uint64(0) - 1,
				0,
				0,
				0,
			},
		},
		{
			name: "max hi limbs",
			a:    uint192{Lo: 0, Mid: 0, Hi: ^uint64(0)},
			b:    uint192{Lo: 0, Mid: 0, Hi: ^uint64(0)},
			want: [7]uint64{
				0, 0, 0, 0,
				1,
				^uint64(0) - 1,
				0,
			},
		},

		{
			name: "zero times zero",
			a:    uint192{Lo: 0, Mid: 0, Hi: 0},
			b:    uint192{Lo: 0, Mid: 0, Hi: 0},
			want: [7]uint64{0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "low limbs only",
			a:    uint192{Lo: 3, Mid: 0, Hi: 0},
			b:    uint192{Lo: 7, Mid: 0, Hi: 0},
			want: [7]uint64{21, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "mid limbs only",
			a:    uint192{Lo: 0, Mid: 1, Hi: 0},  // 2^64
			b:    uint192{Lo: 0, Mid: 1, Hi: 0},  // 2^64
			want: [7]uint64{0, 0, 1, 0, 0, 0, 0}, // 2^128
		},
		{
			name: "hi limbs only",
			a:    uint192{Lo: 0, Mid: 0, Hi: 1},  // 2^128
			b:    uint192{Lo: 0, Mid: 0, Hi: 1},  // 2^128
			want: [7]uint64{0, 0, 0, 0, 1, 0, 0}, // 2^256
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mul192(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("mul192(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestNeg192(t *testing.T) {
	tests := []struct {
		name string
		x    uint192
	}{
		{
			name: "negate simple value",
			x:    uint192{Lo: 1, Mid: 0, Hi: 0},
		},
		{
			name: "negate with mid and hi",
			x:    uint192{Lo: 9, Mid: 7, Hi: 5},
		},
		{
			name: "negate max value",
			x:    uint192{Lo: ^uint64(0), Mid: ^uint64(0), Hi: ^uint64(0)},
		},
	}

	zero := uint192{Lo: 0, Mid: 0, Hi: 0}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := neg192(tt.x)
			sum := add192(tt.x, n)

			if cmp192(sum, zero) != 0 {
				t.Fatalf("x + neg192(x) = %v, want %v (x=%v, neg=%v)",
					sum, zero, tt.x, n)
			}
		})
	}
}

func TestMulMod192(t *testing.T) {
	tests := []struct {
		name string
		a, b uint192
		want uint192
	}{
		{
			name: "zero",
			a:    uint192{Lo: 0, Mid: 0, Hi: 0},
			b:    uint192{Lo: 0, Mid: 0, Hi: 0},
			want: uint192{Lo: 0, Mid: 0, Hi: 0},
		},
		{
			name: "small product fits in 192 bits",
			a:    uint192{Lo: 3, Mid: 0, Hi: 0},
			b:    uint192{Lo: 7, Mid: 0, Hi: 0},
			want: uint192{Lo: 21, Mid: 0, Hi: 0},
		},
		{
			name: "mixed limbs",
			a:    uint192{Lo: 3, Mid: 2, Hi: 1},
			b:    uint192{Lo: 6, Mid: 5, Hi: 4},
			// want computed via mul192 oracle
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// compute expected via mul192 oracle
			want := tt.want
			if tt.name == "mixed limbs" {
				full := mul192(tt.a, tt.b)
				want = uint192{Lo: full[0], Mid: full[1], Hi: full[2]}
			}

			got := mulMod192(tt.a, tt.b)

			if cmp192(got, want) != 0 {
				t.Fatalf("mulMod192(%v, %v) = %v, want %v", tt.a, tt.b, got, want)
			}
		})
	}
}
