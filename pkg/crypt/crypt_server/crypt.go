package cryptserver

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64" // Usado para codificar/decodificar a chave e o payload.
	"errors"
	"fmt"
	"io"
	"strings"
)

// Constantes para modos de criptografia
const (
	CryptModeCBC     = "CBC"
	CryptModeGCM     = "GCM"
	DefaultCryptMode = CryptModeCBC // Manter compatibilidade
)

type CryptData struct {
	Payload    string `json:"payload" binding:"required"`
	encryptKey *string
	derivedKey []byte // Chave derivada SHA-256 para compatibilidade com crypt_client
	cryptMode  string // "CBC" ou "GCM"
}

func InicializationCryptData(encryptKey *string) (CryptDataInterface, error) {

	if encryptKey == nil || *encryptKey == "" {
		return nil, errors.New("token is nil")
	}

	// Limpeza AGRESSIVA: remover TODOS os espaços em branco, quebras de linha, \r, \t, etc.
	newEncryptKey := strings.ReplaceAll(*encryptKey, "\n", "")
	newEncryptKey = strings.ReplaceAll(newEncryptKey, "\r", "")
	newEncryptKey = strings.ReplaceAll(newEncryptKey, "\t", "")
	newEncryptKey = strings.TrimSpace(newEncryptKey)

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

	// Configurar modo padrão (CBC para compatibilidade)
	data.cryptMode = DefaultCryptMode

	return data, nil
}

// InicializationCryptDataWithMode inicializa CryptData com modo específico (CBC ou GCM)
func InicializationCryptDataWithMode(encryptKey *string, mode string) (CryptDataInterface, error) {
	if mode != CryptModeCBC && mode != CryptModeGCM {
		return nil, fmt.Errorf("invalid crypt mode: %s (must be %s or %s)", mode, CryptModeCBC, CryptModeGCM)
	}

	// Usar função base de inicialização
	cryptData, err := InicializationCryptData(encryptKey)
	if err != nil {
		return nil, err
	}

	// Cast para tipo concreto para acessar campos privados
	data, ok := cryptData.(*CryptData)
	if !ok {
		return nil, errors.New("failed to cast CryptData interface")
	}

	// Configurar modo de criptografia
	data.cryptMode = mode

	// Se for modo GCM, derivar chave usando SHA-256 (compatível com crypt_client)
	if mode == CryptModeGCM {
		derivedKey, err := data.deriveKeyFromBase64(*data.encryptKey)
		if err != nil {
			return nil, fmt.Errorf("failed to derive key for GCM mode: %w", err)
		}
		data.derivedKey = derivedKey
	}

	return data, nil
}

// deriveKeyFromBase64 deriva uma chave de 32 bytes usando SHA-256 (compatível com crypt_client)
func (c *CryptData) deriveKeyFromBase64(base64Key string) ([]byte, error) {
	// Decodificar Base64 para bytes originais
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key for derivation: %w", err)
	}

	// Aplicar SHA-256 para derivar chave de 32 bytes
	hash := sha256.Sum256(keyBytes)

	return hash[:], nil
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
// Ela espera que o payload seja uma string Base64 e detecta automaticamente o formato (CBC ou GCM).
// Retorna os dados descriptografados como um slice de bytes ou um erro em caso de falha.
// Se o token for inválido ou a descriptografia falhar, retorna um erro.
func (c *CryptData) PayloadData(base64Payload string) ([]byte, error) {

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

	// Detectar formato automaticamente
	detectedMode, err := c.detectPayloadFormat(base64Payload)
	if err != nil {
		return nil, fmt.Errorf("format detection failed: %w", err)
	}

	// Dispatch para método apropriado baseado na detecção
	var decryptedData []byte
	switch detectedMode {
	case CryptModeCBC:
		decryptedData, err = c.DecryptPayload(base64Payload, *c.encryptKey)
		if err != nil {
			return nil, fmt.Errorf("CBC decryption failed: %w", err)
		}
	case CryptModeGCM:
		decryptedData, err = c.DecryptPayloadGCM(base64Payload, *c.encryptKey)
		if err != nil {
			return nil, fmt.Errorf("GCM decryption failed: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported payload format: %s", detectedMode)
	}

	if len(decryptedData) == 0 {
		return nil, errors.New("decryption resulted in empty data (possible wrong key or corrupted data)")
	}

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

	keyBytes, err := base64.StdEncoding.DecodeString(base64KeyInput)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to decode base64 key: %w", err)
	}

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

	// Validação do tamanho mínimo (deve ter pelo menos o IV)
	if len(combinedBytes) < aes.BlockSize {
		return nil, fmt.Errorf("decrypt: combined payload too short to contain IV (got %d bytes, expected at least %d)", len(combinedBytes), aes.BlockSize)
	}

	iv := combinedBytes[:aes.BlockSize]
	ciphertext := combinedBytes[aes.BlockSize:]

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

	unpaddedData, err := c.pkcs7Unpad(ciphertext, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to unpad data: %w", err)
	}

	// Validar se os dados descriptografados não estão vazios
	if len(unpaddedData) == 0 {
		return nil, errors.New("decrypt: unpaddedData is empty (possible wrong key or corrupted data)")
	}

	return unpaddedData, nil
}

// DecryptPayloadGCM descriptografa um payload que foi criptografado usando AES-GCM.
// Espera-se que base64Payload seja Base64(nonce + ciphertext + tag).
func (c *CryptData) DecryptPayloadGCM(base64Payload string, base64KeyInput string) ([]byte, error) {

	// Validações de entrada
	if base64Payload == "" {
		return nil, errors.New("decrypt GCM: base64 payload is empty")
	}

	if base64KeyInput == "" {
		return nil, errors.New("decrypt GCM: base64 key is empty")
	}

	// Validar formato Base64
	if !isValidBase64(base64Payload) {
		return nil, errors.New("decrypt GCM: payload contains invalid Base64 characters")
	}

	if !isValidBase64(base64KeyInput) {
		return nil, errors.New("decrypt GCM: key contains invalid Base64 characters")
	}

	// Derivar chave usando SHA-256 (compatível com crypt_client)
	derivedKey, err := c.deriveKeyFromBase64(base64KeyInput)
	if err != nil {
		return nil, fmt.Errorf("decrypt GCM: failed to derive key: %w", err)
	}

	// Decodificar payload
	combined, err := base64.StdEncoding.DecodeString(base64Payload)
	if err != nil {
		return nil, fmt.Errorf("decrypt GCM: failed to decode base64 payload: %w", err)
	}

	// Criar cipher AES
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt GCM: failed to create AES cipher: %w", err)
	}

	// Criar modo GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("decrypt GCM: failed to create GCM mode: %w", err)
	}

	// Validar tamanho mínimo (nonce + pelo menos algum dado + tag)
	if len(combined) < gcm.NonceSize() {
		return nil, fmt.Errorf("decrypt GCM: payload too short to contain nonce (got %d bytes, expected at least %d)", len(combined), gcm.NonceSize())
	}

	// Extrair nonce e ciphertext (que inclui o tag no final)
	nonce := combined[:gcm.NonceSize()]
	ciphertext := combined[gcm.NonceSize():]

	// Descriptografar e validar autenticação
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt GCM: authentication failed or decryption error: %w", err)
	}

	// Validar se os dados descriptografados não estão vazios
	if len(plaintext) == 0 {
		return nil, errors.New("decrypt GCM: decryption resulted in empty data")
	}

	return plaintext, nil
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

// EncryptPayloadGCM criptografa os jsonDataBytes fornecidos usando AES-GCM.
// O output é uma string Base64 no formato: Base64(nonce + ciphertext + tag).
func (c *CryptData) EncryptPayloadGCM(jsonDataBytes []byte) (string, error) {

	// Validações de entrada
	if len(jsonDataBytes) == 0 {
		return "", errors.New("encrypt GCM: jsonDataBytes is empty")
	}

	if c.encryptKey == nil || *c.encryptKey == "" {
		return "", errors.New("encrypt GCM: package key (encryptKey) is empty")
	}

	// Usar chave derivada se disponível, senão derivar agora
	var derivedKey []byte
	var err error

	if c.derivedKey != nil {
		derivedKey = c.derivedKey
	} else {
		derivedKey, err = c.deriveKeyFromBase64(*c.encryptKey)
		if err != nil {
			return "", fmt.Errorf("encrypt GCM: failed to derive key: %w", err)
		}
	}

	// Criar cipher AES
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", fmt.Errorf("encrypt GCM: failed to create AES cipher: %w", err)
	}

	// Criar modo GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt GCM: failed to create GCM mode: %w", err)
	}

	// Gerar nonce aleatório (12 bytes para GCM)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt GCM: failed to generate nonce: %w", err)
	}

	// Criptografar dados (já inclui autenticação/tag automaticamente)
	ciphertext := gcm.Seal(nil, nonce, jsonDataBytes, nil)

	// Combinar nonce + ciphertext+tag
	combined := append(nonce, ciphertext...)

	// Codificar para Base64
	base64EncryptedPayload := base64.StdEncoding.EncodeToString(combined)

	// Validação final do resultado
	if base64EncryptedPayload == "" {
		return "", errors.New("encrypt GCM: failed to generate Base64 output")
	}

	return base64EncryptedPayload, nil
}

// EncryptPayloadWithMode criptografa jsonDataBytes usando o modo especificado (CBC ou GCM).
func (c *CryptData) EncryptPayloadWithMode(jsonDataBytes []byte, mode string) (string, error) {

	switch mode {
	case CryptModeCBC:
		return c.EncryptPayload(jsonDataBytes) // Usar implementação CBC existente
	case CryptModeGCM:
		return c.EncryptPayloadGCM(jsonDataBytes)
	default:
		return "", fmt.Errorf("unsupported encryption mode: %s", mode)
	}
}

// SetCryptMode define o modo de criptografia (CBC ou GCM)
func (c *CryptData) SetCryptMode(mode string) error {
	if mode != CryptModeCBC && mode != CryptModeGCM {
		return fmt.Errorf("invalid crypt mode: %s (must be %s or %s)", mode, CryptModeCBC, CryptModeGCM)
	}

	c.cryptMode = mode

	// Se mudou para GCM e não tem chave derivada, derivar agora
	if mode == CryptModeGCM && c.derivedKey == nil && c.encryptKey != nil {
		derivedKey, err := c.deriveKeyFromBase64(*c.encryptKey)
		if err != nil {
			return fmt.Errorf("failed to derive key for GCM mode: %w", err)
		}
		c.derivedKey = derivedKey
	}

	return nil
}

// GetCryptMode retorna o modo de criptografia atual
func (c *CryptData) GetCryptMode() string {
	if c.cryptMode == "" {
		return DefaultCryptMode // Retornar padrão se não configurado
	}
	return c.cryptMode
}

// detectPayloadFormat detecta automaticamente se o payload é CBC ou GCM tentando descriptografar
func (c *CryptData) detectPayloadFormat(base64Payload string) (string, error) {
	combinedBytes, err := base64.StdEncoding.DecodeString(base64Payload)
	if err != nil {
		return "", fmt.Errorf("failed to decode payload for format detection: %w", err)
	}

	if len(combinedBytes) < 29 {
		return "", fmt.Errorf("payload too short for any format: length %d (minimum 29 bytes)", len(combinedBytes))
	}

	// Estratégia pragmática: tentar GCM primeiro (mais específico)
	// Se falhar, tentar CBC

	// Tentar GCM primeiro
	_, err = c.DecryptPayloadGCM(base64Payload, *c.encryptKey)
	if err == nil {
		return CryptModeGCM, nil
	}

	// Se GCM falhou, tentar CBC
	_, err = c.DecryptPayload(base64Payload, *c.encryptKey)
	if err == nil {
		return CryptModeCBC, nil
	}

	// Se ambos falharam, retornar erro
	return "", fmt.Errorf("unable to decrypt payload with either CBC or GCM")
}

// ...existing code...

// ...existing code...

// Helper function para Go < 1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
