package authorization

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ===== AUDIT API STRUCTURES =====

// AuditLogsResponse é a estrutura principal para a resposta da API de logs de auditoria
// Esta é a estrutura que deve ser marshalled para JSON e enviada ao frontend
type AuditLogsResponse struct {
	Logs       []AuditLogData `json:"logs"`
	Pagination PaginationData `json:"pagination"`
}

// AuditLogData representa um evento de auditoria individual
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

// PaginationData contém os metadados de paginação para a lista de logs
type PaginationData struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
	TotalItems  int `json:"totalItems"`
	Limit       int `json:"limit"`
}

// AuditQuery representa uma consulta para buscar logs de auditoria
type AuditLogQuery struct {
	// Filtros básicos
	UserID       string     `json:"userId,omitempty"`
	UserEmail    string     `json:"userEmail,omitempty"`
	ResourceType string     `json:"resourceType,omitempty"`
	ResourceName string     `json:"resourceName,omitempty"`
	Action       string     `json:"action,omitempty"`
	
	// Filtros de data
	StartTime *time.Time `json:"startTime,omitempty"`
	EndTime   *time.Time `json:"endTime,omitempty"`
	
	// Paginação
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
	
	// Ordenação
	SortBy    string `json:"sortBy,omitempty"`    // timestamp, action, resourceType, etc.
	SortOrder string `json:"sortOrder,omitempty"` // asc, desc
	
	// Filtros avançados
	IPAddress string `json:"ipAddress,omitempty"`
	Success   *bool  `json:"success,omitempty"` // true para permitido, false para negado
}

// AuditLogService interface para operações de auditoria
type AuditLogService interface {
	// QueryLogs busca logs de auditoria baseado nos filtros
	QueryLogs(ctx context.Context, query *AuditLogQuery) (*AuditLogsResponse, error)
	
	// GetLogByID busca um log específico por ID
	GetLogByID(ctx context.Context, logID string) (*AuditLogData, error)
	
	// GetUserActivity busca atividades de um usuário específico
	GetUserActivity(ctx context.Context, userID string, limit int) (*AuditLogsResponse, error)
	
	// GetResourceActivity busca atividades de um recurso específico
	GetResourceActivity(ctx context.Context, resourceType, resourceID string, limit int) (*AuditLogsResponse, error)
	
	// GetSuspiciousActivity busca atividades suspeitas
	GetSuspiciousActivity(ctx context.Context, limit int) (*AuditLogsResponse, error)
	
	// ExportLogs exporta logs para arquivo (CSV, JSON)
	ExportLogs(ctx context.Context, query *AuditLogQuery, format string) ([]byte, error)
}

// ===== AUDIT LOG MAPPERS =====

// ToAuditLogData converte AuditEvent para AuditLogData (formato da API)
func ToAuditLogData(event *AuditEvent, userEmail string) *AuditLogData {
	resourceType, resourceName := parseResource(event.Resource)
	
	return &AuditLogData{
		ID:           event.ID,
		ResourceName: resourceName,
		ResourceType: resourceType,
		Action:       formatAction(event.Action),
		UserEmail:    userEmail,
		UserID:       event.UserID,
		IPAddress:    event.ClientIP,
		Timestamp:    event.Timestamp,
		ResourceLink: generateResourceLink(resourceType, resourceName),
		Details:      buildDetails(event),
	}
}

// parseResource extrai tipo e nome do recurso a partir do formato "type:name"
func parseResource(resource string) (resourceType, resourceName string) {
	parts := strings.Split(resource, ":")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "unknown", resource
}

// formatAction formata a ação para exibição amigável
func formatAction(action string) string {
	actionMap := map[string]string{
		"check_permission":         "Check Permission",
		"grant_permission":         "Grant Permission",
		"revoke_permission":        "Revoke Permission",
		"create_vault":             "Create Vault",
		"update_vault":             "Update Vault",
		"delete_vault":             "Delete Vault",
		"share_vault":              "Share Vault",
		"create_secret":            "Create Secret",
		"update_secret":            "Update Secret",
		"delete_secret":            "Delete Secret",
		"read_secret":              "Read Secret",
		"copy_secret":              "Copy Secret",
		"download_secret":          "Download Secret",
		"create_token":             "Create API Token",
		"revoke_token":             "Revoke API Token",
		"login":                    "Login",
		"logout":                   "Logout",
		"password_change":          "Password Change",
		"profile_update":           "Profile Update",
		"tenant_setup":             "Tenant Setup",
		"user_invitation":          "User Invitation",
		"group_create":             "Create Group",
		"group_update":             "Update Group",
		"group_delete":             "Delete Group",
		"external_share_request":   "External Share Request",
		"external_share_approve":   "External Share Approve",
		"external_share_reject":    "External Share Reject",
		"suspicious_activity":      "Suspicious Activity",
		"failed_login":             "Failed Login",
		"multiple_failed_logins":   "Multiple Failed Logins",
		"account_lockout":          "Account Lockout",
		"permission_escalation":    "Permission Escalation",
		"bulk_operation":           "Bulk Operation",
		"system_configuration":     "System Configuration",
	}
	
	if formatted, exists := actionMap[action]; exists {
		return formatted
	}
	
	// Se não encontrar, formatar automaticamente
	return strings.Title(strings.ReplaceAll(action, "_", " "))
}

// generateResourceLink gera um link para o recurso no frontend
func generateResourceLink(resourceType, resourceName string) string {
	if resourceName == "" {
		return ""
	}
	
	switch resourceType {
	case "vault":
		return fmt.Sprintf("/vaults/%s", resourceName)
	case "secret":
		return fmt.Sprintf("/secrets/%s", resourceName)
	case "user":
		return fmt.Sprintf("/users/%s", resourceName)
	case "group":
		return fmt.Sprintf("/groups/%s", resourceName)
	case "tenant":
		return fmt.Sprintf("/tenants/%s", resourceName)
	case "token":
		return fmt.Sprintf("/tokens/%s", resourceName)
	default:
		return ""
	}
}

// buildDetails constrói os detalhes específicos do evento
func buildDetails(event *AuditEvent) map[string]interface{} {
	details := make(map[string]interface{})
	
	// Adicionar metadados do evento
	for key, value := range event.Metadata {
		details[key] = value
	}
	
	// Adicionar informações específicas
	if event.Result != "" {
		details["result"] = event.Result
	}
	
	if event.Duration > 0 {
		details["duration_ms"] = event.Duration.Milliseconds()
	}
	
	if event.RequestID != "" {
		details["request_id"] = event.RequestID
	}
	
	if event.UserAgent != "" {
		details["user_agent"] = event.UserAgent
	}
	
	// Adicionar contexto específico baseado na ação
	switch event.Action {
	case "check_permission":
		if permission, exists := event.Metadata["permission"]; exists {
			details["permission"] = permission
		}
	case "grant_permission", "revoke_permission":
		if grantee, exists := event.Metadata["grantee"]; exists {
			details["target_user"] = grantee
		}
	case "create_vault", "update_vault":
		if name, exists := event.Metadata["vault_name"]; exists {
			details["vault_name"] = name
		}
	case "share_vault":
		if targetUser, exists := event.Metadata["target_user"]; exists {
			details["shared_with"] = targetUser
		}
	case "create_secret", "update_secret":
		if secretType, exists := event.Metadata["secret_type"]; exists {
			details["secret_type"] = secretType
		}
	case "create_token":
		if permissions, exists := event.Metadata["permissions"]; exists {
			details["token_permissions"] = permissions
		}
	case "suspicious_activity":
		if severity, exists := event.Metadata["severity"]; exists {
			details["severity"] = severity
		}
		if reason, exists := event.Metadata["reason"]; exists {
			details["reason"] = reason
		}
	}
	
	return details
}

// ===== AUDIT LOG FILTERS =====

// AuditLogFilter representa filtros para busca de logs
type AuditLogFilter struct {
	UserIDs      []string   `json:"userIds,omitempty"`
	ResourceTypes []string   `json:"resourceTypes,omitempty"`
	Actions      []string   `json:"actions,omitempty"`
	StartTime    *time.Time `json:"startTime,omitempty"`
	EndTime      *time.Time `json:"endTime,omitempty"`
	Success      *bool      `json:"success,omitempty"`
	IPAddresses  []string   `json:"ipAddresses,omitempty"`
	Severity     []string   `json:"severity,omitempty"`
}

// ===== AUDIT LOG AGGREGATIONS =====

// AuditLogStats representa estatísticas de auditoria
type AuditLogStats struct {
	TotalEvents          int64                    `json:"totalEvents"`
	EventsByAction       map[string]int64         `json:"eventsByAction"`
	EventsByResourceType map[string]int64         `json:"eventsByResourceType"`
	EventsByUser         map[string]int64         `json:"eventsByUser"`
	EventsByDay          map[string]int64         `json:"eventsByDay"`
	SuccessRate          float64                  `json:"successRate"`
	SuspiciousEvents     int64                    `json:"suspiciousEvents"`
	UniqueUsers          int64                    `json:"uniqueUsers"`
	UniqueIPs            int64                    `json:"uniqueIPs"`
	LastActivity         time.Time                `json:"lastActivity"`
}

// AuditLogTrend representa tendências de auditoria
type AuditLogTrend struct {
	Date        time.Time `json:"date"`
	EventCount  int64     `json:"eventCount"`
	UserCount   int64     `json:"userCount"`
	ErrorCount  int64     `json:"errorCount"`
	SuccessRate float64   `json:"successRate"`
}

// ===== AUDIT LOG ALERTS =====

// AuditAlert representa um alerta baseado em logs de auditoria
type AuditAlert struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	UserID      string                 `json:"userId"`
	ResourceID  string                 `json:"resourceId"`
	TriggerTime time.Time              `json:"triggerTime"`
	Status      string                 `json:"status"`
	Details     map[string]interface{} `json:"details"`
}

// AuditAlertType representa tipos de alertas
type AuditAlertType string

const (
	AuditAlertTypeMultipleFailedLogins    AuditAlertType = "multiple_failed_logins"
	AuditAlertTypeUnusualAccess          AuditAlertType = "unusual_access"
	AuditAlertTypePermissionEscalation   AuditAlertType = "permission_escalation"
	AuditAlertTypeBulkOperations         AuditAlertType = "bulk_operations"
	AuditAlertTypeOffHoursAccess         AuditAlertType = "off_hours_access"
	AuditAlertTypeGeographicAnomaly      AuditAlertType = "geographic_anomaly"
	AuditAlertTypeRapidOperations        AuditAlertType = "rapid_operations"
	AuditAlertTypeSystemConfigChange     AuditAlertType = "system_config_change"
	AuditAlertTypeDataExfiltration       AuditAlertType = "data_exfiltration"
	AuditAlertTypeUnauthorizedAccess     AuditAlertType = "unauthorized_access"
)

// AuditAlertSeverity representa severidades de alertas
type AuditAlertSeverity string

const (
	AuditAlertSeverityLow      AuditAlertSeverity = "low"
	AuditAlertSeverityMedium   AuditAlertSeverity = "medium"
	AuditAlertSeverityHigh     AuditAlertSeverity = "high"
	AuditAlertSeverityCritical AuditAlertSeverity = "critical"
)

// ===== AUDIT LOG EXPORT =====

// AuditLogExportRequest representa uma solicitação de exportação
type AuditLogExportRequest struct {
	Query     *AuditLogQuery `json:"query"`
	Format    string         `json:"format"`    // csv, json, xlsx
	Fields    []string       `json:"fields"`    // campos específicos para exportar
	Compress  bool           `json:"compress"`  // gzip compress
	EmailTo   string         `json:"emailTo"`   // enviar por email
	FileName  string         `json:"fileName"`  // nome do arquivo
}

// AuditLogExportResponse representa a resposta da exportação
type AuditLogExportResponse struct {
	FileID      string    `json:"fileId"`
	FileName    string    `json:"fileName"`
	FileSize    int64     `json:"fileSize"`
	RecordCount int64     `json:"recordCount"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	DownloadURL string    `json:"downloadUrl"`
}

// ===== VALIDATION FUNCTIONS =====

// Validate valida AuditLogQuery
func (q *AuditLogQuery) Validate() error {
	if q.Limit > 1000 {
		return fmt.Errorf("limit cannot exceed 1000")
	}
	
	if q.Limit <= 0 {
		q.Limit = 50 // default
	}
	
	if q.Page <= 0 {
		q.Page = 1 // default
	}
	
	if q.StartTime != nil && q.EndTime != nil {
		if q.StartTime.After(*q.EndTime) {
			return fmt.Errorf("start time cannot be after end time")
		}
	}
	
	// Validar ordenação
	if q.SortBy != "" {
		validSortFields := []string{"timestamp", "action", "resourceType", "userId", "result"}
		found := false
		for _, field := range validSortFields {
			if q.SortBy == field {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid sort field: %s", q.SortBy)
		}
	}
	
	if q.SortOrder != "" && q.SortOrder != "asc" && q.SortOrder != "desc" {
		return fmt.Errorf("sort order must be 'asc' or 'desc'")
	}
	
	return nil
}

// SetDefaults define valores padrão para AuditLogQuery
func (q *AuditLogQuery) SetDefaults() {
	if q.Limit <= 0 {
		q.Limit = 50
	}
	
	if q.Page <= 0 {
		q.Page = 1
	}
	
	if q.SortBy == "" {
		q.SortBy = "timestamp"
	}
	
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
}

// ===== HELPER FUNCTIONS =====

// CalculateOffset calcula o offset para paginação
func (q *AuditLogQuery) CalculateOffset() int {
	return (q.Page - 1) * q.Limit
}

// BuildWhereClause constrói cláusula WHERE para SQL (exemplo)
func (q *AuditLogQuery) BuildWhereClause() (string, []interface{}) {
	var conditions []string
	var args []interface{}
	
	if q.UserID != "" {
		conditions = append(conditions, "user_id = ?")
		args = append(args, q.UserID)
	}
	
	if q.ResourceType != "" {
		conditions = append(conditions, "resource_type = ?")
		args = append(args, q.ResourceType)
	}
	
	if q.Action != "" {
		conditions = append(conditions, "action = ?")
		args = append(args, q.Action)
	}
	
	if q.StartTime != nil {
		conditions = append(conditions, "timestamp >= ?")
		args = append(args, q.StartTime)
	}
	
	if q.EndTime != nil {
		conditions = append(conditions, "timestamp <= ?")
		args = append(args, q.EndTime)
	}
	
	if q.IPAddress != "" {
		conditions = append(conditions, "client_ip = ?")
		args = append(args, q.IPAddress)
	}
	
	if q.Success != nil {
		if *q.Success {
			conditions = append(conditions, "result = 'allowed'")
		} else {
			conditions = append(conditions, "result = 'denied'")
		}
	}
	
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	
	return whereClause, args
}

// IsValidResourceType verifica se o tipo de recurso é válido
func IsValidResourceType(resourceType string) bool {
	validTypes := []string{"vault", "secret", "user", "group", "tenant", "token", "system"}
	for _, validType := range validTypes {
		if resourceType == validType {
			return true
		}
	}
	return false
}

// IsValidAction verifica se a ação é válida
func IsValidAction(action string) bool {
	validActions := []string{
		"check_permission", "grant_permission", "revoke_permission",
		"create_vault", "update_vault", "delete_vault", "share_vault",
		"create_secret", "update_secret", "delete_secret", "read_secret",
		"copy_secret", "download_secret", "create_token", "revoke_token",
		"login", "logout", "password_change", "profile_update",
		"tenant_setup", "user_invitation", "group_create", "group_update",
		"group_delete", "external_share_request", "external_share_approve",
		"external_share_reject", "suspicious_activity", "failed_login",
	}
	
	for _, validAction := range validActions {
		if action == validAction {
			return true
		}
	}
	return false
}
