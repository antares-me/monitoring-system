package country

import (
	"testing"
)

func TestIsRightCodeExists(t *testing.T) {
	if !IsRightCode("RU") {
		t.Error("Существующий код не распознался")
	}
}

func TestIsRightCodeNotExists(t *testing.T) {
	if IsRightCode("ZU") {
		t.Error("Несуществующий код не распознался")
	}
}
