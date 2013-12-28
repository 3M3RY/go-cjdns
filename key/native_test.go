package key

import (
	"net"
	"testing"
)

var (
	// From the cjdns build tests in test/CryptoAddress_test.c

	pubkeyString = "r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k"
	pubkeyBytes  = []byte{
		0xd7, 0xc0, 0xdf, 0x45, 0x00, 0x1a, 0x5b, 0xe5,
		0xe8, 0x1c, 0x95, 0xe5, 0x19, 0xbe, 0x51, 0x99,
		0x05, 0x52, 0x37, 0xcb, 0x91, 0x16, 0x88, 0x2c,
		0xad, 0xce, 0xfe, 0x48, 0xab, 0x73, 0x51, 0x73,
	}
	pubkeyIPv6 = "fc68:cb2c:60db:cb96:19ac:34a8:fd34:03fc"

	privkeyString = "751d3db85b848deaf221e0ed2b6cc17f587b29057d74cdd4dc0bd18b7157288e"
	privkeyBytes  = []byte{
		0x75, 0x1d, 0x3d, 0xb8, 0x5b, 0x84, 0x8d, 0xea,
		0xf2, 0x21, 0xe0, 0xed, 0x2b, 0x6c, 0xc1, 0x7f,
		0x58, 0x7b, 0x29, 0x05, 0x7d, 0x74, 0xcd, 0xd4,
		0xdc, 0x0b, 0xd1, 0x8b, 0x71, 0x57, 0x28, 0x8e,
	}

	randK = Generate()
)

func TestGenerate(t *testing.T) {
	key := Generate()
	keyPublic := key.Pubkey()

	if !keyPublic.Valid() {
		t.Error("generated public key", keyPublic, "not valid")
	}
	if keyPublic.IP()[0] != 0xFC {
		t.Error("generated public key", keyPublic, "has invalid IP address", keyPublic.IP())
	}
}

func TestDecodePrivate(t *testing.T) {
	key, err := DecodePrivate(privkeyString)
	if err != nil {
		t.Error(err)
		return
	}
	if !key.Valid() {
		t.Error("failed to decode valid private key", privkeyString, "got", key)
	}
	keyPublic := key.Pubkey()
	if !keyPublic.Valid() {
		t.Error("failed to generate valid public key from private key, got", keyPublic)
	}

	ip := keyPublic.IP()
	if ip[0] != 0xFC {
		t.Error("key", keyPublic, "generated from a private key has invalid IP address", keyPublic.IP())
	}
	if !ip.Equal(net.ParseIP(pubkeyIPv6)) {
		t.Error("IP address should be", pubkeyIPv6, "got", ip)
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	// Private
	privateA := Generate()
	b, err := privateA.MarshalText()
	if err != nil {
		t.Error("failed to MarshalText on", privateA, err)
	}
	sA, sB := privateA.String(), string(b)
	if sA != sB {
		t.Error("failed to MarshalText on", sA, "got", sB)
	}

	privateB := new(Private)
	err = privateB.UnmarshalText(b)
	if err != nil {
		t.Error("failed to UnMarshalText of", privateA, err)
	}
	sB = privateB.String()
	if sA != sB {
		t.Error("private key unmarshaling failed,", privateA, "and", privateB, "do not match")
	}

	// Public
	publicA := privateA.Pubkey()
	b, err = publicA.MarshalText()
	if err != nil {
		t.Error("failed to MarshalText on", publicA, err)
	}

	sA, sB = publicA.String(), string(b)
	if sA != sB {
		t.Error("failed to MarshalText on", sA, "got", sB)
	}

	publicB := new(Public)
	err = publicB.UnmarshalText(b)
	if err != nil {
		t.Error("failed to UnMarshalText of", publicA, err)
	}
	sB = publicB.String()
	if sA != sB {
		t.Error("public key unmarshaling failed,", sA, "and", sB, "do not match")
	}
}

func TestDecodePublic(t *testing.T) {
	keyPublic, err := DecodePublic(pubkeyString)
	if err != nil {
		t.Error(err)
		return
	}
	if !keyPublic.Valid() {
		t.Error("decoded public key", keyPublic, "not valid")
	}
	ip := keyPublic.IP()
	if ip[0] != 0xFC {
		t.Error("decoded public key", keyPublic, "has invalid IP address", keyPublic.IP())
	}
	if !ip.Equal(net.ParseIP(pubkeyIPv6)) {
		t.Error("IP address should be", pubkeyIPv6, "got", ip)
	}
}

func BenchmarkPrivateMarshalText(b *testing.B) {
	k := Generate()
	for i := 0; i < b.N; i++ {
		k.MarshalText()
	}
}

func BenchmarkPrivateUnmarshalText(b *testing.B) {
	buf := []byte("751d3db85b848deaf221e0ed2b6cc17f587b29057d74cdd4dc0bd18b7157288e")
	k := new(Private)
	for i := 0; i < b.N; i++ {
		k.UnmarshalText(buf)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate()
	}
}

func BenchmarkPubkey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randK.Pubkey()
	}
}

func BenchmarkPublicMarshalText(b *testing.B) {
	pk := randK.Pubkey()
	for i := 0; i < b.N; i++ {
		pk.MarshalText()
	}
}

func BenchmarkPublicUnmarshalText(b *testing.B) {
	buf := []byte("r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k")
	pk := new(Public)
	for i := 0; i < b.N; i++ {
		pk.UnmarshalText(buf)
	}
}
