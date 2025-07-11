package flag

import (
	"fmt"
	"strconv"
)

// Tristate is a flag compatible with flags and pflags that
// keeps track of whether it had a value supplied or not.
type Tristate int

const (
	Unset Tristate = iota // 0
	True
	False
)

func (f *Tristate) Default(value bool) {
	*f = triFromBool(value)
}

func (f Tristate) String() string {
	b := boolFromTri(f)
	return fmt.Sprintf("%t", b)
}

func (f Tristate) Value() bool {
	b := boolFromTri(f)
	return b
}

func (f *Tristate) Set(value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	*f = triFromBool(boolVal)
	return nil
}

func (f Tristate) Provided() bool {
	return f != Unset
}

func (f *Tristate) Type() string {
	return "tristate"
}

func boolFromTri(t Tristate) bool {
	if t == True {
		return true
	} else {
		return false
	}
}

func triFromBool(b bool) Tristate {
	if b {
		return True
	} else {
		return False
	}
}
