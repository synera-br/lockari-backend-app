package authorization

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

// ===== AUDIT LOG SERVICE IMPLEMENTATION =====

// auditLogService implementa AuditLogService
type auditLogService struct {
	// Storage para eventos de auditoria
	events []AuditEvent
	mu     sync.RWMutex

	// Configurações
	maxEvents int

	// Cache para usuários (email por userID)
	userCache map[string]string
	userMu    sync.RWMutex

	// Métricas
	stats *AuditLogStats
}

// NewAuditLogService cria uma nova instância do serviço de logs de auditoria
func NewAuditLogService(maxEvents int) AuditLogService {
	return &auditLogService{
		events:    make([]AuditEvent, 0),
		maxEvents: maxEvents,
		userCache: make(map[string]string),
		stats:     &AuditLogStats{},
	}
}

// QueryLogs busca logs de auditoria baseado nos filtros
func (s *auditLogService) QueryLogs(ctx context.Context, query *AuditLogQuery) (*AuditLogsResponse, error) {
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	query.SetDefaults()

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Filtrar eventos
	filteredEvents := s.filterEvents(query)

	// Ordenar
	s.sortEvents(filteredEvents, query.SortBy, query.SortOrder)

	// Calcular paginação
	totalItems := len(filteredEvents)
	totalPages := int(math.Ceil(float64(totalItems) / float64(query.Limit)))
	offset := query.CalculateOffset()

	// Paginar resultados
	end := offset + query.Limit
	if end > totalItems {
		end = totalItems
	}

	if offset > totalItems {
		offset = totalItems
	}

	pagedEvents := filteredEvents[offset:end]

	// Converter para AuditLogData
	logs := make([]AuditLogData, len(pagedEvents))
	for i, event := range pagedEvents {
		userEmail := s.getUserEmail(event.UserID)
		logs[i] = *ToAuditLogData(&event, userEmail)
	}

	return &AuditLogsResponse{
		Logs: logs,
		Pagination: PaginationData{
			CurrentPage: query.Page,
			TotalPages:  totalPages,
			TotalItems:  totalItems,
			Limit:       query.Limit,
		},
	}, nil
}

// GetLogByID busca um log específico por ID
func (s *auditLogService) GetLogByID(ctx context.Context, logID string) (*AuditLogData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.ID == logID {
			userEmail := s.getUserEmail(event.UserID)
			return ToAuditLogData(&event, userEmail), nil
		}
	}

	return nil, fmt.Errorf("log not found: %s", logID)
}

// GetUserActivity busca atividades de um usuário específico
func (s *auditLogService) GetUserActivity(ctx context.Context, userID string, limit int) (*AuditLogsResponse, error) {
	query := &AuditLogQuery{
		UserID:    userID,
		Limit:     limit,
		Page:      1,
		SortBy:    "timestamp",
		SortOrder: "desc",
	}

	return s.QueryLogs(ctx, query)
}

// GetResourceActivity busca atividades de um recurso específico
func (s *auditLogService) GetResourceActivity(ctx context.Context, resourceType, resourceID string, limit int) (*AuditLogsResponse, error) {
	query := &AuditLogQuery{
		ResourceType: resourceType,
		ResourceName: resourceID,
		Limit:        limit,
		Page:         1,
		SortBy:       "timestamp",
		SortOrder:    "desc",
	}

	return s.QueryLogs(ctx, query)
}

// GetSuspiciousActivity busca atividades suspeitas
func (s *auditLogService) GetSuspiciousActivity(ctx context.Context, limit int) (*AuditLogsResponse, error) {
	query := &AuditLogQuery{
		Action:    "suspicious_activity",
		Limit:     limit,
		Page:      1,
		SortBy:    "timestamp",
		SortOrder: "desc",
	}

	return s.QueryLogs(ctx, query)
}

// ExportLogs exporta logs para arquivo
func (s *auditLogService) ExportLogs(ctx context.Context, query *AuditLogQuery, format string) ([]byte, error) {
	// Buscar todos os logs sem paginação
	originalLimit := query.Limit
	query.Limit = 10000 // limite máximo para exportação

	response, err := s.QueryLogs(ctx, query)
	if err != nil {
		return nil, err
	}

	query.Limit = originalLimit

	switch strings.ToLower(format) {
	case "json":
		return json.MarshalIndent(response.Logs, "", "  ")
	case "csv":
		return s.exportToCSV(response.Logs), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ===== INTERNAL METHODS =====

// filterEvents filtra eventos baseado na query
func (s *auditLogService) filterEvents(query *AuditLogQuery) []AuditEvent {
	var filtered []AuditEvent

	for _, event := range s.events {
		if !s.matchesFilter(event, query) {
			continue
		}
		filtered = append(filtered, event)
	}

	return filtered
}

// matchesFilter verifica se um evento atende aos filtros
func (s *auditLogService) matchesFilter(event AuditEvent, query *AuditLogQuery) bool {
	// Filtro por usuário
	if query.UserID != "" && event.UserID != query.UserID {
		return false
	}

	// Filtro por email (busca no cache)
	if query.UserEmail != "" {
		userEmail := s.getUserEmail(event.UserID)
		if !strings.Contains(strings.ToLower(userEmail), strings.ToLower(query.UserEmail)) {
			return false
		}
	}

	// Filtro por tipo de recurso
	if query.ResourceType != "" {
		resourceType, _ := parseResource(event.Resource)
		if resourceType != query.ResourceType {
			return false
		}
	}

	// Filtro por nome do recurso
	if query.ResourceName != "" {
		_, resourceName := parseResource(event.Resource)
		if !strings.Contains(strings.ToLower(resourceName), strings.ToLower(query.ResourceName)) {
			return false
		}
	}

	// Filtro por ação
	if query.Action != "" && event.Action != query.Action {
		return false
	}

	// Filtro por data inicial
	if query.StartTime != nil && event.Timestamp.Before(*query.StartTime) {
		return false
	}

	// Filtro por data final
	if query.EndTime != nil && event.Timestamp.After(*query.EndTime) {
		return false
	}

	// Filtro por IP
	if query.IPAddress != "" && event.ClientIP != query.IPAddress {
		return false
	}

	// Filtro por sucesso
	if query.Success != nil {
		isSuccess := (event.Result == "allowed" || event.Result == "success")
		if *query.Success != isSuccess {
			return false
		}
	}

	return true
}

// sortEvents ordena eventos
func (s *auditLogService) sortEvents(events []AuditEvent, sortBy, sortOrder string) {
	sort.Slice(events, func(i, j int) bool {
		var result bool

		switch sortBy {
		case "timestamp":
			result = events[i].Timestamp.Before(events[j].Timestamp)
		case "action":
			result = events[i].Action < events[j].Action
		case "resourceType":
			typeI, _ := parseResource(events[i].Resource)
			typeJ, _ := parseResource(events[j].Resource)
			result = typeI < typeJ
		case "userId":
			result = events[i].UserID < events[j].UserID
		case "result":
			result = events[i].Result < events[j].Result
		default:
			result = events[i].Timestamp.Before(events[j].Timestamp)
		}

		if sortOrder == "desc" {
			return !result
		}
		return result
	})
}

// getUserEmail obtém o email do usuário (mock implementation)
func (s *auditLogService) getUserEmail(userID string) string {
	s.userMu.RLock()
	defer s.userMu.RUnlock()

	if email, exists := s.userCache[userID]; exists {
		return email
	}

	// Mock: gerar email baseado no userID
	email := fmt.Sprintf("%s@example.com", userID)

	s.userMu.Lock()
	s.userCache[userID] = email
	s.userMu.Unlock()

	return email
}

// exportToCSV exporta logs para formato CSV
func (s *auditLogService) exportToCSV(logs []AuditLogData) []byte {
	var buffer strings.Builder
	writer := csv.NewWriter(&buffer)

	// Cabeçalho
	headers := []string{
		"ID", "Timestamp", "User Email", "User ID", "Action",
		"Resource Type", "Resource Name", "IP Address", "Result", "Details",
	}
	writer.Write(headers)

	// Dados
	for _, log := range logs {
		detailsJSON, _ := json.Marshal(log.Details)
		record := []string{
			log.ID,
			log.Timestamp.Format(time.RFC3339),
			log.UserEmail,
			log.UserID,
			log.Action,
			log.ResourceType,
			log.ResourceName,
			log.IPAddress,
			extractResult(log.Details),
			string(detailsJSON),
		}
		writer.Write(record)
	}

	writer.Flush()
	return []byte(buffer.String())
}

// extractResult extrai o resultado dos detalhes
func extractResult(details map[string]interface{}) string {
	if result, exists := details["result"]; exists {
		if resultStr, ok := result.(string); ok {
			return resultStr
		}
	}
	return "unknown"
}

// ===== MOCK DATA FOR TESTING =====

// AddMockEvent adiciona um evento mock para testes
func (s *auditLogService) AddMockEvent(event AuditEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)

	// Limitar número de eventos
	if len(s.events) > s.maxEvents {
		s.events = s.events[1:]
	}
}

// GenerateMockData gera dados mock para testes
func (s *auditLogService) GenerateMockData() {
	mockEvents := []AuditEvent{
		{
			ID:        "audit-001",
			Timestamp: time.Now().Add(-24 * time.Hour),
			UserID:    "alice",
			Action:    "check_permission",
			Resource:  "vault:personal-secrets",
			Result:    "allowed",
			ClientIP:  "192.168.1.100",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-001",
			Duration:  50 * time.Millisecond,
			Metadata: map[string]interface{}{
				"permission": "can_read",
				"tenant":     "company-acme",
			},
		},
		{
			ID:        "audit-002",
			Timestamp: time.Now().Add(-23 * time.Hour),
			UserID:    "bob",
			Action:    "create_vault",
			Resource:  "vault:team-credentials",
			Result:    "allowed",
			ClientIP:  "192.168.1.101",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-002",
			Duration:  200 * time.Millisecond,
			Metadata: map[string]interface{}{
				"vault_name": "Team Credentials",
				"tenant":     "company-acme",
			},
		},
		{
			ID:        "audit-003",
			Timestamp: time.Now().Add(-22 * time.Hour),
			UserID:    "charlie",
			Action:    "read_secret",
			Resource:  "secret:database-password",
			Result:    "denied",
			ClientIP:  "192.168.1.102",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-003",
			Duration:  30 * time.Millisecond,
			Metadata: map[string]interface{}{
				"secret_type": "password",
				"vault_id":    "vault:team-credentials",
				"reason":      "insufficient_permissions",
			},
		},
		{
			ID:        "audit-004",
			Timestamp: time.Now().Add(-21 * time.Hour),
			UserID:    "alice",
			Action:    "share_vault",
			Resource:  "vault:personal-secrets",
			Result:    "allowed",
			ClientIP:  "192.168.1.100",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-004",
			Duration:  100 * time.Millisecond,
			Metadata: map[string]interface{}{
				"target_user": "bob",
				"permission":  "can_read",
			},
		},
		{
			ID:        "audit-005",
			Timestamp: time.Now().Add(-20 * time.Hour),
			UserID:    "david",
			Action:    "suspicious_activity",
			Resource:  "system:lockari",
			Result:    "flagged",
			ClientIP:  "203.0.113.1",
			UserAgent: "curl/7.68.0",
			RequestID: "req-005",
			Duration:  10 * time.Millisecond,
			Metadata: map[string]interface{}{
				"activity":      "multiple_failed_logins",
				"severity":      "high",
				"attempt_count": 5,
				"reason":        "too_many_failed_attempts",
			},
		},
		{
			ID:        "audit-006",
			Timestamp: time.Now().Add(-19 * time.Hour),
			UserID:    "eve",
			Action:    "create_token",
			Resource:  "token:api-automation",
			Result:    "allowed",
			ClientIP:  "192.168.1.103",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-006",
			Duration:  80 * time.Millisecond,
			Metadata: map[string]interface{}{
				"token_name":  "API Automation",
				"permissions": []string{"can_read_secrets", "can_write_secrets"},
				"vault_id":    "vault:team-credentials",
			},
		},
		{
			ID:        "audit-007",
			Timestamp: time.Now().Add(-18 * time.Hour),
			UserID:    "alice",
			Action:    "external_share_request",
			Resource:  "vault:personal-secrets",
			Result:    "pending",
			ClientIP:  "192.168.1.100",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-007",
			Duration:  150 * time.Millisecond,
			Metadata: map[string]interface{}{
				"target_tenant": "company-beta",
				"target_user":   "frank",
				"permission":    "can_read",
			},
		},
		{
			ID:        "audit-008",
			Timestamp: time.Now().Add(-17 * time.Hour),
			UserID:    "bob",
			Action:    "update_secret",
			Resource:  "secret:api-key",
			Result:    "allowed",
			ClientIP:  "192.168.1.101",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-008",
			Duration:  75 * time.Millisecond,
			Metadata: map[string]interface{}{
				"secret_type": "api_key",
				"vault_id":    "vault:team-credentials",
				"version":     "2",
			},
		},
		{
			ID:        "audit-009",
			Timestamp: time.Now().Add(-16 * time.Hour),
			UserID:    "charlie",
			Action:    "copy_secret",
			Resource:  "secret:database-password",
			Result:    "denied",
			ClientIP:  "192.168.1.102",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-009",
			Duration:  25 * time.Millisecond,
			Metadata: map[string]interface{}{
				"secret_type": "password",
				"vault_id":    "vault:team-credentials",
				"reason":      "insufficient_permissions",
			},
		},
		{
			ID:        "audit-010",
			Timestamp: time.Now().Add(-15 * time.Hour),
			UserID:    "alice",
			Action:    "grant_permission",
			Resource:  "vault:personal-secrets",
			Result:    "allowed",
			ClientIP:  "192.168.1.100",
			UserAgent: "Mozilla/5.0",
			RequestID: "req-010",
			Duration:  60 * time.Millisecond,
			Metadata: map[string]interface{}{
				"grantee":    "charlie",
				"permission": "can_read",
			},
		},
	}

	for _, event := range mockEvents {
		s.AddMockEvent(event)
	}
}

// ===== ADDITIONAL HELPER METHODS =====

// GetStats retorna estatísticas dos logs de auditoria
func (s *auditLogService) GetStats(ctx context.Context) (*AuditLogStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &AuditLogStats{
		TotalEvents:          int64(len(s.events)),
		EventsByAction:       make(map[string]int64),
		EventsByResourceType: make(map[string]int64),
		EventsByUser:         make(map[string]int64),
		EventsByDay:          make(map[string]int64),
		UniqueUsers:          0,
		UniqueIPs:            0,
	}

	uniqueUsers := make(map[string]bool)
	uniqueIPs := make(map[string]bool)
	successCount := int64(0)
	suspiciousCount := int64(0)

	for _, event := range s.events {
		// Contadores por ação
		stats.EventsByAction[event.Action]++

		// Contadores por tipo de recurso
		resourceType, _ := parseResource(event.Resource)
		stats.EventsByResourceType[resourceType]++

		// Contadores por usuário
		stats.EventsByUser[event.UserID]++

		// Contadores por dia
		day := event.Timestamp.Format("2006-01-02")
		stats.EventsByDay[day]++

		// Usuários únicos
		uniqueUsers[event.UserID] = true

		// IPs únicos
		if event.ClientIP != "" {
			uniqueIPs[event.ClientIP] = true
		}

		// Taxa de sucesso
		if event.Result == "allowed" || event.Result == "success" {
			successCount++
		}

		// Eventos suspeitos
		if event.Action == "suspicious_activity" {
			suspiciousCount++
		}

		// Última atividade
		if event.Timestamp.After(stats.LastActivity) {
			stats.LastActivity = event.Timestamp
		}
	}

	stats.UniqueUsers = int64(len(uniqueUsers))
	stats.UniqueIPs = int64(len(uniqueIPs))
	stats.SuspiciousEvents = suspiciousCount

	if stats.TotalEvents > 0 {
		stats.SuccessRate = float64(successCount) / float64(stats.TotalEvents) * 100
	}

	return stats, nil
}

// GetTrends retorna tendências de auditoria por período
func (s *auditLogService) GetTrends(ctx context.Context, days int) ([]AuditLogTrend, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	trends := make(map[string]*AuditLogTrend)

	// Inicializar dias
	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		trends[dateStr] = &AuditLogTrend{
			Date:        date,
			EventCount:  0,
			UserCount:   0,
			ErrorCount:  0,
			SuccessRate: 0,
		}
	}

	// Processar eventos
	for _, event := range s.events {
		dateStr := event.Timestamp.Format("2006-01-02")
		if trend, exists := trends[dateStr]; exists {
			trend.EventCount++

			if event.Result == "denied" || event.Result == "error" {
				trend.ErrorCount++
			}
		}
	}

	// Calcular taxa de sucesso
	for _, trend := range trends {
		if trend.EventCount > 0 {
			successCount := trend.EventCount - trend.ErrorCount
			trend.SuccessRate = float64(successCount) / float64(trend.EventCount) * 100
		}
	}

	// Converter para slice ordenado
	var result []AuditLogTrend
	for _, trend := range trends {
		result = append(result, *trend)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result, nil
}

// SetUserEmail define o email de um usuário no cache
func (s *auditLogService) SetUserEmail(userID, email string) {
	s.userMu.Lock()
	defer s.userMu.Unlock()
	s.userCache[userID] = email
}
