#!/bin/bash

echo "=== Teste do CLI de Criptografia/Descriptografia ==="
echo

# Navegar para o diretório correto
cd /home/carlos.tomelin/codes/tomelin/lockari/lockari-backend-app/examples/crypt_test

echo "1. Testando criptografia de JSON..."
PAYLOAD='{"user":"carlos","role":"admin","timestamp":"2025-07-09T21:50:00Z"}'
echo "Payload original: $PAYLOAD"

ENCRYPTED=$(go run main.go -mode encrypt -payload "$PAYLOAD" 2>/dev/null | awk '{print $3}')
echo "Payload criptografado: $ENCRYPTED"
echo

echo "2. Testando descriptografia..."
DECRYPTED=$(go run main.go -mode decrypt -payload "$ENCRYPTED" 2>/dev/null)
echo "Payload descriptografado:"
echo "$DECRYPTED"
echo

echo "3. Testando com chave personalizada..."
CUSTOM_KEY="VGhpc0lzQTE2Qnl0ZUtleVRoaXNJc0ExNkJ5dGVJVgo=" # base64 para "abcdefghijklmnopqrstuvwxyz123456"
ENCRYPTED_CUSTOM=$(go run main.go -key "$CUSTOM_KEY" -mode encrypt -payload '{"test":"custom_key"}' 2>/dev/null | awk '{print $3}')
echo "Criptografado com chave personalizada: $ENCRYPTED_CUSTOM"

DECRYPTED_CUSTOM=$(go run main.go -key "$CUSTOM_KEY" -mode decrypt -payload "$ENCRYPTED_CUSTOM" 2>/dev/null)
echo "Descriptografado com chave personalizada:"
echo "$DECRYPTED_CUSTOM"
echo

echo "4. Testando erro com chave incorreta..."
echo "Tentando descriptografar com chave padrão um payload criptografado com chave personalizada:"
go run main.go -mode decrypt -payload "$ENCRYPTED_CUSTOM" 2>&1 | head -5
echo

echo "=== Teste concluído ==="
echo
echo "Para teste interativo, execute:"
echo "go run main.go -i"
