package cipher

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/skycoin/skycoin/src/cipher/ripemd160"
	"github.com/stretchr/testify/assert"
)

func freshSumRipemd160(b []byte) Ripemd160 {
	sh := ripemd160.New()
	sh.Write(b)
	h := Ripemd160{}
	h.Set(sh.Sum(nil))
	return h
}

func freshSumSHA256(b []byte) SHA256 {
	sh := sha256.New()
	sh.Write(b)
	h := SHA256{}
	h.Set(sh.Sum(nil))
	return h
}

func randBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	x, err := rand.Read(b)
	assert.Equal(t, n, x)
	assert.Nil(t, err)
	return b
}

func TestHashRipemd160(t *testing.T) {
	assert.NotPanics(t, func() { HashRipemd160(randBytes(t, 128)) })
	r := HashRipemd160(randBytes(t, 160))
	assert.NotEqual(t, r, Ripemd160{})
	// 2nd hash should not be affected by previous
	b := randBytes(t, 256)
	r2 := HashRipemd160(b)
	assert.NotEqual(t, r2, Ripemd160{})
	assert.Equal(t, r2, freshSumRipemd160(b))
}

func TestRipemd160Set(t *testing.T) {
	h := Ripemd160{}
	assert.Panics(t, func() {
		h.Set(randBytes(t, 21))
	})
	assert.Panics(t, func() {
		h.Set(randBytes(t, 100))
	})
	assert.Panics(t, func() {
		h.Set(randBytes(t, 19))
	})
	assert.Panics(t, func() {
		h.Set(randBytes(t, 0))
	})
	assert.NotPanics(t, func() {
		h.Set(randBytes(t, 20))
	})
	b := randBytes(t, 20)
	h.Set(b)
	assert.True(t, bytes.Equal(h[:], b))
}

func TestSHA256Set(t *testing.T) {
	h := SHA256{}
	assert.Panics(t, func() {
		h.Set(randBytes(t, 33))
	})
	assert.Panics(t, func() {
		h.Set(randBytes(t, 100))
	})
	assert.Panics(t, func() {
		h.Set(randBytes(t, 31))
	})
	assert.Panics(t, func() {
		h.Set(randBytes(t, 0))
	})
	assert.NotPanics(t, func() {
		h.Set(randBytes(t, 32))
	})
	b := randBytes(t, 32)
	h.Set(b)
	assert.True(t, bytes.Equal(h[:], b))
}

func TestSHA256Hex(t *testing.T) {
	h := SHA256{}
	h.Set(randBytes(t, 32))
	s := h.Hex()
	h2, err := SHA256FromHex(s)
	assert.Nil(t, err)
	assert.Equal(t, h, h2)
	assert.Equal(t, h2.Hex(), s)
}

func TestSHA256KnownValue(t *testing.T) {
	vals := []struct {
		input  string
		output string
	}{
		// These values are generated by
		// echo -n input | sha256sum
		{
			"skycoin",
			"5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d",
		},
		{
			"hello world",
			"b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			"hello world asd awd awd awdapodawpokawpod ",
			"99d71f95cafe05ea2dddebc35b6083bd5af0e44850c9dc5139b4476c99950be4",
		},
	}
	for _, io := range vals {
		assert.Equal(t, io.output, SumSHA256([]byte(io.input)).Hex())
	}
}

func TestSumSHA256(t *testing.T) {
	b := randBytes(t, 256)
	h1 := SumSHA256(b)
	assert.NotEqual(t, h1, SHA256{})
	// A second call to Sum should not be influenced by the original
	c := randBytes(t, 256)
	h2 := SumSHA256(c)
	assert.NotEqual(t, h2, SHA256{})
	assert.Equal(t, h2, freshSumSHA256(c))
}

func TestSHA256FromHex(t *testing.T) {
	// Invalid hex hash
	_, err := SHA256FromHex("cawcd")
	assert.NotNil(t, err)

	// Truncated hex hash
	h := SumSHA256(randBytes(t, 128))
	_, err = SHA256FromHex(hex.EncodeToString(h[:len(h)/2]))
	assert.NotNil(t, err)

	// Valid hex hash
	h2, err := SHA256FromHex(hex.EncodeToString(h[:]))
	assert.Equal(t, h, h2)
	assert.Nil(t, err)
}

func TestMustSHA256FromHex(t *testing.T) {
	// Invalid hex hash
	assert.Panics(t, func() { MustSHA256FromHex("cawcd") })

	// Truncated hex hash
	h := SumSHA256(randBytes(t, 128))
	assert.Panics(t, func() {
		MustSHA256FromHex(hex.EncodeToString(h[:len(h)/2]))
	})

	// Valid hex hash
	h2 := MustSHA256FromHex(hex.EncodeToString(h[:]))
	assert.Equal(t, h, h2)
}

func TestMustSumSHA256(t *testing.T) {
	b := randBytes(t, 128)
	assert.Panics(t, func() { MustSumSHA256(b, 127) })
	assert.Panics(t, func() { MustSumSHA256(b, 129) })
	assert.NotPanics(t, func() { MustSumSHA256(b, 128) })
	h := MustSumSHA256(b, 128)
	assert.NotEqual(t, h, SHA256{})
	assert.Equal(t, h, freshSumSHA256(b))
}

func TestDoubleSHA256(t *testing.T) {
	b := randBytes(t, 128)
	h := DoubleSHA256(b)
	assert.NotEqual(t, h, SHA256{})
	assert.NotEqual(t, h, freshSumSHA256(b))
}

func TestAddSHA256(t *testing.T) {
	b := randBytes(t, 128)
	h := SumSHA256(b)
	c := randBytes(t, 64)
	i := SumSHA256(c)
	add := AddSHA256(h, i)
	assert.NotEqual(t, add, SHA256{})
	assert.NotEqual(t, add, h)
	assert.NotEqual(t, add, i)
	assert.Equal(t, add, SumSHA256(append(h[:], i[:]...)))
}

func TestXorSHA256(t *testing.T) {
	b := randBytes(t, 128)
	c := randBytes(t, 128)
	h := SumSHA256(b)
	i := SumSHA256(c)
	assert.NotEqual(t, h.Xor(i), h)
	assert.NotEqual(t, h.Xor(i), i)
	assert.NotEqual(t, h.Xor(i), SHA256{})
	assert.Equal(t, h.Xor(i), i.Xor(h))
}

func TestNextPowerOfTwo(t *testing.T) {
	inputs := [][]uint64{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{8, 8},
		{14, 16},
		{16, 16},
		{17, 32},
		{43345, 65536},
		{65535, 65536},
		{35657, 65536},
		{65536, 65536},
		{65537, 131072},
	}
	for _, i := range inputs {
		assert.Equal(t, nextPowerOfTwo(i[0]), i[1])
	}
	for i := uint64(2); i < 10000; i++ {
		p := nextPowerOfTwo(i)
		assert.Equal(t, p%2, uint64(0))
		assert.True(t, p >= i)
	}
}

func TestMerkle(t *testing.T) {
	h := SumSHA256(randBytes(t, 128))
	// Single hash input returns hash
	assert.Equal(t, Merkle([]SHA256{h}), h)
	h2 := SumSHA256(randBytes(t, 128))
	// 2 hashes should be AddSHA256 of them
	assert.Equal(t, Merkle([]SHA256{h, h2}), AddSHA256(h, h2))
	// 3 hashes should be Add(Add())
	h3 := SumSHA256(randBytes(t, 128))
	out := AddSHA256(AddSHA256(h, h2), AddSHA256(h3, SHA256{}))
	assert.Equal(t, Merkle([]SHA256{h, h2, h3}), out)
	// 4 hashes should be Add(Add())
	h4 := SumSHA256(randBytes(t, 128))
	out = AddSHA256(AddSHA256(h, h2), AddSHA256(h3, h4))
	assert.Equal(t, Merkle([]SHA256{h, h2, h3, h4}), out)
	// 5 hashes
	h5 := SumSHA256(randBytes(t, 128))
	out = AddSHA256(AddSHA256(h, h2), AddSHA256(h3, h4))
	out = AddSHA256(out, AddSHA256(AddSHA256(h5, SHA256{}),
		AddSHA256(SHA256{}, SHA256{})))
	assert.Equal(t, Merkle([]SHA256{h, h2, h3, h4, h5}), out)
}
