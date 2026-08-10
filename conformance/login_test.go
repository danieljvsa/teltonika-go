package conformance_test

import (
	"testing"

	teltonika "github.com/danieljvsa/teltonika-go/public"
)

// Without the strict option, login encoding preserves the lenient legacy
// behavior: only empty identifiers are rejected.
func TestLoginValidationLenientDefault(t *testing.T) {
	accept := []string{
		"356307042441013",
		"1234",
		"356307ABC441013",
	}
	for _, imei := range accept {
		if _, err := teltonika.Encode(&teltonika.Packet{Kind: teltonika.KindLogin, IMEI: imei}); err != nil {
			t.Errorf("lenient mode rejected %q: %v", imei, err)
		}
	}

	if _, err := teltonika.Encode(&teltonika.Packet{Kind: teltonika.KindLogin, IMEI: ""}); err == nil {
		t.Fatal("expected empty IMEI error in lenient mode")
	}
}

// WithStrictIMEI login identifiers must be exactly 15 decimal digits.
func TestLoginValidationStrict(t *testing.T) {
	tests := []struct {
		name string
		imei string
	}{
		{name: "empty", imei: ""},
		{name: "too short", imei: "1234"},
		{name: "contains letters", imei: "356307ABC441013"},
		{name: "contains spaces", imei: "356307042441 13"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := &teltonika.Packet{Kind: teltonika.KindLogin, IMEI: tt.imei}
			if _, err := teltonika.Encode(packet, teltonika.WithStrictIMEI()); err == nil {
				t.Fatalf("expected invalid IMEI error for %q", tt.imei)
			}
		})
	}

	valid := &teltonika.Packet{Kind: teltonika.KindLogin, IMEI: "356307042441013"}
	if _, err := teltonika.Encode(valid, teltonika.WithStrictIMEI()); err != nil {
		t.Fatalf("strict mode rejected valid IMEI: %v", err)
	}
}
