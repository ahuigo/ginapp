package nettool

import "testing"

func TestGetIp(t *testing.T) {
	ip := GetLocalIP()
	t.Log(ip)
}
