package token

import (
	"encoding/hex"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Expected bytes follow the SPL Token instruction layout
// (https://github.com/solana-program/token/blob/main/interface/src/instruction.rs):
// a one-byte discriminator, then the parameters. UiAmountToAmount carries its string as raw
// UTF-8 with no length prefix; AmountToUiAmount carries a little-endian u64.
func TestAddedInstructionsMatchSPLEncoding(t *testing.T) {
	mint := solana.MustPublicKeyFromBase58("11111111111111111111111111111115")

	cases := []struct {
		name     string
		inst     *Instruction
		wantHex  string
		writable bool
	}{
		{"GetAccountDataSize", NewGetAccountDataSizeInstructionBuilder().SetMint(mint).Build(), "15", false},
		{"InitializeImmutableOwner", NewInitializeImmutableOwnerInstructionBuilder().SetAccount(mint).Build(), "16", true},
		{"AmountToUiAmount", NewAmountToUiAmountInstructionBuilder().SetAmount(1_000_000).SetMintAccount(mint).Build(), "1740420f0000000000", false},
		{"UiAmountToAmount", NewUiAmountToAmountInstructionBuilder().SetUiAmount("1.25").SetMintAccount(mint).Build(), "18312e3235", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.inst.Data()
			require.NoError(t, err)
			assert.Equal(t, tc.wantHex, hex.EncodeToString(data))

			accounts := tc.inst.Accounts()
			require.Len(t, accounts, 1)
			assert.Equal(t, mint, accounts[0].PublicKey)
			assert.Equal(t, tc.writable, accounts[0].IsWritable)

			decoded, err := DecodeInstruction(accounts, data)
			require.NoError(t, err)
			assert.Equal(t, tc.inst.Impl, decoded.Impl)
		})
	}
}

func TestUiAmountToAmountRejectsInvalidUTF8(t *testing.T) {
	_, err := DecodeInstruction([]*solana.AccountMeta{solana.Meta(solana.SystemProgramID)}, []byte{Instruction_UiAmountToAmount, 0xff, 0xfe})
	assert.ErrorContains(t, err, "UTF-8")
}
