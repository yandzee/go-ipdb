package types

import (
	"encoding/binary"
	"fmt"
)

type Uint128 struct {
	Hi uint64
	Lo uint64
}

func U128FromBEBytes(bytes [16]byte) Uint128 {
	return Uint128{
		Hi: binary.BigEndian.Uint64(bytes[:8]),
		Lo: binary.BigEndian.Uint64(bytes[8:]),
	}
}

func (u *Uint128) Hex() string {
	return fmt.Sprintf("0x%016x%016x", u.Hi, u.Lo)
}

func (u Uint128) Compare(other Uint128) int {
	if u.Hi < other.Hi {
		return -1
	}

	if u.Hi > other.Hi {
		return 1
	}

	if u.Lo < other.Lo {
		return -1
	}

	if u.Lo > other.Lo {
		return 1
	}

	return 0
}
