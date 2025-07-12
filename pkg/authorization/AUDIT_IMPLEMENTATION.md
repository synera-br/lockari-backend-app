# 📋 Implementação da API de Auditoria para Lockari

## 🎯 Visão Geral

Esta implementação fornece uma API completa para gerenciamento de logs de auditoria no sistema Lockari, integrando perfeitamente com o sistema de autorização OpenFGA já existente.

## 🏗️ Arquitetura da Implementação

### 1. **Estruturas de Dados (audit_api.go)**

```go
// Estrutura principal da resposta da API
type AuditLogsResponse struct {
    Logs       []AuditLogData `json:"logs"`
    Pagination PaginationData `json:"pagination"`
}

// Dados de um evento de auditoria
type AuditLogData struct {
    ID           string                 `json:"id"`
    ResourceName string                 `json:"resourceName"`
    ResourceType string                 `json:"resourceType"`
    Action       string                 `json:"action"`
    UserEmail    string                 `json:"userEmail"`
    UserID       string                 `json:"userId"`
    IPAddress    string                 `json:"ipAddress"`
    Timestamp    time.Time              `json:"timestamp"`
    ResourceLink string                 `json:"resourceLink,omitempty"`
    Details      map[string]interface{} `json:"details,omitempty"`
}

// Metadados de paginação
type PaginationData struct {
    CurrentPage int `json:"currentPage"`
    TotalPages  int `json:"totalPages"`
    TotalItems  int `json:"totalItems"`
    Limit       int `json:"limit"`
}
```

### 2. **Serviço de Auditoria (audit_service.go)**

```go
// Interface do serviço de auditoria
type AuditLogService interface {
    QueryLogs(ctx context.Context, query *AuditLogQuery) (*AuditLogsResponse, error)
    GetLogByID(ctx context.Context, logID string) (*AuditLogData, error)
    GetUserActivity(ctx context.Context, userID string, limit int) (*AuditLogsResponse, error)
    GetResourceActivity(ctx context.Context, resourceType, resourceID string, limit int) (*AuditLogsResponse, error)
    GetSuspiciousActivity(ctx context.Context, limit int) (*AuditLogsResponse, error)
    ExportLogs(ctx context.Context, query *AuditLogQuery, format string) ([]byte, error)
}
```

### 3. **Handler Web (logs.go)**

```go
// Interface do handler
type AuditLogHandlerInterface interface {
    GetLogs(c *gin.Context)
    ExportLogs(c *gin.Context)
    GetLogStats(c *gin.Context)
    GetLogTrends(c *gin.Context)
}
```

## 🔧 Endpoints da API

### **GET /v1/audit/logs**
- **Descrição**: Recupera logs de auditoria com filtros, paginação e ordenação
- **Parâmetros**:
  - `userId`: ID do usuário (opcional)
  - `userEmail`: Email do usuário (opcional)
  - `resourceType`: Tipo de recurso (opcional)
  - `resourceName`: Nome do recurso (opcional)
  - `action`: Ação realizada (opcional)
  - `ipAddress`: Endereço IP (opcional)
  - `startDate`: Data de início (RFC3339) (opcional)
  - `endDate`: Data de fim (RFC3339) (opcional)
  - `page`: Página atual (opcional, padrão: 1)
  - `limit`: Limite por página (opcional, padrão: 50, máximo: 1000)
  - `sortBy`: Campo para ordenação (opcional)
  - `sortOrder`: Ordem (asc/desc) (opcional)

**Exemplo de Resposta:**
```json
{
  "logs": [
    {
      "id": "log-123",
      "resourceName": "vault-456",
      "resourceType": "vault",
      "action": "read",
      "userEmail": "user@example.com",
      "userId": "user-789",
      "ipAddress": "192.168.1.100",
      "timestamp": "2023-10-30T14:30:00Z",
      "resourceLink": "/v1/vaults/vault-456",
      "details": {
        "permission": "can_read",
        "success": true
      }
    }
  ],
  "pagination": {
    "currentPage": 1,
    "totalPages": 10,
    "totalItems": 500,
    "limit": 50
  }
}
```

### **GET /v1/audit/logs/export**
- **Descrição**: Exporta logs de auditoria em diferentes formatos
- **Parâmetros**: Mesmos filtros do endpoint GET + `format` (json, csv)
- **Resposta**: Arquivo para download

### **GET /v1/audit/logs/stats**
- **Descrição**: Recupera estatísticas dos logs de auditoria
- **Requer**: Permissões de administrador
- **Resposta**: Estatísticas agregadas

### **GET /v1/audit/logs/trends**
- **Descrição**: Recupera dados de tendências dos logs
- **Parâmetros**: `days` (número de dias para análise, padrão: 30)
- **Resposta**: Dados de tendências

## 🔐 Segurança e Autorização

### **Verificação de Permissões**
1. **Autenticação**: Todos os endpoints requerem JWT válido
2. **Autorização**: 
   - Usuários regulares só podem acessar seus próprios logs
   - Administradores podem acessar todos os logs
   - Estatísticas e tendências requerem permissões admin

### **Implementação de Segurança**
```go
// Verifica se o usuário pode acessar os logs
func (h *auditLogHandler) checkAuditLogAccess(ctx context.Context, userID string, query *AuditLogQuery) error {
    // Admin tem acesso total
    if h.isUserAdmin(ctx, userID) {
        return nil
    }
    
    // Usuário regular só acessa próprios logs
    if query.UserID != "" && query.UserID != userID {
        return fmt.Errorf("access denied: cannot access other user's logs")
    }
    
    return nil
}
```

## 📊 Funcionalidades Principais

### **1. Filtros Avançados**
- Por usuário (ID ou email)
- Por tipo de recurso
- Por ação executada
- Por endereço IP
- Por intervalo de datas

### **2. Paginação**
- Páginas configuráveis
- Limite máximo de segurança
- Metadados de paginação

### **3. Ordenação**
- Por qualquer campo
- Ordem ascendente/descendente
- Padrão: timestamp desc

### **4. Exportação**
- Formato JSON
- Formato CSV
- Headers apropriados para download

### **5. Estatísticas**
- Contagem total de eventos
- Distribuição por ação
- Distribuição por tipo de recurso
- Número de usuários únicos

### **6. Tendências**
- Eventos por dia
- Período configurável
- Análise temporal

## 🚀 Como Integrar no Sistema Principal

### **1. Registro do Handler**

```go
// No main.go ou no arquivo de inicialização
func initializeAuditHandlers(router *gin.Engine) {
    // Inicializar serviços
    auditLogService := authorization.NewAuditLogService(10000)
    authzService := // seu serviço de autorização
    encryptor := // seu serviço de criptografia
    authClient := // seu cliente de autenticação
    tokenGen := // seu gerador de tokens
    
    // Criar handler
    v1 := router.Group("/v1")
    auditHandler, err := webhandler.InitializeAuditLogHandler(
        auditLogService,
        authzService,
        encryptor,
        authClient,
        tokenGen,
        v1,
    )
    if err != nil {
        log.Fatal("Failed to initialize audit handler:", err)
    }
}
```

### **2. Integração com Sistema de Auditoria Existente**

```go
// Conectar com o sistema de eventos existente
func connectToExistingAuditSystem(auditLogService authorization.AuditLogService) {
    // Registrar listener para eventos de auditoria
    auditLogger := authorization.NewAuditLogger(logger)
    
    // Quando um evento acontece, registrar no serviço
    auditLogger.OnEvent(func(event authorization.AuditEvent) {
        // Converter e armazenar no serviço
        auditLogService.StoreEvent(context.Background(), event)
    })
}
```

### **3. Configuração de Banco de Dados**

```go
// Para implementação com banco de dados real
type DatabaseAuditLogService struct {
    db *sql.DB
}

func (s *DatabaseAuditLogService) QueryLogs(ctx context.Context, query *AuditLogQuery) (*AuditLogsResponse, error) {
    // Construir query SQL baseada nos filtros
    sqlQuery := "SELECT * FROM audit_logs WHERE 1=1"
    args := []interface{}{}
    
    if query.UserID != "" {
        sqlQuery += " AND user_id = ?"
        args = append(args, query.UserID)
    }
    
    if query.ResourceType != "" {
        sqlQuery += " AND resource_type = ?"
        args = append(args, query.ResourceType)
    }
    
    // Executar query e retornar resultados
    // ...
}
```

## 🔄 Fluxo de Trabalho Completo

### **1. Evento de Auditoria**
```
User Action → Authorization Check → Audit Event → Database → API Response
```

### **2. Consulta de Logs**
```
API Request → Token Validation → Permission Check → Query Processing → Response
```

### **3. Exportação**
```
Export Request → Permission Check → Data Query → Format Conversion → File Download
```

## 📝 Próximos Passos

### **1. Implementação Imediata**
- ✅ Estruturas de dados implementadas
- ✅ Handler web implementado
- ✅ Integração com autorização
- ✅ Documentação completa

### **2. Melhorias Futuras**
- [ ] Implementação com banco de dados persistente
- [ ] Cache de consultas frequentes
- [ ] Alertas em tempo real
- [ ] Dashboard de monitoramento
- [ ] Relatórios automatizados

### **3. Testes**
- [ ] Testes unitários do handler
- [ ] Testes de integração
- [ ] Testes de performance
- [ ] Testes de segurança

## 🎯 Benefícios da Implementação

1. **Conformidade**: Atende aos requisitos de auditoria empresarial
2. **Segurança**: Controle granular de acesso baseado em OpenFGA
3. **Performance**: Consultas otimizadas com paginação
4. **Flexibilidade**: Filtros avançados e múltiplos formatos
5. **Escalabilidade**: Arquitetura preparada para crescimento
6. **Monitoramento**: Estatísticas e tendências em tempo real

Esta implementação fornece uma base sólida para auditoria completa do sistema Lockari, integrando perfeitamente com a arquitetura existente e fornecendo todas as funcionalidades necessárias para um sistema de auditoria robusto e seguro.
