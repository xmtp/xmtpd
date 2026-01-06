package config

import (
	"strconv"
	"strings"
)

// Uint32Slice parses a comma-separated list of uint32 values.
//
// It treats an empty string ("") or whitespace-only input as an empty slice,
// which allows environment variables like XMTPD_PAYER_NODE_SELECTOR_PREFERRED_NODES=""
// to be provided without causing a parse error.
//
// Examples:
//
//	""            -> nil
//	"   "         -> nil
//	"1,2,3"       -> []uint32{1,2,3}
//	"1, 2, 3"     -> []uint32{1,2,3}
//	"1,2,"        -> []uint32{1,2}
//	",1,,2,3,"    -> []uint32{1,2,3}
type Uint32Slice []uint32

// UnmarshalFlag is used by github.com/jessevdk/go-flags to parse flag/env values.
func (s *Uint32Slice) UnmarshalFlag(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		*s = nil
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]uint32, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		u, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			return err
		}
		out = append(out, uint32(u))
	}

	*s = out
	return nil
}

func (s Uint32Slice) String() string {
	if len(s) == 0 {
		return ""
	}
	parts := make([]string, 0, len(s))
	for _, n := range s {
		parts = append(parts, strconv.FormatUint(uint64(n), 10))
	}
	return strings.Join(parts, ",")
}
