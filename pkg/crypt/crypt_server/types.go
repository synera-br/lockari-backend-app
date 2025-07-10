package cryptserver

type CryptDataInternalInterface interface {
	validateTokenFromString(token *string) error
	pkcs7Unpad(data []byte, blockSize int) ([]byte, error)
	pkcs7Pad(data []byte, blockSize int) ([]byte, error)
	deriveKeyFromBase64(base64Key string) ([]byte, error)
	detectPayloadFormat(base64Payload string) (string, error)
}

type CryptDataInterface interface {
	// Métodos existentes (compatibilidade)
	PayloadData(base64Payload string) ([]byte, error)
	DecryptPayload(base64Payload string, base64KeyInput string) ([]byte, error)
	EncryptPayload(jsonDataBytes []byte) (string, error)

	// Novos métodos para suporte AES-GCM
	EncryptPayloadWithMode(jsonDataBytes []byte, mode string) (string, error)
	DecryptPayloadGCM(base64Payload string, base64KeyInput string) ([]byte, error)
	EncryptPayloadGCM(jsonDataBytes []byte) (string, error)
	SetCryptMode(mode string) error
	GetCryptMode() string

	CryptDataInternalInterface
}
