package utility

import "testing"

func TestGetNowTimestamp(t *testing.T) {
	nowTimestamp := GetNowTimestamp()
	t.Log(nowTimestamp)
}
