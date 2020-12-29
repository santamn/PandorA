package pandaapi_test

import (
	pandaapi "pandora/pkg/pandaAPI"
	"testing"
)

func TestLogin(t *testing.T) {
	_, err := pandaapi.NewLoggedInClient("asdfasdfa", "asdfasdfasd")

	if err != nil {
		t.Errorf(err.Error())
	}
}
