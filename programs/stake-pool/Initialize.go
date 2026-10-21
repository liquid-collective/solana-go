// Copyright 2021 github.com/gagliardetto
// Copyright 2025 github.com/liquid-collective
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stakepool

import (
	"errors"
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	ag_treeout "github.com/gagliardetto/treeout"

	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
)

type Initialize struct {
	Fee           *Fee
	WithdrawalFee *Fee
	DepositFee    *Fee
	ReferralFee   *uint8
	MaxValidators *uint32
	// [0] = [WRITE] stakePool
	// [1] = [SIGNER] manager
	// [2] = [] staker
	// [3] = [] withdrawAuthority
	// [4] = [WRITE] validatorList
	// [5] = [] reserveStake
	// [6] = [WRITE] poolMint
	// [7] = [WRITE] managerFeeAccount
	// [8] = [] tokenProgram
	// [9] = [SIGNER] depositAuthority (optional, omitted for a permissionless pool)
	Accounts ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func NewInitializeInstructionBuilder() *Initialize {
	return &Initialize{
		Accounts: make(ag_solanago.AccountMetaSlice, 9),
	}
}

func (i *Initialize) SetStakePoolAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[0] = ag_solanago.Meta(account).WRITE()
	return i
}

func (i *Initialize) GetStakePoolAccount() *ag_solanago.AccountMeta {
	return i.Accounts[0]
}

func (i *Initialize) SetManagerAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[1] = ag_solanago.Meta(account).SIGNER()
	return i
}

func (i *Initialize) GetManagerAccount() *ag_solanago.AccountMeta {
	return i.Accounts[1]
}

func (i *Initialize) SetStakerAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[2] = ag_solanago.Meta(account)
	return i
}

func (i *Initialize) GetStakerAccount() *ag_solanago.AccountMeta {
	return i.Accounts[2]
}

func (i *Initialize) SetWithdrawAuthorityAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[3] = ag_solanago.Meta(account)
	return i
}

func (i *Initialize) GetWithdrawAuthorityAccount() *ag_solanago.AccountMeta {
	return i.Accounts[3]
}

func (i *Initialize) SetValidatorListAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[4] = ag_solanago.Meta(account).WRITE()
	return i
}

func (i *Initialize) GetValidatorListAccount() *ag_solanago.AccountMeta {
	return i.Accounts[4]
}

func (i *Initialize) SetReserveStakeAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[5] = ag_solanago.Meta(account)
	return i
}

func (i *Initialize) GetReserveStakeAccount() *ag_solanago.AccountMeta {
	return i.Accounts[5]
}

func (i *Initialize) SetPoolMintAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[6] = ag_solanago.Meta(account).WRITE()
	return i
}

func (i *Initialize) GetPoolMintAccount() *ag_solanago.AccountMeta {
	return i.Accounts[6]
}

func (i *Initialize) SetManagerFeeAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[7] = ag_solanago.Meta(account).WRITE()
	return i
}

func (i *Initialize) GetManagerFeeAccount() *ag_solanago.AccountMeta {
	return i.Accounts[7]
}

func (i *Initialize) SetTokenProgramAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts[8] = ag_solanago.Meta(account)
	return i
}

func (i *Initialize) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return i.Accounts[8]
}

// SetDepositAuthorityAccount makes deposits require this authority's signature (a private pool).
// Leave it unset for a permissionless pool.
func (i *Initialize) SetDepositAuthorityAccount(account ag_solanago.PublicKey) *Initialize {
	i.Accounts = append(i.Accounts[:9], ag_solanago.Meta(account).SIGNER())
	return i
}

// GetDepositAuthorityAccount returns nil for a permissionless pool.
func (i *Initialize) GetDepositAuthorityAccount() *ag_solanago.AccountMeta {
	if len(i.Accounts) <= 9 {
		return nil
	}
	return i.Accounts[9]
}

func (i *Initialize) ValidateAndBuild() (*Instruction, error) {
	if err := i.Validate(); err != nil {
		return nil, err
	}
	return i.Build(), nil
}

func (i *Initialize) Build() *Instruction {
	return &Instruction{
		BaseVariant: ag_binary.BaseVariant{
			TypeID: ag_binary.TypeIDFromUint8(Instruction_Initialize),
			Impl:   i,
		},
	}
}

func (i *Initialize) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Initialize")).
				ParentFunc(func(instructionBranch ag_treeout.Branches) {
					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						if i.Fee != nil {
							paramsBranch.Child(ag_format.Param("Fee", *i.Fee))
						}
						if i.WithdrawalFee != nil {
							paramsBranch.Child(ag_format.Param("WithdrawalFee", *i.WithdrawalFee))
						}
						if i.DepositFee != nil {
							paramsBranch.Child(ag_format.Param("DepositFee", *i.DepositFee))
						}
						if i.ReferralFee != nil {
							paramsBranch.Child(ag_format.Param("ReferralFee", *i.ReferralFee))
						}
						if i.MaxValidators != nil {
							paramsBranch.Child(ag_format.Param("MaxValidators", *i.MaxValidators))
						}
					})

					// Accounts of the instructions:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("stakePool", i.Accounts[0]))
						accountsBranch.Child(ag_format.Meta("manager", i.Accounts[1]))
						accountsBranch.Child(ag_format.Meta("staker", i.Accounts[2]))
						accountsBranch.Child(ag_format.Meta("withdrawAuthority", i.Accounts[3]))
						accountsBranch.Child(ag_format.Meta("validatorList", i.Accounts[4]))
						accountsBranch.Child(ag_format.Meta("reserveStake", i.Accounts[5]))
						accountsBranch.Child(ag_format.Meta("poolMint", i.Accounts[6]))
						accountsBranch.Child(ag_format.Meta("managerFeeAccount", i.Accounts[7]))
						accountsBranch.Child(ag_format.Meta("tokenProgram", i.Accounts[8]))
						if len(i.Accounts) > 9 {
							accountsBranch.Child(ag_format.Meta("depositAuthority", i.Accounts[9]))
						}
					})
				})
		})
}

func (i *Initialize) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	for name, fee := range map[string]*Fee{"fee": i.Fee, "withdrawalFee": i.WithdrawalFee, "depositFee": i.DepositFee} {
		if fee == nil {
			return fmt.Errorf("%s is not set", name)
		}
	}
	if i.ReferralFee == nil {
		return errors.New("referralFee is not set")
	}
	if i.MaxValidators == nil {
		return errors.New("maxValidators is not set")
	}
	for _, fee := range []*Fee{i.Fee, i.WithdrawalFee, i.DepositFee} {
		if err := writeFee(encoder, *fee); err != nil {
			return err
		}
	}
	if err := encoder.WriteUint8(*i.ReferralFee); err != nil {
		return err
	}
	return encoder.WriteUint32(*i.MaxValidators, ag_binary.LE)
}

func (i *Initialize) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	var fees [3]Fee
	for j := range fees {
		fee, err := readFee(decoder)
		if err != nil {
			return err
		}
		fees[j] = fee
	}
	referralFee, err := decoder.ReadUint8()
	if err != nil {
		return err
	}
	maxValidators, err := decoder.ReadUint32(ag_binary.LE)
	if err != nil {
		return err
	}
	i.Fee, i.WithdrawalFee, i.DepositFee = &fees[0], &fees[1], &fees[2]
	i.ReferralFee = &referralFee
	i.MaxValidators = &maxValidators
	return nil
}

func (i *Initialize) Validate() error {
	if i.Fee == nil {
		return errors.New("fee is not set")
	}
	if i.WithdrawalFee == nil {
		return errors.New("withdrawalFee is not set")
	}
	if i.DepositFee == nil {
		return errors.New("depositFee is not set")
	}
	if i.ReferralFee == nil {
		return errors.New("referralFee is not set")
	}
	if i.MaxValidators == nil {
		return errors.New("maxValidators is not set")
	}
	for j := 0; j < 9; j++ {
		if j >= len(i.Accounts) || i.Accounts[j] == nil {
			return fmt.Errorf("accounts[%v] is not set", j)
		}
	}
	if !i.Accounts[1].IsSigner {
		return errors.New("accounts.Manager should be a signer")
	}
	if len(i.Accounts) > 9 && (i.Accounts[9] == nil || !i.Accounts[9].IsSigner) {
		return errors.New("accounts.DepositAuthority should be a signer")
	}
	return nil
}

func (i *Initialize) GetAccounts() []*ag_solanago.AccountMeta {
	return i.Accounts
}

func (i *Initialize) SetAccounts(accounts []*ag_solanago.AccountMeta) error {
	if len(accounts) < 9 {
		return fmt.Errorf("insufficient accounts, Initialize requires at-least 9 accounts not %d", len(accounts))
	}
	i.Accounts = append(ag_solanago.AccountMetaSlice(nil), accounts...)
	return nil
}

func (i *Initialize) SetFee(fee Fee) *Initialize {
	i.Fee = &fee
	return i
}

func (i *Initialize) SetWithdrawalFee(withdrawalFee Fee) *Initialize {
	i.WithdrawalFee = &withdrawalFee
	return i
}

func (i *Initialize) SetDepositFee(depositFee Fee) *Initialize {
	i.DepositFee = &depositFee
	return i
}

func (i *Initialize) SetReferralFee(referralFee uint8) *Initialize {
	i.ReferralFee = &referralFee
	return i
}

func (i *Initialize) SetMaxValidators(maxValidators uint32) *Initialize {
	i.MaxValidators = &maxValidators
	return i
}

func NewInitializeInstruction(
	// Parameters:
	fee Fee,
	withdrawalFee Fee,
	depositFee Fee,
	referralFee uint8,
	maxValidators uint32,
	// Accounts:
	stakePool ag_solanago.PublicKey,
	manager ag_solanago.PublicKey,
	staker ag_solanago.PublicKey,
	withdrawAuthority ag_solanago.PublicKey,
	validatorList ag_solanago.PublicKey,
	reserveStake ag_solanago.PublicKey,
	poolMint ag_solanago.PublicKey,
	managerFeeAccount ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	// depositAuthority is nil for a permissionless pool.
	depositAuthority *ag_solanago.PublicKey,
) *Initialize {
	inst := NewInitializeInstructionBuilder().
		SetFee(fee).
		SetWithdrawalFee(withdrawalFee).
		SetDepositFee(depositFee).
		SetReferralFee(referralFee).
		SetMaxValidators(maxValidators).
		SetStakePoolAccount(stakePool).
		SetManagerAccount(manager).
		SetStakerAccount(staker).
		SetWithdrawAuthorityAccount(withdrawAuthority).
		SetValidatorListAccount(validatorList).
		SetReserveStakeAccount(reserveStake).
		SetPoolMintAccount(poolMint).
		SetManagerFeeAccount(managerFeeAccount).
		SetTokenProgramAccount(tokenProgram)
	if depositAuthority != nil {
		inst.SetDepositAuthorityAccount(*depositAuthority)
	}
	return inst
}
