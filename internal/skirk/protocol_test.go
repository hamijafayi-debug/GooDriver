package skirk

import (
	"bytes"
	"testing"
)

func TestSealOpenEnvelope(t *testing.T) {
	key, err := DeriveKey("test-secret")
	if err != nil {
		t.Fatal(err)
	}
	sid, err := ParseSessionID("00112233445566778899aabbccddeeff")
	if err != nil {
		t.Fatal(err)
	}
	plaintext := []byte("hello skirk")
	sealed, err := Seal(key, sid, DirectionUp, 7, plaintext, true)
	if err != nil {
		t.Fatal(err)
	}
	env, opened, err := OpenEnvelope(key, sealed)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(opened, plaintext) {
		t.Fatalf("plaintext mismatch: got %q", opened)
	}
	if env.SessionID != sid || env.Direction != DirectionUp || env.Sequence != 7 || env.Flags != FlagFinal {
		t.Fatalf("metadata mismatch: %+v", env)
	}
}

func TestOpenEnvelopeRejectsTamper(t *testing.T) {
	key, err := DeriveKey("test-secret")
	if err != nil {
		t.Fatal(err)
	}
	sid, _ := ParseSessionID("00112233445566778899aabbccddeeff")
	sealed, err := Seal(key, sid, DirectionUp, 1, []byte("payload"), false)
	if err != nil {
		t.Fatal(err)
	}
	sealed[len(sealed)-1] ^= 0x01
	if _, _, err := OpenEnvelope(key, sealed); err == nil {
		t.Fatal("expected tampered envelope to fail authentication")
	}
}

func TestDeriveMuxLaneKeyV4SeparatesClientsRunsAndDirections(t *testing.T) {
	sid, _ := ParseSessionID("00112233445566778899aabbccddeeff")
	upA, err := DeriveMuxLaneKeyV4("test-secret", sid, DirectionUp, "client-a", "run-a", 0)
	if err != nil {
		t.Fatal(err)
	}
	upA2, err := DeriveMuxLaneKeyV4("test-secret", sid, DirectionUp, "client-a", "run-a", 0)
	if err != nil {
		t.Fatal(err)
	}
	upB, err := DeriveMuxLaneKeyV4("test-secret", sid, DirectionUp, "client-b", "run-a", 0)
	if err != nil {
		t.Fatal(err)
	}
	upRunB, err := DeriveMuxLaneKeyV4("test-secret", sid, DirectionUp, "client-a", "run-b", 0)
	if err != nil {
		t.Fatal(err)
	}
	downA, err := DeriveMuxLaneKeyV4("test-secret", sid, DirectionDown, "client-a", "run-a", 0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(upA, upA2) {
		t.Fatal("same mux v4 inputs must derive the same key")
	}
	if bytes.Equal(upA, upB) || bytes.Equal(upA, upRunB) || bytes.Equal(upA, downA) {
		t.Fatal("mux v4 keys must differ by client id, run id, and direction")
	}
}

func TestNonceUsesFullSID(t *testing.T) {
	// Verify that every byte of the 16-byte SID influences the nonce.
	// The original implementation omitted sid[5..10], meaning sessions whose
	// SIDs differed only in those bytes produced identical nonces — a serious
	// AES-GCM nonce-reuse vulnerability.
	var sid [16]byte
	for i := range sid {
		sid[i] = byte(i + 1)
	}
	base := nonce(sid, 0, 0)

	// Flip each byte of the SID independently and confirm the nonce changes.
	for i := 0; i < 16; i++ {
		modified := sid
		modified[i] ^= 0xFF
		n := nonce(modified, 0, 0)
		if bytes.Equal(base, n) {
			t.Fatalf("nonce must change when byte %d of SID changes (was identical — byte not XOR-folded)", i)
		}
	}

	// Verify direction and sequence still differentiate nonces.
	nUp := nonce(sid, DirectionUp, 1)
	nDown := nonce(sid, DirectionDown, 1)
	if bytes.Equal(nUp, nDown) {
		t.Fatal("nonce must differ between DirectionUp and DirectionDown")
	}
	nSeq1 := nonce(sid, DirectionUp, 1)
	nSeq2 := nonce(sid, DirectionUp, 2)
	if bytes.Equal(nSeq1, nSeq2) {
		t.Fatal("nonce must differ between sequence 1 and 2")
	}
}
