package testpkg

import "fmt"

func (u *User) Says(str string) string {
	return fmt.Sprintf("%s says: %s", u.name, str)
}
