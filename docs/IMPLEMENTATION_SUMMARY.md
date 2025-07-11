# Implementação AES-CBC + AES-GCM - Backend

## 📋 Resumo das Modificações Implementadas

O backend `lockari-backend-app/pkg/crypt/crypt_server/crypt.go` foi estendido para suportar tanto **AES-CBC** (compatibilidade existente) quanto **AES-GCM** (nova implementação) com detecção automática de formato.

## 🚀 Funcionalidades Implementadas

### ✅ **1. Suporte Dual CBC/GCM**
- **AES-CBC**: Mantém compatibilidade total com implementação anterior
- **AES-GCM**: Nova implementação com autenticação integrada
- **Detecção Automática**: Identifica formato de payload automaticamente

### ✅ **2. Derivação de Chave SHA-256**
- Compatível com `crypt_client` (GCM)
- Chave original Base64 → SHA-256 → Chave de 32 bytes
- Cache de chave derivada para performance

### ✅ **3. Interface Estendida**
- `EncryptPayloadWithMode(data, mode)`: Criptografia com modo específico
- `DecryptPayloadGCM()`: Descriptografia GCM específica
- `SetCryptMode()` / `GetCryptMode()`: Gerenciamento de modo
- Retrocompatibilidade total com métodos existentes

### ✅ **4. CLI de Teste Atualizado**
- Flag `-crypt CBC|GCM`: Especifica modo de criptografia
- Modo interativo com seleção de criptografia
- Testes automáticos de compatibilidade

## 📊 Comparação CBC vs GCM

| Aspecto | AES-CBC | AES-GCM |
|---------|---------|---------|
| **Segurança** | ✅ Boa | ⭐ Excelente (autenticação integrada) |
| **Performance** | ✅ Rápida | ⭐ Mais rápida (sem padding) |
| **Tamanho Payload** | 🔸 Múltiplo de 16 bytes | ✅ Tamanho exato + 28 bytes overhead |
| **Compatibilidade** | ✅ Frontend atual | 🔸 Requer atualização frontend |
| **Formato** | `Base64(IV + ciphertext)` | `Base64(nonce + ciphertext + tag)` |

## 🔧 Uso da Nova API

### **Inicialização Básica (CBC - Compatibilidade)**
```go
cryptData, err := cryptserver.InicializationCryptData(&key)
```

### **Inicialização com Modo Específico**
```go
// Modo GCM
cryptData, err := cryptserver.InicializationCryptDataWithMode(&key, "GCM")

// Modo CBC explícito
cryptData, err := cryptserver.InicializationCryptDataWithMode(&key, "CBC")
```

### **Criptografia com Modo Específico**
```go
// GCM
encrypted, err := cryptData.EncryptPayloadWithMode(jsonData, "GCM")

// CBC
encrypted, err := cryptData.EncryptPayloadWithMode(jsonData, "CBC")

// Padrão (usa modo configurado na inicialização)
encrypted, err := cryptData.EncryptPayload(jsonData)
```

### **Descriptografia Automática**
```go
// Detecta automaticamente se é CBC ou GCM
decrypted, err := cryptData.PayloadData(base64Payload)
```

## 🧪 Testes Implementados

### **CLI de Teste Estendido**
```bash
# Criptografia GCM
go run main.go -mode encrypt -crypt GCM -payload '{"test": "data"}'

# Criptografia CBC
go run main.go -mode encrypt -crypt CBC -payload '{"test": "data"}'

# Descriptografia automática
go run main.go -mode decrypt -payload "Base64PayloadAqui"

# Modo interativo com seleção de modo
go run main.go -i
```

### **Casos de Teste Validados**
- ✅ **CBC → CBC**: Criptografia e descriptografia CBC
- ✅ **GCM → GCM**: Criptografia e descriptografia GCM  
- ✅ **Detecção Automática**: Distingue payloads CBC vs GCM
- ✅ **Retrocompatibilidade**: Payloads CBC antigos funcionam
- ✅ **Interoperabilidade**: Compatível com crypt_client (GCM)

## 📝 Logs de Debug Implementados

### **Inicialização**
```
DEBUG: Original encrypt key: "..." (length: X)
DEBUG: Cleaned encrypt key: "..." (length: X)
DEBUG: Initialized with default mode: CBC
DEBUG: Derived key for GCM mode: ...
```

### **Detecção de Formato**
```
DEBUG: Format detection - Payload length: X bytes
DEBUG: Format detection - Trying GCM first...
DEBUG: Format detection - GCM decryption successful
DEBUG: Detected payload format: GCM
```

### **Criptografia GCM**
```
DEBUG: GCM generated nonce: ...
DEBUG: GCM plaintext length: X bytes
DEBUG: GCM encryption successful, Base64 output length: X
```

## 🔐 Compatibilidade com Frontend

### **AES-CBC (Atual)**
- ✅ **Mantém total compatibilidade**
- ✅ **Mesmo formato de payload**
- ✅ **Mesma chave Base64**

### **AES-GCM (Migração)**
- 🔄 **Requer implementação frontend**
- ✅ **Derivação SHA-256 compatível**
- ✅ **Formato padronizado**

## 📈 Benefícios da Implementação

### **Imediatos**
1. **Retrocompatibilidade**: Nenhum código existente quebra
2. **Flexibilidade**: Suporte a ambos os modos
3. **Detecção Automática**: Sem necessidade de especificar formato
4. **Logs Detalhados**: Debug facilitado

### **Futuros**
1. **Segurança Aprimorada**: AES-GCM com autenticação
2. **Performance Melhor**: Sem padding PKCS7
3. **Interoperabilidade**: Compatível com crypt_client
4. **Preparação para Migração**: Base sólida para transição

## 🔄 Estratégia de Migração Sugerida

### **Fase 1: Implementação Frontend GCM**
- Implementar funções GCM no frontend (JavaScript/CryptoJS)
- Testes de compatibilidade com backend GCM
- Validação end-to-end

### **Fase 2: Migração Gradual**
- Novos dados: usar GCM por padrão
- Dados antigos: manter CBC (detecção automática)
- Monitoramento e logging

### **Fase 3: Consolidação (Futuro)**
- Eventual descontinuação do CBC
- Migração de dados legacy para GCM
- Simplificação do código

## 🎯 Próximos Passos

1. **Frontend**: Implementar funções AES-GCM compatíveis
2. **Testes**: Validação cruzada frontend ↔ backend
3. **Documentação**: Guias de migração para desenvolvedores
4. **Monitoramento**: Métricas de uso CBC vs GCM
5. **Otimização**: Performance tuning e caching

A implementação está completa e pronta para uso em produção, mantendo total compatibilidade com o sistema atual enquanto prepara o caminho para uma migração segura para AES-GCM.
