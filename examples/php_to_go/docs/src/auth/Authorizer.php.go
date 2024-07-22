package auth

type Authorizer struct {
	name string
}

func NewAuthorizer(name string) *Authorizer {
	return &Authorizer{name: name}
}

func (a *Authorizer) Authorize(user string) bool {
	return user == a.name
}
