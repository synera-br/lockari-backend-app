package main

import (
	"crypto/rand"
	"fmt"
	"strings"
)

func GenerateTestTenantID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate tenant ID: %w", err)
	}
	return fmt.Sprintf("%x", bytes), nil
}

func main() {
	fmt.Println("Testando geração de hashes de 32 caracteres:")
	fmt.Println(strings.Repeat("=", 50))

	for i := 0; i < 10; i++ {
		hash, err := GenerateTestTenantID()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
			continue
		}
		fmt.Printf("Hash %2d: %s (tamanho: %d)\n", i+1, hash, len(hash))
	}

	// Teste de bytes específicos para mostrar o padrão
	fmt.Println("\nExemplos de conversão byte->hex:")
	examples := [][]byte{
		{0x00, 0x00, 0x00, 0x00}, // 4 bytes zeros
		{0xFF, 0xFF, 0xFF, 0xFF}, // 4 bytes máximos
		{0x01, 0x23, 0x45, 0x67}, // 4 bytes sequenciais
	}

	for _, bytes := range examples {
		hex := fmt.Sprintf("%x", bytes)
		fmt.Printf("Bytes: %v -> Hex: %s (tamanho: %d)\n", bytes, hex, len(hex))
	}
}
