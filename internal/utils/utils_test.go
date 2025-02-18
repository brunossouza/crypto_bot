package utils

import (
	"os"
	"os/exec"
	"testing"
)

// TestParseFloat valida a conversão correta de string para float64.
func TestParseFloat(t *testing.T) {
	// ...existing code...
	val := ParseFloat("1.23")
	if val != 1.23 {
		t.Errorf("Esperado 1.23, obteve %f", val)
	}
}

// TestParseFloat_Invalid valida que ParseFloat encerra o programa com um input inválido.
func TestParseFloat_Invalid(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Chama ParseFloat com string inválida para acionar log.Fatal
		ParseFloat("abc")
		return
	}
	// Executa este teste em um subprocesso.
	cmd := exec.Command(os.Args[0], "-test.run=TestParseFloat_Invalid")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Esperava-se que o processo terminasse com erro ao converter string inválida")
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() == 0 {
			t.Fatalf("O código de saída não pode ser zero para input inválido")
		}
	} else {
		t.Fatalf("Erro inesperado: %v", err)
	}
}
