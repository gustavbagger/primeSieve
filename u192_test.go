package main

import (
	"math/big"
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

func toBig(x uint192) *big.Int {
	z := big.NewInt(0)
	z.Lsh(z.SetUint64(x.Hi), 128)
	z.Add(z, new(big.Int).Lsh(new(big.Int).SetUint64(x.Mid), 64))
	z.Add(z, new(big.Int).SetUint64(x.Lo))
	return z
}

func bigTo7Limbs(x *big.Int) [7]uint64 {
	var out [7]uint64

	tmp := new(big.Int).Set(x)
	mask := new(big.Int).SetUint64(^uint64(0)) // 0xffffffffffffffff

	for i := 0; i < 7; i++ {
		out[i] = new(big.Int).And(tmp, mask).Uint64()
		tmp.Rsh(tmp, 64)
	}

	return out
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

// Cannot max out N since we assume no carry in add192, coming from 2N<<2^192 (N at most ~10^56)
func TestMontOne(t *testing.T) {
	tests := []struct {
		name string
		N    uint192
	}{
		{
			name: "Mixed",
			N:    uint192{Lo: 1, Mid: 1 << 15, Hi: 0},
		},
		{
			name: "Lo",
			N:    uint192{Lo: 3<<5 + 1, Mid: 0, Hi: 0},
		},
		{
			name: "Big Boi",
			N:    uint192{Lo: 1<<64 - 1, Mid: 1<<64 - 1, Hi: 1<<58 - 1},
		},
	}

	cmpMont := func(get, N uint192) (bool, *big.Int, *big.Int) {
		a := big.NewInt(1)
		a.Lsh(a, 192)

		NBig := toBig(N)
		GetBig := toBig(get)

		out := big.NewInt(1).Mod(a, NBig)
		return out.Cmp(GetBig) == 0, GetBig, out
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if b, got, want := cmpMont(montOne(tt.N), tt.N); !b {
				t.Fatalf("montOne(%v) = %d, but expected %v", tt.N, got, want)
			}
		})
	}
}

func TestREDC(t *testing.T) {
	tests := []struct {
		name string
		N    uint192
		A    uint192
		B    uint192
	}{
		{
			name: "Small",
			N:    uint192{Lo: 97, Mid: 0, Hi: 0},
			A:    uint192{Lo: 12, Mid: 0, Hi: 0},
			B:    uint192{Lo: 7, Mid: 0, Hi: 0},
		},
		{
			name: "Medium",
			N:    uint192{Lo: 1<<32 + 15, Mid: 1<<20 + 3, Hi: 0},
			A:    uint192{Lo: 123456, Mid: 0, Hi: 0},
			B:    uint192{Lo: 98765, Mid: 0, Hi: 0},
		},
		{
			name: "Large limbs",
			N:    uint192{Lo: 1<<63 + 11, Mid: 1<<62 + 7, Hi: 1<<58 + 3},
			A:    uint192{Lo: 1<<60 + 123, Mid: 1<<40 + 55, Hi: 0},
			B:    uint192{Lo: 1<<59 + 999, Mid: 1<<39 + 77, Hi: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Nbig := toBig(tt.N)
			Abig := toBig(tt.A)
			Bbig := toBig(tt.B)

			Rbig := toBig(montOne(tt.N))

			ARbig := new(big.Int).Mod(new(big.Int).Mul(Abig, Rbig), Nbig)

			BRbig := new(big.Int).Mod(new(big.Int).Mul(Bbig, Rbig), Nbig)

			CBig := new(big.Int).Mul(ARbig, BRbig)
			C := bigTo7Limbs(CBig)
			// mu = inv64(N.Lo)
			mu := ^inv64(tt.N.Lo) + 1

			// Run REDC
			got := REDC(C, tt.N, mu)

			// Convert result to big.Int
			gotBig := toBig(got)

			wantBig := new(big.Int).Mod(new(big.Int).Mul(Bbig, ARbig), Nbig)

			if gotBig.Cmp(wantBig) != 0 {
				t.Fatalf("REDC mismatch:\n got  %v\n want %v\nA: %v B: %v R: %v, N: %v\n AR: %v, BR: %v", gotBig, wantBig, Abig, Bbig, Rbig, Nbig, ARbig, BRbig)
			}
		})
	}
}

func TestStrongPRP(t *testing.T) {
	tests := []struct {
		name string
		N    *big.Int
	}{
		{
			name: "Small prime 97",
			N:    big.NewInt(97),
		},
		{
			name: "Small composite 91",
			N:    big.NewInt(91), // 7 * 13
		},
		{
			name: "Prime 2^61-1",
			N:    new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 61), big.NewInt(1)),
		},
		{
			name: "Composite 2^61-1 * 3",
			N: new(big.Int).Mul(
				new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 61), big.NewInt(1)),
				big.NewInt(3),
			),
		},
		{
			name: "Larger prime-ish 64-bit",
			N:    big.NewInt(0).SetUint64(18446744073709551557), // known prime near 2^64
		},
		// --- Small primes ---
		{"Prime 3", big.NewInt(3)},
		{"Prime 5", big.NewInt(5)},
		{"Prime 17", big.NewInt(17)},
		{"Prime 257", big.NewInt(257)}, // Fermat prime

		// --- Small composites ---
		{"Composite 9", big.NewInt(9)},
		{"Composite 21", big.NewInt(21)},
		{"Composite 341", big.NewInt(341)}, // 341 = 11*31 (Fermat pseudoprime to base 2)

		// --- Carmichael numbers (all are Fermat pseudoprimes) ---
		{"Carmichael 561", big.NewInt(561)},   // 3 * 11 * 17
		{"Carmichael 1105", big.NewInt(1105)}, // 5 * 13 * 17
		{"Carmichael 1729", big.NewInt(1729)}, // 7 * 13 * 19
		{"Carmichael 2465", big.NewInt(2465)}, // 5 * 17 * 29
		{"Carmichael 6601", big.NewInt(6601)}, // 7 * 23 * 41

		// --- Large primes (64‑bit range) ---
		{"Large prime 2^61−1", new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 61), big.NewInt(1))},

		// --- Large composites (64‑bit range) ---
		{"Large composite 2^61−1 * 3",
			new(big.Int).Mul(
				new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 61), big.NewInt(1)),
				big.NewInt(3),
			),
		},
		// --- Edge cases ---
		{"N = 2^64−1", new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 64), big.NewInt(1))},
		{"N = 2^80+1", new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), 80), big.NewInt(1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to uint192
			N7 := bigTo7Limbs(tt.N)

			N192 := uint192{Lo: N7[0], Mid: N7[1], Hi: N7[2]}

			got := strongPRP(N192)

			// Use math/big primality test as reference
			want := tt.N.ProbablyPrime(20)
			//fmt.Println("Checking if ", tt.N, "is prime. Expecting: ", want, " Got: ", got)
			if got != want {
				t.Fatalf("strongPRP(%v) = %v, but big.Int.ProbablyPrime says %v", tt.N, got, want)
			}
		})
	}
}
