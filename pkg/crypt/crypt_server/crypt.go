package cryptserver

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64" // Usado para codificar/decodificar a chave e o payload.
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type CryptData struct {
	Payload    string `json:"payload" binding:"required"`
	encryptKey *string
}

func InicializationCryptData(encryptKey *string) (CryptDataInterface, error) {

	if encryptKey == nil || *encryptKey == "" {
		return nil, errors.New("token is nil")
	}

	// Log da chave original para debug
	log.Printf("DEBUG: Original encrypt key: %q (length: %d)", *encryptKey, len(*encryptKey))

	// Limpar espaços em branco e quebras de linha (como no frontend)
	newEncryptKey := strings.TrimSpace(strings.ReplaceAll(*encryptKey, "\n", ""))

	// Log da chave após limpeza
	log.Printf("DEBUG: Cleaned encrypt key: %q (length: %d)", newEncryptKey, len(newEncryptKey))

	// Validar se a chave após limpeza não ficou vazia
	if newEncryptKey == "" {
		return nil, errors.New("token is empty after trimming whitespace")
	}

	data := &CryptData{}
	err := data.validateTokenFromString(&newEncryptKey)
	if err != nil {
		return nil, err
	}

	// CRÍTICO: Armazenar a chave LIMPA, não a original
	data.encryptKey = &newEncryptKey

	// Log da chave final armazenada
	log.Printf("DEBUG: Final stored encrypt key: %q", *data.encryptKey)

	return data, nil
}

func (c *CryptData) validateTokenFromString(token *string) error {
	if token == nil || *token == "" {
		return errors.New("token is nil or empty")
	}

	// Validar formato Base64 básico antes de tentar decodificar
	tokenValue := *token
	if len(tokenValue) == 0 {
		return errors.New("token is empty")
	}

	// Validar caracteres Base64 válidos
	if !isValidBase64(tokenValue) {
		return errors.New("token contains invalid Base64 characters")
	}

	decodeToken, err := base64.StdEncoding.DecodeString(tokenValue)
	if err != nil {
		return fmt.Errorf("failed to decode base64 key: %w", err)
	}

	if len(decodeToken) == 0 {
		return errors.New("decoded token is empty")
	}

	// Validar tamanhos de chave AES suportados (mesma validação do frontend)
	switch len(decodeToken) {
	case 16, 24, 32:
		// Tamanhos válidos para AES-128, AES-192, AES-256
		return nil
	default:
		return fmt.Errorf("invalid AES key size: %d bytes (must be 16, 24, or 32)", len(decodeToken))
	}
}

// isValidBase64 verifica se a string contém apenas caracteres Base64 válidos
func isValidBase64(s string) bool {
	// Regex para Base64 válido: A-Z, a-z, 0-9, +, /, e = no final
	for _, char := range s {
		if !((char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '+' || char == '/' || char == '=') {
			return false
		}
	}
	return true
}

// PayloadData é uma função wrapper que descriptografa usando a chave global do pacote.
// Ela espera que o payload seja uma string Base64 no formato: Base64(bytes_crus_IV + bytes_crus_Ciphertext).
// Retorna os dados descriptografados como um slice de bytes ou um erro em caso de falha.
// Se o token for inválido ou a descriptografia falhar, retorna um erro.
func (c *CryptData) PayloadData(base64Payload string) ([]byte, error) {
	log.Println("Decrypting payload with global key...")

	// Validações de entrada (consistentes com o frontend)
	if base64Payload == "" {
		return nil, errors.New("decrypt: base64 payload is empty")
	}

	if !isValidBase64(base64Payload) {
		return nil, errors.New("decrypt: payload contains invalid Base64 characters")
	}

	if c.encryptKey == nil || *c.encryptKey == "" {
		return nil, errors.New("decrypt: package key (encryptKey) is empty")
	}

	log.Println("Using package key for decryption:", *c.encryptKey)
	decryptedData, err := c.DecryptPayload(base64Payload, *c.encryptKey)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	if len(decryptedData) == 0 {
		return nil, errors.New("decryption resulted in empty data (possible wrong key or corrupted data)")
	}
	log.Println("Decryption successful, returning data...")
	log.Println("Decrypted data size:", len(decryptedData))

	return decryptedData, nil
}

// pkcs7Unpad remove o padding PKCS7 dos dados.
func (c *CryptData) pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("pkcs7Unpad: input data is empty")
	}

	if blockSize <= 0 || blockSize > 255 {
		return nil, fmt.Errorf("pkcs7Unpad: invalid block size %d", blockSize)
	}

	paddingLen := int(data[len(data)-1])

	// Validações mais rigorosas do padding
	if paddingLen == 0 || paddingLen > blockSize || paddingLen > len(data) {
		return nil, errors.New("pkcs7Unpad: invalid padding length (possible wrong key or corrupted data)")
	}

	// Validar se todos os bytes de padding são iguais
	for i := len(data) - paddingLen; i < len(data); i++ {
		if data[i] != byte(paddingLen) {
			return nil, errors.New("pkcs7Unpad: invalid padding bytes (possible wrong key or corrupted data)")
		}
	}

	return data[:len(data)-paddingLen], nil
}

// DecryptPayload descriptografa um payload que foi criptografado usando AES-CBC.
// Espera-se que base64Payload seja Base64(bytes_crus_IV + bytes_crus_Ciphertext).
func (c *CryptData) DecryptPayload(base64Payload string, base64KeyInput string) ([]byte, error) {
	// Validações de entrada mais rigorosas
	if base64Payload == "" {
		return nil, errors.New("decrypt: base64 payload is empty")
	}

	if base64KeyInput == "" {
		return nil, errors.New("decrypt: base64 key is empty")
	}

	// Validar formato Base64 do payload
	if !isValidBase64(base64Payload) {
		return nil, errors.New("decrypt: payload contains invalid Base64 characters")
	}

	// Validar formato Base64 da chave
	if !isValidBase64(base64KeyInput) {
		return nil, errors.New("decrypt: key contains invalid Base64 characters")
	}

	// Log detalhado da chave para debug
	log.Printf("DEBUG: Raw base64 key input: %q (length: %d)", base64KeyInput, len(base64KeyInput))

	keyBytes, err := base64.StdEncoding.DecodeString(base64KeyInput)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to decode base64 key: %w", err)
	}

	log.Printf("DEBUG: Decoded key bytes: %v (length: %d)", keyBytes, len(keyBytes))
	log.Printf("DEBUG: Key as string: %q", string(keyBytes))
	log.Printf("DEBUG: Key hex: %x", keyBytes)

	// Validação rigorosa do tamanho da chave (mesma do frontend)
	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		return nil, fmt.Errorf("decrypt: invalid AES key size: %d bytes (must be 16, 24, or 32)", len(keyBytes))
	}

	combinedBytes, err := base64.StdEncoding.DecodeString(base64Payload)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to decode base64 payload: %w", err)
	}

	log.Printf("DEBUG: Combined bytes length: %d", len(combinedBytes))
	log.Printf("DEBUG: First 32 bytes (IV): %x", combinedBytes[:min(32, len(combinedBytes))])

	// Validação do tamanho mínimo (deve ter pelo menos o IV)
	if len(combinedBytes) < aes.BlockSize {
		return nil, fmt.Errorf("decrypt: combined payload too short to contain IV (got %d bytes, expected at least %d)", len(combinedBytes), aes.BlockSize)
	}

	iv := combinedBytes[:aes.BlockSize]
	ciphertext := combinedBytes[aes.BlockSize:]

	log.Printf("DEBUG: IV: %x", iv)
	log.Printf("DEBUG: Ciphertext length: %d", len(ciphertext))

	// Validar se o ciphertext não está vazio
	if len(ciphertext) == 0 {
		return nil, errors.New("decrypt: ciphertext is empty after IV extraction")
	}

	// Validar se o ciphertext é múltiplo do block size
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("decrypt: ciphertext length (%d) is not a multiple of AES block size (%d)", len(ciphertext), aes.BlockSize)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to create AES cipher: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext) // Descriptografa in-place

	log.Printf("DEBUG: Decrypted ciphertext (before unpadding): %x", ciphertext)
	log.Printf("DEBUG: Last 16 bytes of decrypted data: %x", ciphertext[len(ciphertext)-16:])

	unpaddedData, err := c.pkcs7Unpad(ciphertext, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to unpad data: %w", err)
	}

	// Validar se os dados descriptografados não estão vazios
	if len(unpaddedData) == 0 {
		return nil, errors.New("decrypt: unpaddedData is empty (possible wrong key or corrupted data)")
	}

	log.Printf("DEBUG: Final unpaddedData length: %d", len(unpaddedData))
	log.Printf("DEBUG: Final unpaddedData as string: %s", string(unpaddedData))

	return unpaddedData, nil
}

// Helper function para Go < 1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// pkcs7Pad adiciona padding PKCS7 aos dados.
func (c *CryptData) pkcs7Pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 || blockSize > 255 { // blockSize > 255 porque o byte de padding não pode exceder 255
		return nil, fmt.Errorf("pkcs7Pad: invalid block size %d (must be > 0 and <= 255)", blockSize)
	}
	padding := blockSize - (len(data) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...), nil
}

// EncryptPayload criptografa os jsonDataBytes fornecidos (que devem ser dados JSON já serializados)
// usando AES-CBC. O output é uma string Base64 no formato: Base64(bytes_crus_IV + bytes_crus_Ciphertext).
// Utiliza a chave global 'base64Key' definida no pacote.
func (c *CryptData) EncryptPayload(jsonDataBytes []byte) (string, error) {
	// Validações de entrada (consistentes com o frontend)
	if len(jsonDataBytes) == 0 {
		return "", errors.New("encrypt: jsonDataBytes is empty")
	}

	if c.encryptKey == nil || *c.encryptKey == "" {
		return "", errors.New("encrypt: package key (encryptKey) is empty")
	}

	// 1. Decodificar a chave de Base64 para bytes (usando a constante do pacote)
	keyBytes, err := base64.StdEncoding.DecodeString(*c.encryptKey)
	if err != nil {
		// Este erro seria inesperado se a constante base64Key for válida
		return "", fmt.Errorf("encrypt: failed to decode package base64 key: %w", err)
	}

	// Validar o tamanho da chave (mesmo que frontend)
	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		// Este erro seria inesperado se a constante base64Key for válida
		return "", fmt.Errorf("encrypt: invalid package AES key size: %d bytes (must be 16, 24, or 32)", len(keyBytes))
	}

	// 2. Os jsonDataBytes já são fornecidos, não precisamos fazer json.Marshal aqui.

	// 3. Aplicar Padding PKCS7 aos dados JSON
	paddedData, err := c.pkcs7Pad(jsonDataBytes, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("encrypt: failed to apply PKCS7 padding: %w", err)
	}

	// 4. Criar o cipher AES
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("encrypt: failed to create AES cipher: %w", err)
	}

	// 5. Gerar um IV aleatório
	iv := make([]byte, aes.BlockSize) // aes.BlockSize é 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("encrypt: failed to generate IV: %w", err)
	}

	// 6. Criptografar os dados usando o modo CBC
	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	// --- INÍCIO DAS ALTERAÇÕES CRÍTICAS ---
	// 7. Combinar IV (bytes crus) e Ciphertext (bytes crus)
	combinedRawBytes := append(iv, ciphertext...)

	// 8. Codificar os bytes combinados para Base64
	// Este formato é Base64(bytes_crus_IV + bytes_crus_Ciphertext),
	// que é o que DecryptPayload espera e o que o frontend (CryptoJS) originalmente produzia.
	base64EncryptedPayload := base64.StdEncoding.EncodeToString(combinedRawBytes)
	// --- FIM DAS ALTERAÇÕES CRÍTICAS ---

	// Validação final do resultado (como no frontend)
	if base64EncryptedPayload == "" {
		return "", errors.New("encrypt: failed to generate Base64 output")
	}

	return base64EncryptedPayload, nil
}
