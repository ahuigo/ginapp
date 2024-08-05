package conf

import "testing"

func TestConf(t *testing.T) {
	if GetConf().Http.ReadTimeout == 0 {
		t.Errorf("GetConf().Http.ReadTimeout == 0")
	}
}
