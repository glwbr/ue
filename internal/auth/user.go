package auth

type User struct {
	FirstName string
	LastName  string
	Email     string
}

func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Email
	}
	return u.FirstName + " " + u.LastName
}
