package stakepool

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testKey(i int) ag_solanago.PublicKey {
	var k ag_solanago.PublicKey
	k[0] = byte(i)
	k[31] = 7
	return k
}

// fill calls every setter that takes a single public key with a distinct key, the way a caller
// provides each account, and returns the builder. Optional-account setters are left alone so a
// test opts into them explicitly.
func fill[T any](b T) T {
	optional := map[string]string{
		"Initialize":                 "SetDepositAuthorityAccount",
		"WithdrawSol":                "SetSolWithdrawAuthority",
		"UpdateValidatorListBalance": "SetValidatorStakeAccount",
		"DepositStake":               "SetStakeDepositAuthoritySigner",
		"DepositStakeWithSlippage":   "SetStakeDepositAuthoritySigner",
	}
	skip := []string{optional[reflect.TypeOf(b).Elem().Name()]}
	v := reflect.ValueOf(b)
	var names []string
	for i := 0; i < v.NumMethod(); i++ {
		m := v.Type().Method(i)
		if !strings.HasPrefix(m.Name, "Set") || m.Name == "SetAccounts" || m.Name == "SetStakePoolProgramID" {
			continue
		}
		if mt := v.Method(i).Type(); mt.NumIn() == 1 && mt.In(0) == reflect.TypeOf(ag_solanago.PublicKey{}) {
			names = append(names, m.Name)
		}
	}
	sort.Strings(names)
	n := 1
	for _, name := range names {
		skipped := false
		for _, s := range skip {
			skipped = skipped || s == name
		}
		if skipped {
			continue
		}
		v.MethodByName(name).Call([]reflect.Value{reflect.ValueOf(testKey(n))})
		n++
	}
	return b
}

func signerIndexes(accounts []*ag_solanago.AccountMeta) []int {
	out := []int{}
	for i, a := range accounts {
		if a.IsSigner {
			out = append(out, i)
		}
	}
	return out
}

// Expected data is canonical Borsh computed outside this package from the SPL definitions
// (https://docs.rs/spl-stake-pool/2.0.3/src/spl_stake_pool/instruction.rs.html).
func TestInstructionsMatchSPLEncoding(t *testing.T) {
	newPK := testKey(9)
	funding := FundingTypeSolDeposit
	cases := []struct {
		name     string
		inst     *Instruction
		wantHex  string
		accounts int
		signers  []int
	}{
		{"AddValidatorToPool", fill(NewAddValidatorToPoolBuilder()).SetOptionalSeed(u32(7)).Build(), "0107000000", 13, []int{1}},
		{"CleanupRemovedValidatorEntries", fill(NewCleanupRemovedValidatorEntriesInstructionBuilder()).Build(), "08", 2, []int{}},
		{"CreateTokenMetadata", fill(NewCreateTokenMetadataInstructionBuilder()).SetName("A").SetSymbol("B").SetURI("C").Build(), "11010000004101000000420100000043", 8, []int{1, 4}},
		{"CreateTokenMetadata long", fill(NewCreateTokenMetadataInstructionBuilder()).SetName("Alluvial LsSOL").SetSymbol("LsSOL").SetURI("https://example.com/m.json").Build(), "110e000000416c6c757669616c204c73534f4c050000004c73534f4c1a00000068747470733a2f2f6578616d706c652e636f6d2f6d2e6a736f6e", 8, []int{1, 4}},
		{"UpdateTokenMetadata", fill(NewUpdateTokenMetadataInstructionBuilder()).SetName("A").SetSymbol("B").SetURI("C").Build(), "12010000004101000000420100000043", 5, []int{1}},
		{"DecreaseAdditionalValidatorStake", fill(NewDecreaseAdditionalValidatorStakeBuilder()).SetArgs(&AdditionalValidatorStakeArgs{Lamports: 1000, TransientStakeSeed: 2, EphemeralStakeSeed: 3}).Build(), "14e80300000000000002000000000000000300000000000000", 12, []int{1}},
		{"DecreaseValidatorStake", fill(NewDecreaseValidatorStakeInstructionBuilder()).SetLamports(1000).SetTransientStakeSeed(2).Build(), "03e8030000000000000200000000000000", 10, []int{1}},
		{"DecreaseValidatorStakeWithReserve", fill(NewDecreaseValidatorStakeWithReserveWithReserveInstructionBuilder()).SetLamports(1000).SetTransientStakeSeed(2).Build(), "15e8030000000000000200000000000000", 11, []int{1}},
		{"DepositSol", fill(NewDepositSolBuilder()).SetLamports(1_000_000_000).Build(), "0e00ca9a3b00000000", 10, []int{3}},
		{"DepositSolWithSlippage", fill(NewDepositSolWithSlippageBuilder()).SetLamportsIn(1000).SetMinTokensOut(900).Build(), "19e8030000000000008403000000000000", 10, []int{3}},
		{"DepositStake permissionless", fill(NewDepositStakeBuilder()).Build(), "09", 15, []int{}},
		{"DepositStake private pool", fill(NewDepositStakeBuilder()).SetStakeDepositAuthoritySigner(newPK).Build(), "09", 15, []int{2}},
		{"DepositStakeWithSlippage", fill(NewDepositStakeWithSlippageBuilder()).SetMinTokensOut(900).Build(), "178403000000000000", 15, []int{}},
		{"IncreaseAdditionalValidatorStake", fill(NewIncreaseAdditionalValidatorStakeInstructionBuilder()).SetArgs(&AdditionalValidatorStakeArgs{Lamports: 1000, TransientStakeSeed: 2, EphemeralStakeSeed: 3}).Build(), "13e80300000000000002000000000000000300000000000000", 14, []int{1}},
		{"IncreaseValidatorStake", fill(NewIncreaseValidatorStakeInstructionBuilder()).SetLamports(1000).SetTransientStakeSeed(2).Build(), "04e8030000000000000200000000000000", 14, []int{1}},
		{"Initialize permissionless", fill(NewInitializeInstructionBuilder()).SetFee(Fee{Denominator: 10, Numerator: 1}).SetWithdrawalFee(Fee{Denominator: 100, Numerator: 2}).SetDepositFee(Fee{Denominator: 1000, Numerator: 3}).SetReferralFee(5).SetMaxValidators(2000).Build(), "000a00000000000000010000000000000064000000000000000200000000000000e803000000000000030000000000000005d0070000", 9, []int{1}},
		{"Initialize private pool", fill(NewInitializeInstructionBuilder()).SetFee(Fee{Denominator: 10, Numerator: 1}).SetWithdrawalFee(Fee{Denominator: 100, Numerator: 2}).SetDepositFee(Fee{Denominator: 1000, Numerator: 3}).SetReferralFee(5).SetMaxValidators(2000).SetDepositAuthorityAccount(newPK).Build(), "000a00000000000000010000000000000064000000000000000200000000000000e803000000000000030000000000000005d0070000", 10, []int{1, 9}},
		{"Redelegate", fill(NewRedelegateBuilder()).SetLamports(1000).SetSourceTransientStakeSeed(2).SetEphemeralStakeSeed(3).SetDestinationTransientStakeSeed(4).Build(), "16e803000000000000020000000000000003000000000000000400000000000000", 15, []int{1}},
		{"RemoveValidatorFromPool", fill(NewRemoveValidatorFromPoolInstructionBuilder()).Build(), "02", 8, []int{1}},
		{"SetFee epoch", fill(NewSetFeeInstructionBuilder()).SetFee(Epoch{Fee: Fee{Denominator: 100, Numerator: 5}}).Build(), "0c0264000000000000000500000000000000", 2, []int{1}},
		{"SetFee sol referral", fill(NewSetFeeInstructionBuilder()).SetFee(SolReferral{Fee: 10}).Build(), "0c000a", 2, []int{1}},
		{"SetFee stake referral", fill(NewSetFeeInstructionBuilder()).SetFee(StakeReferral{Fee: 20}).Build(), "0c0114", 2, []int{1}},
		{"SetFee stake withdrawal", fill(NewSetFeeInstructionBuilder()).SetFee(StakeWithdrawal{Fee: Fee{Denominator: 1, Numerator: 200}}).Build(), "0c030100000000000000c800000000000000", 2, []int{1}},
		{"SetFee sol deposit", fill(NewSetFeeInstructionBuilder()).SetFee(SolDeposit{Fee: Fee{Denominator: 1, Numerator: 300}}).Build(), "0c0401000000000000002c01000000000000", 2, []int{1}},
		{"SetFee stake deposit", fill(NewSetFeeInstructionBuilder()).SetFee(StakeDeposit{Fee: Fee{Denominator: 1, Numerator: 400}}).Build(), "0c0501000000000000009001000000000000", 2, []int{1}},
		{"SetFee sol withdrawal", fill(NewSetFeeInstructionBuilder()).SetFee(SolWithdrawal{Fee: Fee{Denominator: 1, Numerator: 500}}).Build(), "0c060100000000000000f401000000000000", 2, []int{1}},
		{"SetFundingAuthority", fill(NewSetFundingAuthorityInstructionBuilder()).SetAuth(&funding).Build(), "0f01", 3, []int{1}},
		{"SetManager", fill(NewSetManagerInstructionBuilder()).Build(), "0b", 4, []int{1, 2}},
		{"SetStaker", fill(NewSetStakerInstructionBuilder()).Build(), "0d", 3, []int{1}},
		{"SetPreferredValidator some", fill(NewSetPreferredValidatorInstructionBuilder()).SetValidatorType(PreferredValidatorTypeDeposit).SetValidatorVoteAddress(&newPK).Build(), "0500010900000000000000000000000000000000000000000000000000000000000007", 3, []int{1}},
		{"SetPreferredValidator none", fill(NewSetPreferredValidatorInstructionBuilder()).SetValidatorType(PreferredValidatorTypeWithdraw).Build(), "050100", 3, []int{1}},
		{"UpdateStakePoolBalance", fill(NewUpdateStakePoolBalanceInstructionBuilder()).Build(), "07", 7, []int{}},
		{"UpdateValidatorListBalance", fill(NewUpdateValidatorListBalanceInstructionBuilder()).SetArgs(UpdateValidatorListBalanceArgs{StartIndex: 3, NoMerge: true}).Build(), "060300000001", 7, []int{}},
		{"WithdrawSol", fill(NewWithdrawSolBuilder()).SetArg(1000).Build(), "10e803000000000000", 12, []int{2}},
		{"WithdrawSol restricted pool", fill(NewWithdrawSolBuilder()).SetArg(1000).SetSolWithdrawAuthority(newPK).Build(), "10e803000000000000", 13, []int{2, 12}},
		{"WithdrawSolWithSlippage", fill(NewWithdrawSolWithSlippageInstructionBuilder()).SetTokensIn(1000).SetMinLamportsOut(900).Build(), "1ae8030000000000008403000000000000", 12, []int{2}},
		{"WithdrawStake", fill(NewWithdrawStakeBuilder()).SetAmount(1000).Build(), "0ae803000000000000", 13, []int{6}},
		{"WithdrawStakeWithSlippage", fill(NewWithdrawStakeWithSlippageInstructionBuilder()).SetPoolTokensIn(1000).SetMinLamportsOut(900).Build(), "18e8030000000000008403000000000000", 13, []int{6}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.inst.Data()
			require.NoError(t, err)
			assert.Equal(t, tc.wantHex, hex.EncodeToString(data))

			accounts := tc.inst.Accounts()
			assert.Len(t, accounts, tc.accounts)
			assert.Equal(t, tc.signers, signerIndexes(accounts))

			// Decoding the canonical bytes with the accounts yields the same instruction.
			decoded, err := DecodeInstruction(accounts, data)
			require.NoError(t, err)
			assert.Equal(t, tc.inst.Impl, decoded.Impl)
			assert.Equal(t, accounts, decoded.Accounts())
		})
	}
}

func u32(v uint32) *uint32 { return &v }

// Real instructions from our devnet stake pools, decoded and re-encoded unchanged.
func TestInstructionsRoundTripOnChainData(t *testing.T) {
	raw, err := os.ReadFile("testdata/onchain_instructions.json")
	require.NoError(t, err)
	var fixtures map[string][]struct {
		Accounts int    `json:"accounts"`
		Data     string `json:"data"`
	}
	require.NoError(t, json.Unmarshal(raw, &fixtures))
	require.NotEmpty(t, fixtures)

	for name, list := range fixtures {
		for _, f := range list {
			t.Run(name, func(t *testing.T) {
				data, err := hex.DecodeString(f.Data)
				require.NoError(t, err)
				accounts := make([]*ag_solanago.AccountMeta, f.Accounts)
				for i := range accounts {
					accounts[i] = ag_solanago.Meta(testKey(i + 1))
				}

				decoded, err := DecodeInstruction(accounts, data)
				require.NoError(t, err)
				assert.Equal(t, name, InstructionIDToName(decoded.TypeID.Uint8()))
				assert.Len(t, decoded.Accounts(), f.Accounts)

				again, err := decoded.Data()
				require.NoError(t, err)
				assert.Equal(t, f.Data, hex.EncodeToString(again))
			})
		}
	}
}

func TestInitializeOptionalDepositAuthority(t *testing.T) {
	b := fill(NewInitializeInstructionBuilder())
	require.Nil(t, b.GetDepositAuthorityAccount())
	require.Len(t, b.Accounts, 9, "an absent authority is not an account")

	b.SetDepositAuthorityAccount(testKey(40))
	require.Len(t, b.Accounts, 10)
	assert.True(t, b.GetDepositAuthorityAccount().IsSigner)
}

func TestInstructionValidationRejectsMissingParameters(t *testing.T) {
	_, err := fill(NewWithdrawStakeBuilder()).ValidateAndBuild()
	assert.ErrorContains(t, err, "amount is not set")

	_, err = fill(NewDepositSolBuilder()).ValidateAndBuild()
	assert.ErrorContains(t, err, "lamports is not set")

	b := fill(NewWithdrawStakeBuilder()).SetAmount(1)
	b.Accounts[6] = ag_solanago.Meta(testKey(60)) // user transfer authority without signing
	_, err = b.ValidateAndBuild()
	assert.ErrorContains(t, err, "should be a signer")
}

func TestDecodeRejectsTruncatedData(t *testing.T) {
	accounts := fill(NewWithdrawStakeBuilder()).Accounts
	for _, data := range [][]byte{{0x0a}, {0x0a, 0x01, 0x02}, {0x11, 0x05, 0x00, 0x00, 0x00, 'a'}} {
		_, err := DecodeInstruction(accounts, data)
		assert.Error(t, err, "%x", data)
	}
}

func TestAccountPrivileges(t *testing.T) {
	t.Run("sol deposit source is writable", func(t *testing.T) {
		for _, accounts := range [][]*ag_solanago.AccountMeta{
			fill(NewDepositSolBuilder()).Accounts,
			fill(NewDepositSolWithSlippageBuilder()).Accounts,
		} {
			assert.True(t, accounts[3].IsSigner)
			assert.True(t, accounts[3].IsWritable)
		}
	})
	t.Run("stake deposit authority signs only for private pools", func(t *testing.T) {
		assert.False(t, fill(NewDepositStakeBuilder()).Accounts[2].IsSigner)
		assert.False(t, fill(NewDepositStakeWithSlippageBuilder()).Accounts[2].IsSigner)
		assert.True(t, NewDepositStakeBuilder().SetStakeDepositAuthoritySigner(testKey(1)).Accounts[2].IsSigner)
		assert.True(t, NewDepositStakeWithSlippageBuilder().SetStakeDepositAuthoritySigner(testKey(1)).Accounts[2].IsSigner)
	})
	t.Run("metadata payer signs and pays", func(t *testing.T) {
		payer := fill(NewCreateTokenMetadataInstructionBuilder()).Accounts[4]
		assert.True(t, payer.IsSigner)
		assert.True(t, payer.IsWritable)
	})
}

func TestUpdateValidatorListBalanceAppendsValidatorAccounts(t *testing.T) {
	b := fill(NewUpdateValidatorListBalanceInstructionBuilder()).SetArgs(UpdateValidatorListBalanceArgs{})
	require.Len(t, b.Accounts, 7)
	assert.Equal(t, ag_solanago.PublicKey{}, b.GetValidatorStakeAccount())

	b.SetValidatorStakeAccount(testKey(70))
	b.SetValidatorAndTransientAccounts([]ag_solanago.PublicKey{testKey(71), testKey(72)})
	require.Len(t, b.Accounts, 10)
	assert.Equal(t, testKey(70), b.GetValidatorStakeAccount())
	_, err := b.ValidateAndBuild()
	assert.NoError(t, err)
}

func TestSetProgramIDReturnsRegistrationResult(t *testing.T) {
	original := ProgramID
	t.Cleanup(func() { ProgramID = original })

	other := testKey(80)
	require.NoError(t, SetProgramID(other))
	assert.Equal(t, other, ProgramID)
	assert.NoError(t, SetProgramID(other), "registering the same decoder again is a no-op")
}
