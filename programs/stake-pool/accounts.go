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
	bin "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

type AccountType interface {
	isAccountType()
}

type Uninitialized struct {
}

func (u Uninitialized) isAccountType() {}

type Lockup struct {
	UnixTimestamp uint64
	Epoch         uint64
	Custodian     ag_solanago.PublicKey
}

type FutureEpoch[T any] struct {
	Enum bin.BorshEnum `borsh_enum:"true"`
	None bin.EmptyVariant
	One  T
	Two  T
}

func (f *FutureEpoch[T]) Value() *T {
	switch f.Enum {
	case 0:
		return nil
	case 1:
		return &f.One
	case 2:
		return &f.Two
	default:
		return nil
	}
}

type StakePool struct {
	// Account type, must be `StakePool` currently
	AccountType uint8

	// Manager authority, allows for updating the staker, manager, and fee account
	Manager ag_solanago.PublicKey

	// Staker authority, allows for adding and removing validators, and managing stake distribution
	Staker ag_solanago.PublicKey

	// Stake deposit authority
	StakeDepositAuthority ag_solanago.PublicKey

	// Stake withdrawal authority bump seed for `create_program_address(&[state::StakePool account, "withdrawal"])`
	StakeWithdrawBumpSeed uint8

	// Validator stake list storage account
	ValidatorList ag_solanago.PublicKey

	// Reserve stake account, holds deactivated stake
	ReserveStake ag_solanago.PublicKey

	// Pool Mint
	PoolMint ag_solanago.PublicKey

	// Manager fee account
	ManagerFeeAccount ag_solanago.PublicKey

	// Pool token program id
	TokenProgramID ag_solanago.PublicKey

	// Total stake under management.
	TotalLamports uint64

	// Total supply of pool tokens (should always match the supply in the Pool Mint)
	PoolTokenSupply uint64

	// Last epoch the `total_lamports` field was updated
	LastUpdateEpoch uint64

	// Lockup that all stakes in the pool must have
	Lockup Lockup

	// Fee taken as a proportion of rewards each epoch
	EpochFee Fee

	// Fee for next epoch
	NextEpochFee FutureEpoch[Fee]

	// Preferred deposit validator vote account pubkey
	PreferredDepositValidatorVoteAddress *ag_solanago.PublicKey `bin:"optional"`

	// Preferred withdraw validator vote account pubkey
	PreferredWithdrawValidatorVoteAddress *ag_solanago.PublicKey `bin:"optional"`

	// Fee assessed on stake deposits
	StakeDepositFee Fee

	// Fee assessed on withdrawals
	StakeWithdrawalFee Fee

	// Future stake withdrawal fee, to be set for the following epoch
	NextStakeWithdrawalFee FutureEpoch[Fee]

	// Fees paid out to referrers on referred stake deposits.
	// Expressed as a percentage (0 - 100) of deposit fees.
	StakeReferralFee uint8

	// Toggles whether the `DepositSol` instruction requires a signature from this `sol_deposit_authority`
	SolDepositAuthority *ag_solanago.PublicKey `bin:"optional"`

	// Fee assessed on SOL deposits
	SolDepositFee Fee

	// Fees paid out to referrers on referred SOL deposits.
	// Expressed as a percentage (0 - 100) of SOL deposit fees.
	SolReferralFee uint8

	// Toggles whether the `WithdrawSol` instruction requires a signature from the `deposit_authority`
	SolWithdrawAuthority *ag_solanago.PublicKey `bin:"optional"`

	// Fee assessed on SOL withdrawals
	SolWithdrawalFee Fee

	// Future SOL withdrawal fee, to be set for the following epoch
	NextSolWithdrawalFee FutureEpoch[Fee]

	// Last epoch's total pool tokens, used only for APR estimation
	LastEpochPoolTokenSupply uint64

	// Last epoch's total lamports, used only for APR estimation
	LastEpochTotalLamports uint64
}

func (s StakePool) isAccountType() {}

type ValidatorList struct {
	AccountType   uint8
	MaxValidators uint32
	Validators    []ValidatorStakeInfo
}

func (v ValidatorList) isAccountType() {}
