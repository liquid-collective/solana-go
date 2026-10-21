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

type AddValidatorToPool struct {
	OptionalSeed *uint32
	// [0] = [WRITE] stakePool
	// [1] = [SIGNER] staker
	// [2] = [WRITE] reserveStake
	// [3] = [] withdrawAuthority
	// [4] = [WRITE] validatorList
	// [5] = [WRITE] validatorStakeAccount
	// [6] = [] voteAccount
	// [7] = [] rent
	// [8] = [] clock
	// [9] = [] stakeHistory
	// [10] = [] stakeConfig
	// [11] = [] systemProgram
	// [12] = [] stakeProgram
	Accounts ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func NewAddValidatorToPoolInstruction(
	// Parameters:
	optionalSeed *uint32,
	// Accounts:
	stakePool ag_solanago.PublicKey,
	staker ag_solanago.PublicKey,
	reserveStake ag_solanago.PublicKey,
	withdrawAuthority ag_solanago.PublicKey,
	validatorList ag_solanago.PublicKey,
	validatorStakeAccount ag_solanago.PublicKey,
	voteAccount ag_solanago.PublicKey,
	rent ag_solanago.PublicKey,
	clock ag_solanago.PublicKey,
	stakeHistory ag_solanago.PublicKey,
	stakeConfig ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey,
	stakeProgram ag_solanago.PublicKey,
) *AddValidatorToPool {
	return NewAddValidatorToPoolBuilder().
		SetOptionalSeed(optionalSeed).
		SetStakePool(stakePool).
		SetStaker(staker).
		SetReserveStake(reserveStake).
		SetWithdrawAuthority(withdrawAuthority).
		SetValidatorList(validatorList).
		SetValidatorStakeAccount(validatorStakeAccount).
		SetVoteAccount(voteAccount).
		SetRent(rent).
		SetClock(clock).
		SetStakeHistory(stakeHistory).
		SetStakeConfig(stakeConfig).
		SetSystemProgram(systemProgram).
		SetStakeProgram(stakeProgram)
}

func NewAddValidatorToPoolBuilder() *AddValidatorToPool {
	return &AddValidatorToPool{
		Accounts: make(ag_solanago.AccountMetaSlice, 13),
	}
}

func (inst *AddValidatorToPool) SetOptionalSeed(seed *uint32) *AddValidatorToPool {
	inst.OptionalSeed = seed
	return inst
}

func (inst *AddValidatorToPool) SetStakePool(pool ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[0] = ag_solanago.Meta(pool).WRITE()
	return inst
}

func (inst *AddValidatorToPool) SetStaker(staker ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[1] = ag_solanago.Meta(staker).SIGNER()
	return inst
}

func (inst *AddValidatorToPool) SetReserveStake(reserveStake ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[2] = ag_solanago.Meta(reserveStake).WRITE()
	return inst
}

func (inst *AddValidatorToPool) SetWithdrawAuthority(withdrawAuthority ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[3] = ag_solanago.Meta(withdrawAuthority)
	return inst
}

func (inst *AddValidatorToPool) SetValidatorList(validatorList ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[4] = ag_solanago.Meta(validatorList).WRITE()
	return inst
}

func (inst *AddValidatorToPool) SetValidatorStakeAccount(validatorStakeAccount ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[5] = ag_solanago.Meta(validatorStakeAccount).WRITE()
	return inst
}

func (inst *AddValidatorToPool) SetVoteAccount(voteAccount ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[6] = ag_solanago.Meta(voteAccount)
	return inst
}

func (inst *AddValidatorToPool) SetRent(rent ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[7] = ag_solanago.Meta(rent)
	return inst
}

func (inst *AddValidatorToPool) SetClock(clock ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[8] = ag_solanago.Meta(clock)
	return inst
}

func (inst *AddValidatorToPool) SetStakeHistory(stakeHistory ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[9] = ag_solanago.Meta(stakeHistory)
	return inst
}

func (inst *AddValidatorToPool) SetStakeConfig(stakeConfig ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[10] = ag_solanago.Meta(stakeConfig)
	return inst
}

func (inst *AddValidatorToPool) SetSystemProgram(systemProgram ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[11] = ag_solanago.Meta(systemProgram)
	return inst
}

func (inst *AddValidatorToPool) SetStakeProgram(stakeProgram ag_solanago.PublicKey) *AddValidatorToPool {
	inst.Accounts[12] = ag_solanago.Meta(stakeProgram)
	return inst
}

func (inst *AddValidatorToPool) GetOptionalSeed() *uint32 {
	return inst.OptionalSeed
}

func (inst *AddValidatorToPool) GetStakePool() ag_solanago.PublicKey {
	return inst.Accounts[0].PublicKey
}

func (inst *AddValidatorToPool) GetStaker() ag_solanago.PublicKey {
	return inst.Accounts[1].PublicKey
}

func (inst *AddValidatorToPool) GetReserveStake() ag_solanago.PublicKey {
	return inst.Accounts[2].PublicKey
}

func (inst *AddValidatorToPool) GetWithdrawAuthority() ag_solanago.PublicKey {
	return inst.Accounts[3].PublicKey
}

func (inst *AddValidatorToPool) GetValidatorList() ag_solanago.PublicKey {
	return inst.Accounts[4].PublicKey
}

func (inst *AddValidatorToPool) GetValidatorStakeAccount() ag_solanago.PublicKey {
	return inst.Accounts[5].PublicKey
}

func (inst *AddValidatorToPool) GetVoteAccount() ag_solanago.PublicKey {
	return inst.Accounts[6].PublicKey
}

func (inst *AddValidatorToPool) GetRent() ag_solanago.PublicKey {
	return inst.Accounts[7].PublicKey
}

func (inst *AddValidatorToPool) GetClock() ag_solanago.PublicKey {
	return inst.Accounts[8].PublicKey
}

func (inst *AddValidatorToPool) GetStakeHistory() ag_solanago.PublicKey {
	return inst.Accounts[9].PublicKey
}

func (inst *AddValidatorToPool) GetStakeConfig() ag_solanago.PublicKey {
	return inst.Accounts[10].PublicKey
}

func (inst *AddValidatorToPool) GetSystemProgram() ag_solanago.PublicKey {
	return inst.Accounts[11].PublicKey
}

func (inst *AddValidatorToPool) GetStakeProgram() ag_solanago.PublicKey {
	return inst.Accounts[12].PublicKey
}

func (inst *AddValidatorToPool) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AddValidatorToPool) Build() *Instruction {
	return &Instruction{
		BaseVariant: ag_binary.BaseVariant{
			TypeID: ag_binary.TypeIDFromUint8(Instruction_AddValidatorToPool),
			Impl:   inst,
		},
	}
}

func (inst *AddValidatorToPool) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("AddValidatorToPool")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {
					//
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						if inst.OptionalSeed != nil {
							paramsBranch.Child(ag_format.Param("OptionalSeed", *inst.OptionalSeed))
						}
					})

					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("             stake_pool", inst.Accounts.Get(0)))
						accountsBranch.Child(ag_format.Meta("                 staker", inst.Accounts.Get(1)))
						accountsBranch.Child(ag_format.Meta("          reserve_stake", inst.Accounts.Get(2)))
						accountsBranch.Child(ag_format.Meta("     withdraw_authority", inst.Accounts.Get(3)))
						accountsBranch.Child(ag_format.Meta("         validator_list", inst.Accounts.Get(4)))
						accountsBranch.Child(ag_format.Meta("validator_stake_account", inst.Accounts.Get(5)))
						accountsBranch.Child(ag_format.Meta("           vote_account", inst.Accounts.Get(6)))
						accountsBranch.Child(ag_format.Meta("                   rent", inst.Accounts.Get(7)))
						accountsBranch.Child(ag_format.Meta("                  clock", inst.Accounts.Get(8)))
						accountsBranch.Child(ag_format.Meta("          stake_history", inst.Accounts.Get(9)))
						accountsBranch.Child(ag_format.Meta("           stake_config", inst.Accounts.Get(10)))
						accountsBranch.Child(ag_format.Meta("         system_program", inst.Accounts.Get(11)))
						accountsBranch.Child(ag_format.Meta("          stake_program", inst.Accounts.Get(12)))

					})
				})
		})
}

func (inst *AddValidatorToPool) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	var optionalSeed uint32
	if inst.OptionalSeed != nil {
		optionalSeed = *inst.OptionalSeed
	}
	if err := encoder.WriteUint32(optionalSeed, ag_binary.LE); err != nil {
		return err
	}
	return nil
}

func (inst *AddValidatorToPool) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	optionalSeed, err := decoder.ReadUint32(ag_binary.LE)
	if err != nil {
		return err
	}
	inst.OptionalSeed = &optionalSeed
	return nil
}

func (inst *AddValidatorToPool) Validate() error {
	for i := 0; i < 13; i++ {
		if i >= len(inst.Accounts) || inst.Accounts[i] == nil {
			return fmt.Errorf("accounts[%v] is not set", i)
		}
	}
	for _, account := range inst.Accounts[13:] {
		if account == nil {
			return errors.New("an optional account is not set")
		}
	}
	if !inst.Accounts[1].IsSigner {
		return errors.New("accounts.Staker should be a signer")
	}
	return nil
}

func (inst *AddValidatorToPool) GetAccounts() []*ag_solanago.AccountMeta {
	return inst.Accounts
}

func (inst *AddValidatorToPool) SetAccounts(accounts []*ag_solanago.AccountMeta) error {
	if len(accounts) < 13 {
		return fmt.Errorf("insufficient accounts, AddValidatorToPool requires at-least 13 accounts not %d", len(accounts))
	}
	inst.Accounts = append(ag_solanago.AccountMetaSlice(nil), accounts...)
	return nil
}
