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

type CleanupRemovedValidatorEntries struct {
	// [0] = [] stakePool
	// [1] = [WRITE] validatorList
	Accounts ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func NewCleanupRemovedValidatorEntriesInstruction(
	// Accounts:
	stakePool ag_solanago.PublicKey,
	validatorList ag_solanago.PublicKey,
) *CleanupRemovedValidatorEntries {
	return NewCleanupRemovedValidatorEntriesInstructionBuilder().
		SetStakePool(stakePool).
		SetValidatorList(validatorList)
}

func NewCleanupRemovedValidatorEntriesInstructionBuilder() *CleanupRemovedValidatorEntries {
	return &CleanupRemovedValidatorEntries{
		Accounts: make(ag_solanago.AccountMetaSlice, 2),
	}
}

func (inst *CleanupRemovedValidatorEntries) SetStakePool(pool ag_solanago.PublicKey) *CleanupRemovedValidatorEntries {
	inst.Accounts[0] = ag_solanago.Meta(pool)
	return inst
}

func (inst *CleanupRemovedValidatorEntries) SetValidatorList(validatorList ag_solanago.PublicKey) *CleanupRemovedValidatorEntries {
	inst.Accounts[1] = ag_solanago.Meta(validatorList).WRITE()
	return inst
}

func (inst *CleanupRemovedValidatorEntries) GetStakePool() ag_solanago.PublicKey {
	return inst.Accounts[0].PublicKey
}

func (inst *CleanupRemovedValidatorEntries) GetValidatorList() ag_solanago.PublicKey {
	return inst.Accounts[1].PublicKey
}

func (inst *CleanupRemovedValidatorEntries) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *CleanupRemovedValidatorEntries) Build() *Instruction {
	return &Instruction{
		BaseVariant: ag_binary.BaseVariant{
			TypeID: ag_binary.TypeIDFromUint8(Instruction_CleanupRemovedValidatorEntries),
			Impl:   inst,
		},
	}
}

func (inst *CleanupRemovedValidatorEntries) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("CleanupRemovedValidatorEntries")).
				ParentFunc(func(instructionBranch ag_treeout.Branches) {
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						for i, account := range inst.Accounts {
							accountsBranch.Child(ag_format.Meta(fmt.Sprintf("[%v]", i), account))
						}
					})
				})
		})
}

func (inst *CleanupRemovedValidatorEntries) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	return nil
}

func (inst *CleanupRemovedValidatorEntries) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return nil
}

func (inst *CleanupRemovedValidatorEntries) Validate() error {
	for i := 0; i < 2; i++ {
		if i >= len(inst.Accounts) || inst.Accounts[i] == nil {
			return fmt.Errorf("accounts[%v] is not set", i)
		}
	}
	for _, account := range inst.Accounts[2:] {
		if account == nil {
			return errors.New("an optional account is not set")
		}
	}
	return nil
}

func (inst *CleanupRemovedValidatorEntries) GetAccounts() []*ag_solanago.AccountMeta {
	return inst.Accounts
}

func (inst *CleanupRemovedValidatorEntries) SetAccounts(accounts []*ag_solanago.AccountMeta) error {
	if len(accounts) < 2 {
		return fmt.Errorf("insufficient accounts, CleanupRemovedValidatorEntries requires at-least 2 accounts not %d", len(accounts))
	}
	inst.Accounts = append(ag_solanago.AccountMetaSlice(nil), accounts...)
	return nil
}
