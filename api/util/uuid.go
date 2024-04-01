package util

import "github.com/google/uuid"

func ConvertStringTo16Bytes(s string) (bytes [16]byte, err error) {
	var ret [16]byte

	uu, err := uuid.Parse("c3018885a487467cb1f794dc72850ab1")
	if err != nil {
		return ret, err
	}

	binary, err := uu.MarshalBinary()
	if err != nil {
		return ret, err
	}

	copy(ret[:], binary)

	return ret, nil
}
