package stakepool

import (
	"errors"
	"fmt"
	"math"

	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

// Instruction data is Borsh: little-endian integers, strings with a 4-byte length prefix and
// options with a one-byte tag. These helpers are shared by the instruction codecs.

func writeBorshString(encoder *ag_binary.Encoder, s string) error {
	if uint64(len(s)) > math.MaxUint32 {
		return errors.New("string is too long")
	}
	if err := encoder.WriteUint32(uint32(len(s)), ag_binary.LE); err != nil {
		return err
	}
	return encoder.WriteBytes([]byte(s), false)
}

func readBorshString(decoder *ag_binary.Decoder) (string, error) {
	n, err := decoder.ReadUint32(ag_binary.LE)
	if err != nil {
		return "", err
	}
	if uint64(n) > uint64(decoder.Remaining()) {
		return "", fmt.Errorf("string length %d exceeds the %d remaining bytes", n, decoder.Remaining())
	}
	b, err := decoder.ReadNBytes(int(n))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func writeOptionPublicKey(encoder *ag_binary.Encoder, pk *ag_solanago.PublicKey) error {
	if pk == nil {
		return encoder.WriteUint8(0)
	}
	if err := encoder.WriteUint8(1); err != nil {
		return err
	}
	return encoder.WriteBytes(pk[:], false)
}

func readOptionPublicKey(decoder *ag_binary.Decoder) (*ag_solanago.PublicKey, error) {
	tag, err := decoder.ReadUint8()
	if err != nil {
		return nil, err
	}
	switch tag {
	case 0:
		return nil, nil
	case 1:
		b, err := decoder.ReadNBytes(ag_solanago.PublicKeyLength)
		if err != nil {
			return nil, err
		}
		var pk ag_solanago.PublicKey
		copy(pk[:], b)
		return &pk, nil
	default:
		return nil, fmt.Errorf("invalid option tag %d", tag)
	}
}

func writeFee(encoder *ag_binary.Encoder, f Fee) error {
	if err := encoder.WriteUint64(f.Denominator, ag_binary.LE); err != nil {
		return err
	}
	return encoder.WriteUint64(f.Numerator, ag_binary.LE)
}

func readFee(decoder *ag_binary.Decoder) (Fee, error) {
	denominator, err := decoder.ReadUint64(ag_binary.LE)
	if err != nil {
		return Fee{}, err
	}
	numerator, err := decoder.ReadUint64(ag_binary.LE)
	if err != nil {
		return Fee{}, err
	}
	return Fee{Denominator: denominator, Numerator: numerator}, nil
}

func writeAdditionalValidatorStakeArgs(encoder *ag_binary.Encoder, a *AdditionalValidatorStakeArgs) error {
	for _, v := range []uint64{a.Lamports, a.TransientStakeSeed, a.EphemeralStakeSeed} {
		if err := encoder.WriteUint64(v, ag_binary.LE); err != nil {
			return err
		}
	}
	return nil
}

func readAdditionalValidatorStakeArgs(decoder *ag_binary.Decoder) (*AdditionalValidatorStakeArgs, error) {
	var vals [3]uint64
	for i := range vals {
		v, err := decoder.ReadUint64(ag_binary.LE)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return &AdditionalValidatorStakeArgs{Lamports: vals[0], TransientStakeSeed: vals[1], EphemeralStakeSeed: vals[2]}, nil
}

// FeeType is an enum with payloads; the tags follow the declaration order of the SPL enum.
const (
	feeTypeSolReferral uint8 = iota
	feeTypeStakeReferral
	feeTypeEpoch
	feeTypeStakeWithdrawal
	feeTypeSolDeposit
	feeTypeStakeDeposit
	feeTypeSolWithdrawal
)

func writeFeeType(encoder *ag_binary.Encoder, t FeeType) error {
	writeFeeVariant := func(tag uint8, f Fee) error {
		if err := encoder.WriteUint8(tag); err != nil {
			return err
		}
		return writeFee(encoder, f)
	}
	writeReferral := func(tag uint8, percent uint8) error {
		if err := encoder.WriteUint8(tag); err != nil {
			return err
		}
		return encoder.WriteUint8(percent)
	}

	switch v := t.(type) {
	case SolReferral:
		return writeReferral(feeTypeSolReferral, v.Fee)
	case *SolReferral:
		return writeReferral(feeTypeSolReferral, v.Fee)
	case StakeReferral:
		return writeReferral(feeTypeStakeReferral, v.Fee)
	case *StakeReferral:
		return writeReferral(feeTypeStakeReferral, v.Fee)
	case Epoch:
		return writeFeeVariant(feeTypeEpoch, v.Fee)
	case *Epoch:
		return writeFeeVariant(feeTypeEpoch, v.Fee)
	case StakeWithdrawal:
		return writeFeeVariant(feeTypeStakeWithdrawal, v.Fee)
	case *StakeWithdrawal:
		return writeFeeVariant(feeTypeStakeWithdrawal, v.Fee)
	case SolDeposit:
		return writeFeeVariant(feeTypeSolDeposit, v.Fee)
	case *SolDeposit:
		return writeFeeVariant(feeTypeSolDeposit, v.Fee)
	case StakeDeposit:
		return writeFeeVariant(feeTypeStakeDeposit, v.Fee)
	case *StakeDeposit:
		return writeFeeVariant(feeTypeStakeDeposit, v.Fee)
	case SolWithdrawal:
		return writeFeeVariant(feeTypeSolWithdrawal, v.Fee)
	case *SolWithdrawal:
		return writeFeeVariant(feeTypeSolWithdrawal, v.Fee)
	default:
		return fmt.Errorf("unsupported fee type %T", t)
	}
}

func readFeeType(decoder *ag_binary.Decoder) (FeeType, error) {
	tag, err := decoder.ReadUint8()
	if err != nil {
		return nil, err
	}
	switch tag {
	case feeTypeSolReferral, feeTypeStakeReferral:
		percent, err := decoder.ReadUint8()
		if err != nil {
			return nil, err
		}
		if tag == feeTypeSolReferral {
			return SolReferral{Fee: percent}, nil
		}
		return StakeReferral{Fee: percent}, nil
	case feeTypeEpoch, feeTypeStakeWithdrawal, feeTypeSolDeposit, feeTypeStakeDeposit, feeTypeSolWithdrawal:
		fee, err := readFee(decoder)
		if err != nil {
			return nil, err
		}
		switch tag {
		case feeTypeEpoch:
			return Epoch{Fee: fee}, nil
		case feeTypeStakeWithdrawal:
			return StakeWithdrawal{Fee: fee}, nil
		case feeTypeSolDeposit:
			return SolDeposit{Fee: fee}, nil
		case feeTypeStakeDeposit:
			return StakeDeposit{Fee: fee}, nil
		default:
			return SolWithdrawal{Fee: fee}, nil
		}
	default:
		return nil, fmt.Errorf("invalid fee type tag %d", tag)
	}
}

// setAccount puts meta at idx: it replaces an existing account and appends otherwise, so setting an
// optional or variable account on a populated or decoded instruction never drops the accounts around it.
func setAccount(accounts ag_solanago.AccountMetaSlice, idx int, meta *ag_solanago.AccountMeta) ag_solanago.AccountMetaSlice {
	for len(accounts) <= idx {
		accounts = append(accounts, nil)
	}
	accounts[idx] = meta
	return accounts
}
