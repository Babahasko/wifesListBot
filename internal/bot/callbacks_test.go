package bot

import (
	"testing"
)

func TestListCompactCBPack(t *testing.T) {
	listCBService := NewListCallbackService()
	listName := "Еда"
	cbStr, _ := listCBService.Pack(listName)
	expect := "list_Еда"
	if cbStr != expect {
		t.Errorf(`listCbService.Pack("Еда") expect: %q, get: %q`, expect, cbStr)
	}
}

func TestListCompactCBUnpack(t *testing.T) {
	cbStr := "list_Еда"
	listCBService := NewListCallbackService()
	listCBData := listCBService.Unpack(cbStr)
	expect := "Еда"
	if listCBData.Name != expect {
		t.Errorf(`listCbService.Unpack("Еда") expect :%q, get: %q`, expect, listCBData.Name)
	}
}
