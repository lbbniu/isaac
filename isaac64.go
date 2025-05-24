package isaac

import "math"

type UINT64_C = uint64

// Isaac64 corresponds to struct isaac_state
type Isaac64 struct {
	m [ISAAC_WORDS]uint64 // state table
	r []uint64            // result table
	a uint64
	b uint64
	c uint64
}

func just64(a uint64) uint64 {
	// return a & ((1 << 1 << (ISAAC_BITS - 1)) - 1)
	return a & math.MaxUint64
}

// ind64 corresponds to the C macro: ind64(mm, x) = *(ub8*)((ub1*)(mm) + ((x) & ((RANDSIZ-1)<<3)))
// Explanation: Perform byte-level offset on mm, then take 64-bit integer.
// Equivalent in Go: mm[( (x) & ((RANDSIZ-1)<<3)) >> 3].
func ind64(m [ISAAC_WORDS]uint64, x uint64) uint64 {
	return m[(x&((ISAAC_WORDS-1)*8))>>3]
}

// mix64 corresponds to the C macro mix64(a,b,c,d,e,f,g,h)
func mix64(a, b, c, d, e, f, g, h uint64) (na, nb, nc, nd, ne, nf, ng, nh uint64) {
	a -= e
	f ^= (just64(h) >> 9)
	h += a
	b -= f
	g ^= (a << 9)
	a += b
	c -= g
	h ^= (just64(b) >> 23)
	b += c
	d -= h
	a ^= (c << 15)
	c += d
	e -= a
	b ^= (just64(d) >> 14)
	d += e
	f -= b
	c ^= (e << 20)
	e += f
	g -= c
	d ^= (just64(f) >> 17)
	f += g
	h -= d
	e ^= (g << 14)
	g += h
	return a, b, c, d, e, f, g, h
}

// isaac_refill corresponds to the C version of isaac_refill function
func (s *Isaac64) isaac_refill(r *[ISAAC_WORDS]uint64) {
	a := s.a
	b := s.b + (s.c + 1)
	s.c++

	HALF := ISAAC_WORDS / 2

	// isaac_step corresponds to the C ISAAC_STEP macro
	step := func(i int, off int, mix uint64) {
		a = (0 ^ mix) + s.m[off+i]
		x := s.m[i]
		y := ind64(s.m, x) + a + b
		s.m[i] = y
		b = just64(ind64(s.m, y>>ISAAC_WORDS_LOG) + x)
		r[i] = b
	}

	// First half
	for i := 0; i < HALF; i += 4 {
		// step1: a = ^ (a ^ (a << 21))
		step(i, HALF, ^(a ^ (a << 21)))
		// step2: a = a ^ (a >> 5)
		step(i+1, HALF, a^(just64(a)>>5))
		// step3: a = a ^ (a << 12)
		step(i+2, HALF, a^(a<<12))
		// step4: a = a ^ (a >> 33)
		step(i+3, HALF, a^(just64(a)>>33))
	}

	// Second half
	for i := HALF; i < ISAAC_WORDS; i += 4 {
		// step1: a = ^ (a ^ (a << 21))
		step(i, -HALF, ^(a ^ (a << 21)))
		// step2: a = a ^ (just (a) >>  5)
		step(i+1, -HALF, a^(just64(a)>>5))
		// step3: a = a^(a<<12)
		step(i+2, -HALF, a^(a<<12))
		// step4: a = a^(just64(a)>>33)
		step(i+3, -HALF, a^(just64(a)>>33))
	}

	s.a = a
	s.b = b
}

// NewIsaac64 creates a new ISAAC64 instance
func NewIsaac64() *Isaac64 {
	return &Isaac64{}
}

// Seed initializes ISAAC64
// Corresponds to the C isaac_seed function
func (s *Isaac64) Seed(seed [ISAAC_WORDS]uint64, initValues ...uint64) {
	if len(initValues) > 0 && len(initValues) != 8 {
		panic("isaac: need exactly 8 initial values for uint64")
	}

	// Use the same initial values as the C version
	var a, b, c, d, e, f, g, h uint64
	if len(initValues) == 8 {
		a = initValues[0]
		b = initValues[1]
		c = initValues[2]
		d = initValues[3]
		e = initValues[4]
		f = initValues[5]
		g = initValues[6]
		h = initValues[7]
	} else {
		a = 0x647c4677a2884b7c
		b = 0xb9f8b322c73ac862
		c = 0x8c0ea5053d4712a0
		d = 0xb29b2e824a595524
		e = 0x82f053db8355e0ce
		f = 0x48fe4a0fa5a09315
		g = 0xae985bf2cbfc89ed
		h = 0x98f5704f6c44c0ab
	}

	// Initialize m array
	for i := 0; i < ISAAC_WORDS; i++ {
		s.m[i] = seed[i]
	}

	// Mix S->m so that every part of the seed affects every part of the state
	// Two rounds of mixing
	for range [2]struct{}{} {
		for i := 0; i < ISAAC_WORDS; i += 8 {
			a += s.m[i]
			b += s.m[i+1]
			c += s.m[i+2]
			d += s.m[i+3]
			e += s.m[i+4]
			f += s.m[i+5]
			g += s.m[i+6]
			h += s.m[i+7]
			a, b, c, d, e, f, g, h = mix64(a, b, c, d, e, f, g, h)
			s.m[i] = a
			s.m[i+1] = b
			s.m[i+2] = c
			s.m[i+3] = d
			s.m[i+4] = e
			s.m[i+5] = f
			s.m[i+6] = g
			s.m[i+7] = h
		}
	}

	s.a = 0
	s.b = 0
	s.c = 0
}

func (s *Isaac64) Refill(r *[ISAAC_WORDS]uint64) {
	s.isaac_refill(r)
}

func (s *Isaac64) Uint64() uint64 {
	if r := s.r; len(r) == 0 {
		var r [ISAAC_WORDS]uint64
		s.Refill(&r)
		s.r = r[:]
	}
	r := s.r[0]
	s.r = s.r[1:]
	return r
}
