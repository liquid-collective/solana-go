package stake

// InstructionIDToName maps a stake program instruction id to its name.
// liquid-collective patch: not part of upstream solana-go.
func InstructionIDToName(id uint32) string {
	switch id {
	case Instruction_Initialize:
		return "Initialize"
	case Instruction_Authorize:
		return "Authorize"
	case Instruction_DelegateStake:
		return "DelegateStake"
	case Instruction_Split:
		return "Split"
	case Instruction_Withdraw:
		return "Withdraw"
	case Instruction_Deactivate:
		return "Deactivate"
	case Instruction_SetLockup:
		return "SetLockup"
	case Instruction_Merge:
		return "Merge"
	case Instruction_AuthorizeWithSeed:
		return "AuthorizeWithSeed"
	case Instruction_InitializeChecked:
		return "InitializeChecked"
	case Instruction_AuthorizeChecked:
		return "AuthorizeChecked"
	case Instruction_AuthorizeCheckedWithSeed:
		return "AuthorizeCheckedWithSeed"
	case Instruction_SetLockupChecked:
		return "SetLockupChecked"
	case Instruction_GetMinimumDelegation:
		return "GetMinimumDelegation"
	case Instruction_DeactivateDelinquent:
		return "DeactivateDelinquent"
	case Instruction_Redelegate:
		return "Redelegate"
	case Instruction_MoveStake:
		return "MoveStake"
	case Instruction_MoveLamports:
		return "MoveLamports"
	default:
		return ""
	}
}
