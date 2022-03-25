package provider

import (
	"testing"
)

func TestIsRightSmsCodeExists(t *testing.T) {
	if !IsRightCode("Rond", "Sms") {
		t.Error("Существующий код не распознался")
	}
}

func TestIsRightSmsCodeNotExists(t *testing.T) {
	if IsRightCode("Megafon", "Sms") {
		t.Error("Несуществующий код не распознался")
	}
}

func TestIsRightMmsCodeExists(t *testing.T) {
	if !IsRightCode("Topolo", "Mms") {
		t.Error("Существующий код не распознался")
	}
}

func TestIsRightMmsCodeNotExists(t *testing.T) {
	if IsRightCode("Tele2", "Mms") {
		t.Error("Несуществующий код не распознался")
	}
}

func TestIsRightVoiceCallCodeExists(t *testing.T) {
	if !IsRightCode("JustPhone", "VoiceCall") {
		t.Error("Существующий код не распознался")
	}
}

func TestIsRightVoiceCallCodeNotExists(t *testing.T) {
	if IsRightCode("Beeline", "VoiceCall") {
		t.Error("Несуществующий код не распознался")
	}
}

func TestIsRightEmailCodeExists(t *testing.T) {
	if !IsRightCode("Gmail", "Email") {
		t.Error("Существующий код не распознался")
	}
}

func TestIsRightEmailCodeNotExists(t *testing.T) {
	if IsRightCode("Rambler", "Email") {
		t.Error("Несуществующий код не распознался")
	}
}
