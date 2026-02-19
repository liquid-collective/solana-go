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
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	ag_treeout "github.com/gagliardetto/treeout"

	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
)

type SetStaker struct {
	// [0] = [WRITE] stakePool
	// [1] = [SIGNER] currentStaker
	// [2] = [] newStaker
	ag_solanago.AccountMetaSlice `bin:"-"`
}

func NewSetStakerInstruction(
// Accounts:
	stakePool ag_solanago.PublicKey,
	signer ag_solanago.PublicKey,
	newStaker ag_solanago.PublicKey,
) *SetStaker {
	return NewSetStakerInstructionBuilder().
		SetStakePool(stakePool).
		SetCurrentStaker(signer).
		SetNewStaker(newStaker)
}

func NewSetStakerInstructionBuilder() *SetStaker {
	return &SetStaker{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
}

func (inst *SetStaker) SetStakePool(pool ag_solanago.PublicKey) *SetStaker {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(pool).WRITE()
	return inst
}

func (inst *SetStaker) SetCurrentStaker(staker ag_solanago.PublicKey) *SetStaker {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(staker).SIGNER()
	return inst
}

func (inst *SetStaker) SetNewStaker(newStaker ag_solanago.PublicKey) *SetStaker {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(newStaker)
	return inst
}

func (inst *SetStaker) GetStakePool() ag_solanago.PublicKey {
	return inst.AccountMetaSlice[0].PublicKey
}

func (inst *SetStaker) GetCurrentStaker() ag_solanago.PublicKey {
	return inst.AccountMetaSlice[1].PublicKey
}

func (inst *SetStaker) GetNewStaker() ag_solanago.PublicKey {
	return inst.AccountMetaSlice[2].PublicKey
}

func (inst *SetStaker) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetStaker) Build() *Instruction {
	return &Instruction{
		BaseVariant: ag_binary.BaseVariant{
			TypeID: ag_binary.TypeIDFromUint8(Instruction_SetStaker),
			Impl:   inst,
		},
	}
}

func (inst *SetStaker) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("SetStaker")).
				ParentFunc(func(instructionBranch ag_treeout.Branches) {
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("StakePool", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("CurrentStaker", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("NewStaker", inst.AccountMetaSlice.Get(2)))
					})
				})
		})
}

func (inst *SetStaker) Validate() error {
	for i, account := range inst.AccountMetaSlice {
		if account == nil {
			return fmt.Errorf("ins.AccountMetaSlice[%v] is not set", i)
		}
	}
	return nil
}
