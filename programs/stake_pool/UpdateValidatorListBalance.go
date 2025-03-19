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
	"encoding/binary"
	"errors"
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	ag_treeout "github.com/gagliardetto/treeout"

	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
)

type UpdateValidatorListBalance struct {
	Args *UpdateValidatorListBalanceArgs

	// [0] = [] stakePool
	// [1] = [] withdrawAuthority
	// [2] = [WRITE] validatorList
	// [3] = [WRITE] reserveStake
	// [4] = [] clock
	// [5] = [] stakeHistory
	// [6] = [] stakeProgram
	// [7..N] = [] validatorAndTransientStakeAccounts
	Accounts ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func NewUpdateValidatorListBalanceInstruction(
	args UpdateValidatorListBalanceArgs,

	// Accounts:
	stakePool ag_solanago.PublicKey,
	withdrawAuthority ag_solanago.PublicKey,
	validatorList ag_solanago.PublicKey,
	reserveStake ag_solanago.PublicKey,
	clock ag_solanago.PublicKey,
	stakeHistory ag_solanago.PublicKey,
	stakeProgram ag_solanago.PublicKey,
	validatorAndTransientStakeAccounts []ag_solanago.PublicKey,
	transientStakeAccount ag_solanago.PublicKey,
) *UpdateValidatorListBalance {
	return NewUpdateValidatorListBalanceInstructionBuilder().
		SetArgs(args).
		SetStakePool(stakePool).
		SetWithdrawAuthority(withdrawAuthority).
		SetValidatorList(validatorList).
		SetReserveStake(reserveStake).
		SetClock(clock).
		SetStakeHistory(stakeHistory).
		SetStakeProgram(stakeProgram).
		SetValidatorAndTransientAccounts(validatorAndTransientStakeAccounts)
}

func NewUpdateValidatorListBalanceInstructionBuilder() *UpdateValidatorListBalance {
	return &UpdateValidatorListBalance{
		Accounts: make(ag_solanago.AccountMetaSlice, 7),
		Signers:  make(ag_solanago.AccountMetaSlice, 1),
	}
}

func (u *UpdateValidatorListBalance) GetAccounts() []*ag_solanago.AccountMeta {
	return u.Accounts
}

func (u *UpdateValidatorListBalance) SetAccounts(accounts []*ag_solanago.AccountMeta) error {
	u.Accounts = accounts
	return nil
}

func (u *UpdateValidatorListBalance) SetArgs(args UpdateValidatorListBalanceArgs) *UpdateValidatorListBalance {
	u.Args = &args
	return u
}

func (u *UpdateValidatorListBalance) SetStakePool(pool ag_solanago.PublicKey) *UpdateValidatorListBalance {
	u.Accounts[0] = ag_solanago.Meta(pool)
	return u
}

func (u *UpdateValidatorListBalance) SetWithdrawAuthority(withdrawAuthority ag_solanago.PublicKey) *UpdateValidatorListBalance {
	u.Accounts[1] = ag_solanago.Meta(withdrawAuthority)
	return u
}

func (u *UpdateValidatorListBalance) SetValidatorList(validatorList ag_solanago.PublicKey) *UpdateValidatorListBalance {
	u.Accounts[2] = ag_solanago.Meta(validatorList).WRITE()
	return u
}

func (u *UpdateValidatorListBalance) SetReserveStake(reserveStake ag_solanago.PublicKey) *UpdateValidatorListBalance {
	u.Accounts[3] = ag_solanago.Meta(reserveStake).WRITE()
	return u
}

func (u *UpdateValidatorListBalance) SetValidatorAndTransientAccounts(accounts []ag_solanago.PublicKey) *UpdateValidatorListBalance {
	for _, acc := range accounts {
		u.Accounts = append(u.Accounts, ag_solanago.Meta(acc).WRITE())
	}

	return u
}

func (u *UpdateValidatorListBalance) SetValidatorStakeAccount(validatorStakeAccount ag_solanago.PublicKey) *UpdateValidatorListBalance {
	u.Accounts[7] = ag_solanago.Meta(validatorStakeAccount).WRITE()
	return u
}

func (u *UpdateValidatorListBalance) SetClock(clock ag_solanago.PublicKey) *UpdateValidatorListBalance {
	u.Accounts[4] = ag_solanago.Meta(clock)
	return u
}

func (u *UpdateValidatorListBalance) SetStakeHistory(stakeHistory ag_solanago.PublicKey) *UpdateValidatorListBalance {
	u.Accounts[5] = ag_solanago.Meta(stakeHistory)
	return u
}

func (u *UpdateValidatorListBalance) SetStakeProgram(stakeProgram ag_solanago.PublicKey) *UpdateValidatorListBalance {
	u.Accounts[6] = ag_solanago.Meta(stakeProgram)
	return u
}

func (u *UpdateValidatorListBalance) GetStakePool() ag_solanago.PublicKey {
	return u.Accounts[0].PublicKey
}

func (u *UpdateValidatorListBalance) GetWithdrawAuthority() ag_solanago.PublicKey {
	return u.Accounts[1].PublicKey
}

func (u *UpdateValidatorListBalance) GetValidatorList() ag_solanago.PublicKey {
	return u.Accounts[2].PublicKey
}

func (u *UpdateValidatorListBalance) GetReserveStake() ag_solanago.PublicKey {
	return u.Accounts[3].PublicKey
}

func (u *UpdateValidatorListBalance) GetClock() ag_solanago.PublicKey {
	return u.Accounts[4].PublicKey
}

func (u *UpdateValidatorListBalance) GetStakeHistory() ag_solanago.PublicKey {
	return u.Accounts[5].PublicKey
}

func (u *UpdateValidatorListBalance) GetStakeProgram() ag_solanago.PublicKey {
	return u.Accounts[6].PublicKey
}

func (u *UpdateValidatorListBalance) GetValidatorStakeAccount() ag_solanago.PublicKey {
	return u.Accounts[7].PublicKey
}

func (u *UpdateValidatorListBalance) GetTransientStakeAccount() ag_solanago.PublicKey {
	return u.Accounts[8].PublicKey
}

func (u *UpdateValidatorListBalance) ValidateAndBuild() (*Instruction, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u.Build(), nil
}

func (u *UpdateValidatorListBalance) Build() *Instruction {
	return &Instruction{
		BaseVariant: ag_binary.BaseVariant{
			TypeID: ag_binary.TypeIDFromUint8(Instruction_UpdateValidatorListBalance),
			Impl:   u,
		},
	}
}

func (u *UpdateValidatorListBalance) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("UpdateValidatorListBalance")).
				ParentFunc(func(instructionBranch ag_treeout.Branches) {
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						for i, account := range u.Accounts {
							accountsBranch.Child(ag_format.Meta(fmt.Sprintf("[%v]", i), account))
						}
						signersBranch := accountsBranch.Child(fmt.Sprintf("signers[len=%v]", len(u.Signers)))
						for j, signer := range u.Signers {
							signersBranch.Child(ag_format.Meta(fmt.Sprintf("[%v]", j), signer))
						}
					})
				})
		})
}

func (u *UpdateValidatorListBalance) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	// Serialize `Args` param:
	return encoder.Encode(u.Args)
}

func (u *UpdateValidatorListBalance) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return decoder.Decode(&u.Args)
}

func (u *UpdateValidatorListBalance) Validate() error {
	if u.Args == nil {
		return errors.New("Args parameter is not set")
	}

	for i, account := range u.Accounts {
		if account == nil {
			return fmt.Errorf("accounts[%v] is not set", i)
		}
	}

	return nil
}

func (u *UpdateValidatorListBalance) FindTransientStakeAccount(programID, voteAccountAddress, stakePoolAddress ag_solanago.PublicKey, validatorTransitSuffix uint32) (ag_solanago.PublicKey, uint8, error) {
	seedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint32(seedBytes, validatorTransitSuffix)

	seeds := [][]byte{
		[]byte("transient"),
		voteAccountAddress.Bytes(),
		stakePoolAddress.Bytes(),
		seedBytes,
	}

	// Find Program Address (PDA)
	return ag_solanago.FindProgramAddress(seeds, programID)
}

func (u *UpdateValidatorListBalance) FindStakeProgramAddress(programID ag_solanago.PublicKey, voteAccountAddress ag_solanago.PublicKey, stakePoolAddress ag_solanago.PublicKey) (ag_solanago.PublicKey, uint8, error) {
	seeds := [][]byte{
		voteAccountAddress.Bytes(), // Validator Vote Account
		stakePoolAddress.Bytes(),   // Stake Pool Address
	}

	// Find Program Address (PDA)
	return ag_solanago.FindProgramAddress(seeds, programID)
}
