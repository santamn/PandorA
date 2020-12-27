package pandaapi

import "fmt"

// DeadPandAError PandAが死んでいる時に返すエラー
type DeadPandAError struct {
	code int
	err  error
	url  string
}

func (d *DeadPandAError) Error() string {
	if d.err != nil {
		return fmt.Sprintf("Panda is dead. Status code %d in %s\nerror: %s\n", d.code, d.url, d.err)
	}

	return fmt.Sprintf("Panda is dead. Status code %d: in %s\n", d.code, d.url)
}

// FailedLoginError ログインに失敗したときのエラー
type FailedLoginError struct {
	EscID    string
	Password string
}

func (f *FailedLoginError) Error() string {
	return fmt.Sprintf(
		"Login failed. Please confirm your EcsID and password.\nEcsID: %s\nPassword: %s",
		f.EscID,
		f.Password,
	)
}

// NetworkError ネットの接続状態のエラー
type NetworkError struct {
	err error
}

func (n *NetworkError) Error() string {
	return fmt.Sprintf("Network Error:%s", n.err.Error())
}
