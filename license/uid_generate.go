package license

import (
	"strconv"
)

//Decode是toString的逆方法，NewUID不是toString的逆方法
func Decode(uid string) *UID {
	if len(uid) < 2 {
		return nil
	}
	t, err := strconv.ParseInt(uid[:2], 16, 16)
	if err != nil {
		return nil
	}
	id := uid[2:]
	for i := 0; i < int(UIDTypeMaxLimit); i++ {
		if CheckInput[i](id) {
			if int(t) != i {
				return nil
			}
			return &UID{
				T:  UIDType(i),
				id: id,
			}
		}
	}
	return nil
}
