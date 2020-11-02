package resource_test

import (
	"pandora/pkg/resource"
	"testing"
)

func TestDownload(t *testing.T) {
	resource.Download("a0180935", "SCP-8900-ex")
}
