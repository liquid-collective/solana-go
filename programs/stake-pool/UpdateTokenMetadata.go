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

type UpdateTokenMetadata struct {
	Name   *string
	Symbol *string
	URI    *string
	// [0] = [WRITE] stakePool
	// [1] = [SIGNER] manager
	// [2] = [WRITE] stakePoolWithdrawAuthority
	// [3] = [WRITE] tokenMetadata
	// [4] = [] mplTokenMetadata
	Accounts ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func NewUpdateTokenMetadataInstruction(
	name string,
	symbol string,
	uri string,
	stakePool ag_solanago.PublicKey,
	manager ag_solanago.PublicKey,
	stakePoolWithdrawAuthority ag_solanago.PublicKey,
	tokenMetadata ag_solanago.PublicKey,
	mplTokenMetadata ag_solanago.PublicKey,
) *UpdateTokenMetadata {
	return NewUpdateTokenMetadataInstructionBuilder().
		SetName(name).
		SetSymbol(symbol).
		SetURI(uri).
		SetStakePool(stakePool).
		SetManager(manager).
		SetStakePoolWithdrawAuthority(stakePoolWithdrawAuthority).
		SetTokenMetadata(tokenMetadata).
		SetMplTokenMetadata(mplTokenMetadata)
}

func NewUpdateTokenMetadataInstructionBuilder() *UpdateTokenMetadata {
	return &UpdateTokenMetadata{
		Accounts: make(ag_solanago.AccountMetaSlice, 5),
	}
}

func (inst *UpdateTokenMetadata) SetName(name string) *UpdateTokenMetadata {
	inst.Name = &name
	return inst
}

func (inst *UpdateTokenMetadata) SetSymbol(symbol string) *UpdateTokenMetadata {
	inst.Symbol = &symbol
	return inst
}

func (inst *UpdateTokenMetadata) SetURI(uri string) *UpdateTokenMetadata {
	inst.URI = &uri
	return inst
}

func (inst *UpdateTokenMetadata) SetStakePool(stakePool ag_solanago.PublicKey) *UpdateTokenMetadata {
	inst.Accounts[0] = ag_solanago.Meta(stakePool)
	return inst
}

func (inst *UpdateTokenMetadata) SetManager(manager ag_solanago.PublicKey) *UpdateTokenMetadata {
	inst.Accounts[1] = ag_solanago.Meta(manager).SIGNER()
	return inst
}

func (inst *UpdateTokenMetadata) SetStakePoolWithdrawAuthority(stakePoolWithdrawAuthority ag_solanago.PublicKey) *UpdateTokenMetadata {
	inst.Accounts[2] = ag_solanago.Meta(stakePoolWithdrawAuthority)
	return inst
}

func (inst *UpdateTokenMetadata) SetTokenMetadata(tokenMetadata ag_solanago.PublicKey) *UpdateTokenMetadata {
	inst.Accounts[3] = ag_solanago.Meta(tokenMetadata).WRITE()
	return inst
}

func (inst *UpdateTokenMetadata) SetMplTokenMetadata(mplTokenMetadata ag_solanago.PublicKey) *UpdateTokenMetadata {
	inst.Accounts[4] = ag_solanago.Meta(mplTokenMetadata)
	return inst
}

func (inst *UpdateTokenMetadata) GetName() *string {
	return inst.Name
}

func (inst *UpdateTokenMetadata) GetSymbol() *string {
	return inst.Symbol
}

func (inst *UpdateTokenMetadata) GetURI() *string {
	return inst.URI
}

func (inst *UpdateTokenMetadata) GetStakePool() ag_solanago.PublicKey {
	return inst.Accounts[0].PublicKey
}

func (inst *UpdateTokenMetadata) GetManager() ag_solanago.PublicKey {
	return inst.Accounts[1].PublicKey
}

func (inst *UpdateTokenMetadata) GetStakePoolWithdrawAuthority() ag_solanago.PublicKey {
	return inst.Accounts[2].PublicKey
}

func (inst *UpdateTokenMetadata) GetTokenMetadata() ag_solanago.PublicKey {
	return inst.Accounts[3].PublicKey
}

func (inst *UpdateTokenMetadata) GetMplTokenMetadata() ag_solanago.PublicKey {
	return inst.Accounts[4].PublicKey
}

func (inst *UpdateTokenMetadata) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *UpdateTokenMetadata) Build() *Instruction {
	return &Instruction{
		BaseVariant: ag_binary.BaseVariant{
			TypeID: ag_binary.TypeIDFromUint8(Instruction_UpdateTokenMetadata),
			Impl:   inst,
		},
	}
}

func (inst *UpdateTokenMetadata) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("UpdateTokenMetadata")).
				ParentFunc(func(instructionBranch ag_treeout.Branches) {
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						if inst.Name != nil {
							paramsBranch.Child(ag_format.Param("Name", *inst.Name))
						}
						if inst.Symbol != nil {
							paramsBranch.Child(ag_format.Param("Symbol", *inst.Symbol))
						}
						if inst.URI != nil {
							paramsBranch.Child(ag_format.Param("URI", *inst.URI))
						}
					})
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						for i, account := range inst.Accounts {
							accountsBranch.Child(ag_format.Meta(fmt.Sprintf("[%v]", i), account))
						}
					})
				})
		})
}

func (inst *UpdateTokenMetadata) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	if inst.Name == nil {
		return errors.New("name is not set")
	}
	if err := writeBorshString(encoder, *inst.Name); err != nil {
		return err
	}
	if inst.Symbol == nil {
		return errors.New("symbol is not set")
	}
	if err := writeBorshString(encoder, *inst.Symbol); err != nil {
		return err
	}
	if inst.URI == nil {
		return errors.New("uRI is not set")
	}
	if err := writeBorshString(encoder, *inst.URI); err != nil {
		return err
	}
	return nil
}

func (inst *UpdateTokenMetadata) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	name, err := readBorshString(decoder)
	if err != nil {
		return err
	}
	inst.Name = &name
	symbol, err := readBorshString(decoder)
	if err != nil {
		return err
	}
	inst.Symbol = &symbol
	uRI, err := readBorshString(decoder)
	if err != nil {
		return err
	}
	inst.URI = &uRI
	return nil
}

func (inst *UpdateTokenMetadata) Validate() error {
	if inst.Name == nil {
		return errors.New("name is not set")
	}
	if inst.Symbol == nil {
		return errors.New("symbol is not set")
	}
	if inst.URI == nil {
		return errors.New("uRI is not set")
	}
	for i := 0; i < 5; i++ {
		if i >= len(inst.Accounts) || inst.Accounts[i] == nil {
			return fmt.Errorf("accounts[%v] is not set", i)
		}
	}
	for _, account := range inst.Accounts[5:] {
		if account == nil {
			return errors.New("an optional account is not set")
		}
	}
	if !inst.Accounts[1].IsSigner {
		return errors.New("accounts.Manager should be a signer")
	}
	return nil
}

func (inst *UpdateTokenMetadata) GetAccounts() []*ag_solanago.AccountMeta {
	return inst.Accounts
}

func (inst *UpdateTokenMetadata) SetAccounts(accounts []*ag_solanago.AccountMeta) error {
	if len(accounts) < 5 {
		return fmt.Errorf("insufficient accounts, UpdateTokenMetadata requires at-least 5 accounts not %d", len(accounts))
	}
	inst.Accounts = append(ag_solanago.AccountMetaSlice(nil), accounts...)
	return nil
}
