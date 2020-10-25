package data

import (
	"fmt"
	"testing"
)

func TestDownload(t *testing.T) {
	Download("a0180935", "SCP-8900-ex")
}

func TestCollectSites(t *testing.T) {
	lic, err := newLoggedInClient("a0180935", "SCP-8900-ex")
	if err != nil {
		fmt.Println(err)
	}

	if _, err := collectSites(lic); err != nil {
		fmt.Println(err)
	}
}
