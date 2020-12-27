package resource_test

import (
	"fmt"
	"log"
	"pandora/pkg/resource"
	"testing"
)

func TestEncoding(t *testing.T) {
	code := 14

	fmt.Printf("code :%b\n", uint(code))

	r := resource.DecodeRejectableType(uint(code))

	log.Println("Video:", r.Video)
	log.Println("Audio:", r.Audio)
	log.Println("Excel:", r.Excel)
	log.Println("PowerPoint:", r.PowerPoint)
	log.Println("Word:", r.Word)

	log.Println("encode:", r.Encode())

	if r.Encode() != fmt.Sprint(code) {
		t.Error("encode is failed")
	}
}
