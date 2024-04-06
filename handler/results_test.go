package handler

import "testing"

func Test_toID(t *testing.T) {
	actual := toID("https://docs.google.com/spreadsheets/d/testid/edit#gid=3799310")
	if actual != "testid" {
		t.Errorf("actual: %s, expected: %s", actual, "testid")
	}
}
