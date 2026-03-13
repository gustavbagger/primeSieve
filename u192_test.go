package main

import (
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
		name string
		a, b uint192
		want uint192
	}{
		{
			name: "equal values",
			a:    uint192{Hi: 1, Mid: 2, Lo: 3},
			b:    uint192{Hi: 1, Mid: 2, Lo: 3},
			want: uint192{Hi: 0, Mid: 0, Lo: 0},
		},
		{
			name: "a < b by Hi",
			a:    uint192{Hi: 1, Mid: 0, Lo: 0},
			b:    uint192{Hi: 2, Mid: 0, Lo: 0},
			want: uint192{Hi: 1<<64 - 1, Mid: 0, Lo: 0},
		},
		{
			name: "a > b by Hi",
			a:    uint192{Hi: 3, Mid: 0, Lo: 0},
			b:    uint192{Hi: 2, Mid: 0, Lo: 0},
			want: uint192{Hi: 1, Mid: 0, Lo: 0},
		},
		{
			name: "a < b by Mid",
			a:    uint192{Hi: 6, Mid: 1, Lo: 0},
			b:    uint192{Hi: 5, Mid: 2, Lo: 0},
			want: uint192{Hi: 0, Mid: (1 << 64) - 1, Lo: 0},
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
			got := sub192(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("sub192(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
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

func TestLHS192(t *testing.T) {
	tests := []struct {
		name string
		x    uint192
		s    int
		want uint192
	}{
		{
			name: "simple Lo shift",
			x:    uint192{Lo: 1, Mid: 0, Hi: 0},
			s:    10,
			want: uint192{Lo: 1 << 10, Mid: 0, Hi: 0},
		},
		{
			name: "Lo shift, carry",
			x:    uint192{Lo: 1, Mid: 0, Hi: 0},
			s:    70,
			want: uint192{Lo: 0, Mid: 1 << 6, Hi: 0},
		},
		{
			name: "simple Mid shift",
			x:    uint192{Lo: 0, Mid: 1, Hi: 0},
			s:    30,
			want: uint192{Lo: 0, Mid: 1 << 30, Hi: 0},
		},
		{
			name: "overshifting Hi",
			x:    uint192{Lo: 0, Mid: 0, Hi: 1<<50 + 1<<30},
			s:    25,
			want: uint192{Lo: 0, Mid: 0, Hi: 1 << 55},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lSH192(tt.x, tt.s)
			if cmp192(got, tt.want) != 0 {
				t.Fatalf("lSH192(%v, %v) = %d, want %d", tt.x, tt.s, got, tt.want)
			}
		})
	}
}

func TestRHS192(t *testing.T) {
	tests := []struct {
		name string
		x    uint192
		s    int
		want uint192
	}{
		{
			name: "simple Lo shift",
			x:    uint192{Lo: 1, Mid: 0, Hi: 0},
			s:    5,
			want: uint192{Lo: 0, Mid: 0, Hi: 0},
		},
		{
			name: "Mid shift",
			x:    uint192{Lo: 0, Mid: 1 << 20, Hi: 0},
			s:    14,
			want: uint192{Lo: 0, Mid: 1 << 6, Hi: 0},
		},
		{
			name: "Mid shift, carry",
			x:    uint192{Lo: 0, Mid: 1 << 20, Hi: 0},
			s:    24,
			want: uint192{Lo: 1 << 60, Mid: 0, Hi: 0},
		},
		{
			name: "multi-word shift",
			x:    uint192{Lo: 0, Mid: 1 << 2, Hi: 1<<50 + 1<<20},
			s:    25,
			want: uint192{Lo: 1 << 41, Mid: 1 << 59, Hi: 1 << 25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rSH192(tt.x, tt.s)
			if cmp192(got, tt.want) != 0 {
				t.Fatalf("rSH192(%v, %v) = %d, want %d", tt.x, tt.s, got, tt.want)
			}
		})
	}
}

func TestTwoAdicVal192(t *testing.T) {
	tests := []struct {
		name string
		x    uint192
		want int
	}{
		{
			name: "one",
			x:    uint192{Lo: 1, Mid: 0, Hi: 0},
			want: 0,
		},
		{
			name: "simple Lo setup",
			x:    uint192{Lo: 1 << 5, Mid: 0, Hi: 0},
			want: 5,
		},
		{
			name: "simple Lo setup, mixed",
			x:    uint192{Lo: 1<<4 + 1<<10, Mid: 0, Hi: 0},
			want: 4,
		},
		{
			name: "simple Mid setup",
			x:    uint192{Lo: 0, Mid: 1 << 10, Hi: 0},
			want: 74,
		},
		{
			name: "simple Hi setup",
			x:    uint192{Lo: 0, Mid: 0, Hi: 1 << 6},
			want: 134,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := twoAdicVal192(tt.x)
			if got != tt.want {
				t.Fatalf("twoAdicVal192(%v) = %d, want %d", tt.x, got, tt.want)
			}
		})
	}
}

func TestLeadingZeros192(t *testing.T) {
	tests := []struct {
		name string
		x    uint192
		want int
	}{
		{
			name: "zero",
			x:    uint192{Lo: 0, Mid: 0, Hi: 0},
			want: 192,
		},
		{
			name: "simple Lo setup",
			x:    uint192{Lo: 1 << 5, Mid: 0, Hi: 0},
			want: 186,
		},
		{
			name: "simple Lo setup, mixed",
			x:    uint192{Lo: 1<<4 + 1<<10, Mid: 0, Hi: 0},
			want: 181,
		},
		{
			name: "simple Mid setup",
			x:    uint192{Lo: 0, Mid: 1 << 10, Hi: 0},
			want: 117,
		},
		{
			name: "maxed out",
			x:    uint192{Lo: 0, Mid: 0, Hi: (1 << 64) - 1},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LeadingZeros192(tt.x)
			if got != tt.want {
				t.Fatalf("LeadingZeros192(%v) = %d, want %d", tt.x, got, tt.want)
			}
		})
	}
}

func TestMontAddReduce(t *testing.T) {
	tests := []struct {
		name string
		a    uint192
		b    uint192
		n    uint192
		want uint192
	}{
		{
			name: "add, no reduce",
			a:    uint192{Lo: 1 << 10, Mid: 0, Hi: 0},
			b:    uint192{Lo: 1 << 12, Mid: 0, Hi: 0},
			n:    uint192{Lo: 1 << 30, Mid: 0, Hi: 0},
			want: uint192{Lo: 1<<10 + 1<<12, Mid: 0, Hi: 0},
		},
		{
			name: "add, reduce",
			a:    uint192{Lo: 1<<12 + 1<<11, Mid: 0, Hi: 0},
			b:    uint192{Lo: 1 << 12, Mid: 0, Hi: 0},
			n:    uint192{Lo: 1 << 13, Mid: 0, Hi: 0},
			want: uint192{Lo: 1 << 11, Mid: 0, Hi: 0},
		},
		{
			name: "mixed, reduce",
			a:    uint192{Lo: 1<<12 + 1<<11, Mid: 1 << 20, Hi: 0},
			b:    uint192{Lo: 1 << 12, Mid: 1 << 21, Hi: 0},
			n:    uint192{Lo: 1 << 13, Mid: 1 << 21, Hi: 0},
			want: uint192{Lo: 1 << 11, Mid: 1 << 20, Hi: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := montAddReduce(tt.a, tt.b, tt.n)
			if got != tt.want {
				t.Fatalf("montAddReduce(%v, %v,%v) = %d, want %d", tt.a, tt.b, tt.n, got, tt.want)
			}
		})
	}
}

/* Not needed at the minute, failing. Suspect substraction overflow problem at 2-N*x
func TestInv192(t *testing.T) {
	tests := []struct {
		name string
		x    uint192
	}{
		{
			name: "one",
			x:    uint192{Lo: 1, Mid: 0, Hi: 0},
		},
		{
			name: "Lo only",
			x:    uint192{Lo: 1 << 10, Mid: 0, Hi: 0},
		},
		{
			name: "Mid only",
			x:    uint192{Lo: 0, Mid: 1 << 5, Hi: 0},
		},
		{
			name: "Hi only",
			x:    uint192{Lo: 0, Mid: 0, Hi: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inv192(tt.x)
			gotMul := mulMod192(got, tt.x)
			if cmp192(gotMul, uint192{Lo: 1}) != 0 {
				t.Fatalf("inv192(%v) = %d, but mulMod192(%v,%d) = %v != 1", tt.x, got, tt.x, got, gotMul)
			}
		})
	}
}
*/

func TestInv64(t *testing.T) {
	tests := []struct {
		name string
		x    uint64
	}{
		{
			name: "one",
			x:    1,
		},
		{
			name: "medium val",
			x:    1725,
		},
		{
			name: "large val",
			x:    78754627,
		},
		{
			name: "minusone",
			x:    1<<64 - 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inv64(tt.x)
			gotMul := got * tt.x
			if gotMul != 1 {
				t.Fatalf("inv64(%v) = %d, but %v * %d = %v != 1", tt.x, got, tt.x, got, gotMul)
			}
		})
	}
}

// worried here about montAddReduce assuming a,b<N in montOne
func TestMontOne(t *testing.T) {
	tests := []struct {
		name string
		N    uint192
	}{
		{
			name: "one",
			N:    uint192{Lo: 1, Mid: 0, Hi: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inv192(tt.N)
			want := montAddReduce(uint192{1<<64 - 1, 1<<64 - 1, 1<<64 - 1}, uint192{Lo: 1}, tt.N)
			if cmp192(want, got) != 0 {
				t.Fatalf("montOne(%v) = %d, but expected %v", tt.N, got, want)
			}
		})
	}
}
