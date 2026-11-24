// Copyright The KCL Authors. All rights reserved.

package libkcl

import (
	"debug/elf"
	"errors"
	"strings"
)

// exportsSymbol reports whether path is a shared object that defines name.
// Versioned names such as kcl_service_new@KCL match.
func exportsSymbol(path, name string) (bool, error) {
	f, err := elf.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if ok, err := symbolTableHas(f.DynamicSymbols, name); ok || (err != nil && !errors.Is(err, elf.ErrNoSymbols)) {
		return ok, err
	}
	ok, err := symbolTableHas(f.Symbols, name)
	if errors.Is(err, elf.ErrNoSymbols) {
		return false, nil
	}
	return ok, err
}

func symbolTableHas(load func() ([]elf.Symbol, error), name string) (bool, error) {
	syms, err := load()
	if err != nil {
		return false, err
	}
	for _, sym := range syms {
		if sym.Section == elf.SHN_UNDEF {
			continue
		}
		if sym.Name == name || strings.HasPrefix(sym.Name, name+"@") {
			return true, nil
		}
	}
	return false, nil
}
