package ftp

type AuthAnonymous struct{}

func (a AuthAnonymous) CheckPasswd(login string, pass string) (success bool, err error) {
	success = login == "anonymous" || login == "ftp"
	return
}
