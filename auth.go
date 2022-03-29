package main

type Authenticator struct{}

func (a *Authenticator) Authenticate(username, password []byte) bool {
	
	return true
}

// ACL returns true if a user has access permissions to read or write on a topic.
// Allow always returns true.
func (a *Authenticator) ACL(user []byte, topic string, write bool) bool {
	return true
}
