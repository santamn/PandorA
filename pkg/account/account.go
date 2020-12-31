package account

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"pandora/pkg/dir"
	"pandora/pkg/resource"
	"strconv"
	"strings"
)

const (
	// アカウント情報を記録するファイルの名前
	accountFile = "account.dat"
)

// WriteAccountInfo アカウント情報を書き込む
func WriteAccountInfo(ecsID, password string, rejectable *resource.RejectableType) error {
	rejectNum := rejectable.Encode()
	data := []byte(ecsID + ":" + password + ":" + rejectNum)

	file, err := dir.FetchSettingsFile(accountFile)
	if err != nil {
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, rot47(data)); err != nil {
		return err
	}

	return nil
}

// ReadAccountInfo アカウント情報の読み出しを行う
func ReadAccountInfo() (ecsID, password string, rejectable *resource.RejectableType, err error) {
	file, err := dir.FetchSettingsFile(accountFile)
	if err != nil {
		return "", "", nil, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "", nil, err
	}

	text := strings.SplitN(string(rot47(content)), ":", 3)
	if len(text) != 3 {
		err = errors.New("Invalid format")
		return
	}

	var rejectNum string
	ecsID, password, rejectNum = text[0], text[1], text[2]
	num, err := strconv.Atoi(rejectNum)
	if err != nil {
		return ecsID, password, nil, err
	}
	rejectable = resource.DecodeRejectableType(uint(num))

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
