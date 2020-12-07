package ftp

// Anonymous user authenticator, implements goftp/server Auth interface
// Respond success on login request for anonymous user
type AuthAnonymous struct{}

// CheckPasswd answers success if login is "anonymous" or "ftp"
func (a AuthAnonymous) CheckPasswd(login string, pass string) (success bool, err error) {
	success = login == "anonymous" || login == "ftp"
	return
}
