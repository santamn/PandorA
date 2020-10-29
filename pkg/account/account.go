package account

import (
	"bufio"
	"encoding/binary"
	"errors"
	"pandora/pkg/dir"
	"strings"
)

const (
	// アカウント情報を記録するファイルの名前
	accountFile = "account.dat"
)

// WriteAccountInfo アカウント情報を書き込む
func WriteAccountInfo(ecsID, password string) error {
	data := []byte(ecsID + ":" + password)

	file, err := dir.FetchFile(accountFile, "")
	if err != nil {
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, rot47(data)); err != nil {
		return err
	}

	return nil
}

// ReadAccountInfo アカウント情報の読み出しを行う
func ReadAccountInfo() (ecsID, password string, err error) {
	content := make([]byte, 0, 32)

	file, err := dir.FetchFile(accountFile, "")
	if err != nil {
		return "", "", err
	}

	buff := bufio.NewScanner(file)
	content = buff.Bytes()

	text := strings.Split(string(rot47(content)), ":")
	if len(text) != 2 {
		err = errors.New("Invalid format")
		return
	}

	ecsID, password = text[0], text[1]

	return
}

// ASCIIコードで33(!)-126(~)をrotする
func rot47(target []byte) (result []byte) {
	result = make([]byte, len(target))

	for i, t := range target {
		result[i] = (t-33+47)%94 + 33
	}

	return
}
