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

func TestCopyright(t *testing.T) {
	lic, err := newLoggedInClient("a0180935", "SCP-8900-ex")
	if err != nil {
		fmt.Println(err)
	}

	s := site{
		Title: "[2020後期水１]統計物理学",
		ID:    "2020-888-N234-001",
	}
	r, err := collectUnacquiredResouceInfo(lic, []site{s})
	if err != nil {
		fmt.Println(err)
	}

	if errors := paraDownloadPDF(lic, r); len(errors) > 0 {
		fmt.Println(errors)
	}
}
