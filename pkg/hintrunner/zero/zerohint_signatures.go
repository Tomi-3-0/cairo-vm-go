package zero

import (
	"fmt"

	"github.com/NethermindEth/cairo-vm-go/pkg/hintrunner/hinter"
	"github.com/NethermindEth/cairo-vm-go/pkg/hintrunner/utils"
	VM "github.com/NethermindEth/cairo-vm-go/pkg/vm"
	"github.com/NethermindEth/cairo-vm-go/pkg/vm/builtins"
)

func newVerifyECDSASignatureHinter(ecdsaPtr, signature_r, signature_s hinter.ResOperander) hinter.Hinter {
	return &GenericZeroHinter{
		Name: "VerifyECDSASignature",
		Op: func(vm *VM.VirtualMachine, _ *hinter.HintRunnerContext) error {
			//> ecdsa_builtin.add_signature(ids.ecdsa_ptr.address_, (ids.signature_r, ids.signature_s))
			ecdsaPtrAddr, err := hinter.ResolveAsAddress(vm, ecdsaPtr)
			if err != nil {
				return err
			}
			signature_rFelt, err := hinter.ResolveAsFelt(vm, signature_r)
			if err != nil {
				return err
			}
			signature_sFelt, err := hinter.ResolveAsFelt(vm, signature_s)
			if err != nil {
				return err
			}
			ECDSA_segment, ok := vm.Memory.FindSegmentWithBuiltin(builtins.ECDSAName)
			if !ok {
				return fmt.Errorf("ECDSA segment not found")
			}
			ECDSA_builtinRunner := (ECDSA_segment.BuiltinRunner).(*builtins.ECDSA)
			return ECDSA_builtinRunner.AddSignature(ecdsaPtrAddr.Offset, signature_rFelt, signature_sFelt)
		},
	}
}

func createVerifyECDSASignatureHinter(resolver hintReferenceResolver) (hinter.Hinter, error) {
	ecdsaPtr, err := resolver.GetResOperander("ecdsa_ptr")
	if err != nil {
		return nil, err
	}
	signature_r, err := resolver.GetResOperander("signature_r")
	if err != nil {
		return nil, err
	}
	signature_s, err := resolver.GetResOperander("signature_s")
	if err != nil {
		return nil, err
	}
	return newVerifyECDSASignatureHinter(ecdsaPtr, signature_r, signature_s), nil
}

func newImportSecp256R1PHinter() hinter.Hinter {
	return &GenericZeroHinter{
		Name: "Secp256R1",
		Op: func(vm *VM.VirtualMachine, ctx *hinter.HintRunnerContext) error {
			//> from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_P as SECP_P
			SECP256R1_PBig, ok := utils.GetSecp256R1_P()
			if !ok {
				return fmt.Errorf("SECP256R1_P failed.")
			}
			return ctx.ScopeManager.AssignVariable("SECP_P", SECP256R1_PBig)
		},
	}
}

func createImportSecp256R1PHinter() (hinter.Hinter, error) {
	return newImportSecp256R1PHinter(), nil
}
