package embeddedmls

/*
#cgo LDFLAGS: -L./lib -lbindings_c
#include <stdbool.h>
#include <stdlib.h>
#include "lib/libxmtpmls.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

func ValidateKeyPackage(data []byte) error {
	result := C.validate_inbox_id_key_package_ffi(
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.ulong(len(data)),
	)
	defer C.free_c_string(result.message)
	if !result.ok {
		return errors.New(C.GoString(result.message))
	}

	return nil
}
