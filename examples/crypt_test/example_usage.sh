#!/bin/bash

echo "=== Exemplo de Uso Real com Payload do Frontend ==="
echo

cd /home/carlos.tomelin/codes/tomelin/lockari/lockari-backend-app/examples/crypt_test

# Simulando um payload complexo que poderia vir do frontend
FRONTEND_PAYLOAD='{"user_id":"12345","session_token":"abc123xyz","permissions":["read","write","admin"],"metadata":{"ip":"192.168.1.100","user_agent":"Mozilla/5.0","timestamp":"2025-07-09T21:55:00Z"}}'

echo "=== Cenário: Payload vindo do Frontend ==="
echo "Payload original do frontend:"
echo "$FRONTEND_PAYLOAD" | jq .
echo

echo "1. Criptografando payload antes de enviar ao backend..."
ENCRYPTED_PAYLOAD=$(go run main.go -mode encrypt -payload "$FRONTEND_PAYLOAD" 2>/dev/null | awk '{print $3}')
echo "Payload criptografado (que seria enviado via HTTP):"
echo "$ENCRYPTED_PAYLOAD"
echo

echo "2. No backend, descriptografando payload recebido..."
echo "Comando para descriptografar:"
echo "go run main.go -mode decrypt -payload \"$ENCRYPTED_PAYLOAD\""
echo

echo "Resultado da descriptografia:"
DECRYPTED_PAYLOAD=$(go run main.go -mode decrypt -payload "$ENCRYPTED_PAYLOAD" 2>/dev/null)
echo "$DECRYPTED_PAYLOAD"
echo

echo "=== Teste com chave específica do ambiente ==="
echo "Simulando chave de produção diferente..."
PROD_KEY="VGhpc0lzQTE2Qnl0ZUtleVRoaXNJc0ExNkJ5dGVJVgo=" # base64 para "productionKey123456789012345678901234567890"

echo "Criptografando com chave de produção..."
PROD_ENCRYPTED=$(go run main.go -key "$PROD_KEY" -mode encrypt -payload "$FRONTEND_PAYLOAD" 2>/dev/null | awk '{print $3}')
echo "Payload criptografado com chave de produção:"
echo "$PROD_ENCRYPTED"
echo

echo "Descriptografando com chave de produção..."
PROD_DECRYPTED=$(go run main.go -key "$PROD_KEY" -mode decrypt -payload "$PROD_ENCRYPTED" 2>/dev/null)
echo "Resultado:"
echo "$PROD_DECRYPTED"
echo

echo "=== Exemplo de Erro: Chave Incorreta ==="
echo "Tentando descriptografar payload de produção com chave padrão..."
echo "Comando: go run main.go -mode decrypt -payload \"$PROD_ENCRYPTED\""
echo "Erro esperado:"
go run main.go -mode decrypt -payload "$PROD_ENCRYPTED" 2>&1 | grep -E "(Erro|failed)" | head -1
echo

echo "=== Como usar na prática ==="
echo "1. Para testar payload do frontend:"
echo "   go run main.go -mode decrypt -payload \"PAYLOAD_AQUI\""
echo
echo "2. Para testar com chave específica:"
echo "   go run main.go -key \"SUA_CHAVE_BASE64\" -mode decrypt -payload \"PAYLOAD_AQUI\""
echo
echo "3. Para modo interativo (múltiplos testes):"
echo "   go run main.go -i"
echo
echo "4. Para criptografar dados de teste:"
echo "   go run main.go -mode encrypt -payload '{\"test\":\"data\"}'"
