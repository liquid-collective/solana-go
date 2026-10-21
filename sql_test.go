package solana

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLScanValueRoundTrip(t *testing.T) {
	pk := MustPublicKeyFromBase58("SPoo1Ku8WFXoNDMHPsrGSTSG1Y47rzgn41SLUNakuHy")
	hash := MustHashFromBase58("4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM")
	sig := MustSignatureFromBase58("5VERv8NMvzbJMEkV8xnrLkEaWRtSz9CosKDYjCJjBRnbJLgp8uirBgmQpjKhoR4tjF3ZpRzrFmBV6UjKdiSZkQUW")

	cases := []struct {
		name string
		val  driver.Valuer
		dst  interface{ Scan(any) error }
		want any
	}{
		{"PublicKey", pk, new(PublicKey), pk},
		{"Hash", hash, new(Hash), hash},
		{"Signature", sig, new(Signature), sig},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := tc.val.Value()
			require.NoError(t, err)
			require.IsType(t, "", v)
			require.NoError(t, tc.dst.Scan(v))
			require.EqualValues(t, tc.want, derefScanned(tc.dst))
		})
	}
}

func TestSQLScanAcceptsByteBackedText(t *testing.T) {
	pk := MustPublicKeyFromBase58("SPoo1Ku8WFXoNDMHPsrGSTSG1Y47rzgn41SLUNakuHy")
	hash := MustHashFromBase58("4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM")
	sig := MustSignatureFromBase58("5VERv8NMvzbJMEkV8xnrLkEaWRtSz9CosKDYjCJjBRnbJLgp8uirBgmQpjKhoR4tjF3ZpRzrFmBV6UjKdiSZkQUW")

	var gotPK PublicKey
	require.NoError(t, gotPK.Scan([]byte(pk.String())))
	require.Equal(t, pk, gotPK)
	var gotHash Hash
	require.NoError(t, gotHash.Scan([]byte(hash.String())))
	require.Equal(t, hash, gotHash)
	var gotSig Signature
	require.NoError(t, gotSig.Scan([]byte(sig.String())))
	require.Equal(t, sig, gotSig)
}

func TestSQLScanRejectsOtherTypes(t *testing.T) {
	require.Error(t, new(PublicKey).Scan(7))
	require.Error(t, new(Hash).Scan(42))
	require.Error(t, new(Signature).Scan(nil))
}

func TestSQLScanRejectsInvalidBase58(t *testing.T) {
	require.Error(t, new(PublicKey).Scan("not-a-pubkey"))
}

func derefScanned(dst any) any {
	switch d := dst.(type) {
	case *PublicKey:
		return *d
	case *Hash:
		return *d
	case *Signature:
		return *d
	}
	return nil
}
