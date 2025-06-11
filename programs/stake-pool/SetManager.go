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

type SetManager struct {
	// [0] = [WRITE] stakePool
	// [1] = [SIGNER] manager
	// [2] = [SIGNER] newManager
	// [3] = [] newManagerFeeAccount
	ag_solanago.AccountMetaSlice `bin:"-"`
}

func NewSetManagerInstruction(
// Accounts:
	stakePool ag_solanago.PublicKey,
	manager ag_solanago.PublicKey,
	newManager ag_solanago.PublicKey,
	newManagerFeeAccount ag_solanago.PublicKey,
) *SetManager {
	return NewSetManagerInstructionBuilder().
		SetStakePool(stakePool).
		SetManager(manager).
		SetNewManager(newManager).
		SetNewManagerFeeAccount(newManagerFeeAccount)
}

func NewSetManagerInstructionBuilder() *SetManager {
	return &SetManager{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
}

func (inst *SetManager) SetStakePool(pool ag_solanago.PublicKey) *SetManager {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(pool).WRITE()
	return inst
}

func (inst *SetManager) SetManager(manager ag_solanago.PublicKey) *SetManager {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(manager).SIGNER()
	return inst
}

func (inst *SetManager) SetNewManager(newManager ag_solanago.PublicKey) *SetManager {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(newManager).SIGNER()
	return inst
}

func (inst *SetManager) SetNewManagerFeeAccount(newManagerFeeAccount ag_solanago.PublicKey) *SetManager {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(newManagerFeeAccount)
	return inst
}

func (inst *SetManager) GetStakePool() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

func (inst *SetManager) GetManager() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst *SetManager) GetNewManager() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst *SetManager) GetNewManagerFeeAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[3]
}

func (inst *SetManager) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetManager) Build() *Instruction {
	return &Instruction{
		BaseVariant: ag_binary.BaseVariant{
			TypeID: ag_binary.TypeIDFromUint8(Instruction_SetManager),
			Impl:   inst,
		},
	}
}

func (inst *SetManager) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("SetManager")).
				ParentFunc(func(instructionBranch ag_treeout.Branches) {
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("StakePool", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("Manager", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("NewManager", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("NewManagerFeeAccount", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (inst *SetManager) Validate() error {
	for accIndex, acc := range inst.AccountMetaSlice {
		if acc == nil {
			return fmt.Errorf("ins.AccountMetaSlice[%v] is not set", accIndex)
		}
	}
	return nil
}
