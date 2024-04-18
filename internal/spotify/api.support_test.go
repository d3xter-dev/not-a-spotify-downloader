package Spotify

import (
	"bytes"
	"testing"
)

func TestConvert62(t *testing.T) {
	if !bytes.Equal([]byte{0x00, 0x0d, 0x53, 0x65, 0x35, 0x86, 0x4e, 0x0f, 0x99, 0x76, 0x1f, 0x9d, 0xa9, 0x00, 0xb1, 0xc1}, Convert62("0065zxtT6XKaQww7cLne0h")) {
		t.Fail()
	}
}

func TestConvertTo62(t *testing.T) {
	if ConvertTo62([]byte{0x00, 0x0d, 0x53, 0x65, 0x35, 0x86, 0x4e, 0x0f, 0x99, 0x76, 0x1f, 0x9d, 0xa9, 0x00, 0xb1, 0xc1}) != "0065zxtT6XKaQww7cLne0h" {
		t.Fail()
	}
}
