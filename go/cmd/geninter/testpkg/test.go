package testpkg

import "context"

type User struct {
	name string
	age  int
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) Age() int {
	return u.age
}

func (u *User) SetAge(age int) {
	u.age = age
}

func (u *User) Info() (name string, age int) {
	return u.name, u.age
}

func (u *User) SetInfo(name string, age int) {
	u.name = name
	u.age = age
}

func (u *User) IsAdult(ctx context.Context) bool {
	return u.age >= 18
}
