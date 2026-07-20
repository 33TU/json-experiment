package internal

import "unsafe"

// emptyInterface mirrors the runtime representation of an empty interface.
type emptyInterface struct {
	typ  unsafe.Pointer // concrete type descriptor
	data unsafe.Pointer // interface data word
}

// nonEmptyInterface mirrors the runtime representation of an interface with methods.
type nonEmptyInterface struct {
	tab  unsafe.Pointer // pointer to interfaceTable
	data unsafe.Pointer // interface data word
}

// interfaceTable mirrors the leading words of the runtime interface table.
type interfaceTable struct {
	inter unsafe.Pointer // interface type descriptor
	typ   unsafe.Pointer // concrete type descriptor
}

// InterfaceData returns the data word from v's empty-interface representation.
func InterfaceData(v any) unsafe.Pointer {
	return (*emptyInterface)(unsafe.Pointer(&v)).data
}

// NonEmptyInterfaceValue converts the non-empty interface at ptr to an empty interface.
// It relies on the runtime layouts of non-empty interfaces and interface tables.
func NonEmptyInterfaceValue(ptr unsafe.Pointer) any {
	iface := (*nonEmptyInterface)(ptr)
	if iface.tab == nil {
		return nil
	}

	var value any
	eface := (*emptyInterface)(unsafe.Pointer(&value))
	eface.typ = (*interfaceTable)(iface.tab).typ
	eface.data = iface.data

	return value
}
