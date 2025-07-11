package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// AuditLogger é uma implementação de auditoria que usa logger
type AuditLogger struct {
	logger Logger
	mu     sync.RWMutex
}

// NewAuditLogger cria um novo audit logger
func NewAuditLogger(logger Logger) *AuditLogger {
	return &AuditLogger{
		logger: logger,
	}
}

// LogPermissionCheck registra uma verificação de permissão
func (al *AuditLogger) LogPermissionCheck(ctx context.Context, event PermissionCheckEvent) {
	al.mu.RLock()
	defer al.mu.RUnlock()

	if al.logger == nil {
		return
	}

	logData := map[string]interface{}{
		"event_type":  "permission_check",
		"user":        event.User,
		"relation":    event.Relation,
		"object":      event.Object,
		"result":      event.Result,
		"timestamp":   event.Timestamp,
		"duration_ms": event.Duration.Milliseconds(),
	}

	if event.Error != "" {
		logData["error"] = event.Error
	}

	// Extrair informações do contexto se disponíveis
	if requestID := ctx.Value("request_id"); requestID != nil {
		logData["request_id"] = requestID
	}

	if userAgent := ctx.Value("user_agent"); userAgent != nil {
		logData["user_agent"] = userAgent
	}

	if clientIP := ctx.Value("client_ip"); clientIP != nil {
		logData["client_ip"] = clientIP
	}

	al.logger.Info("Permission check", logData)
}

// LogPermissionGrant registra uma concessão de permissão
func (al *AuditLogger) LogPermissionGrant(ctx context.Context, event PermissionGrantEvent) {
	al.mu.RLock()
	defer al.mu.RUnlock()

	if al.logger == nil {
		return
	}

	logData := map[string]interface{}{
		"event_type": "permission_grant",
		"grantor":    event.Grantor,
		"grantee":    event.Grantee,
		"relation":   event.Relation,
		"object":     event.Object,
		"timestamp":  event.Timestamp,
	}

	// Extrair informações do contexto se disponíveis
	if requestID := ctx.Value("request_id"); requestID != nil {
		logData["request_id"] = requestID
	}

	if userAgent := ctx.Value("user_agent"); userAgent != nil {
		logData["user_agent"] = userAgent
	}

	if clientIP := ctx.Value("client_ip"); clientIP != nil {
		logData["client_ip"] = clientIP
	}

	al.logger.Info("Permission granted", logData)
}

// LogPermissionRevoke registra uma revogação de permissão
func (al *AuditLogger) LogPermissionRevoke(ctx context.Context, event PermissionRevokeEvent) {
	al.mu.RLock()
	defer al.mu.RUnlock()

	if al.logger == nil {
		return
	}

	logData := map[string]interface{}{
		"event_type": "permission_revoke",
		"revoker":    event.Revoker,
		"revokee":    event.Revokee,
		"relation":   event.Relation,
		"object":     event.Object,
		"timestamp":  event.Timestamp,
	}

	// Extrair informações do contexto se disponíveis
	if requestID := ctx.Value("request_id"); requestID != nil {
		logData["request_id"] = requestID
	}

	if userAgent := ctx.Value("user_agent"); userAgent != nil {
		logData["user_agent"] = userAgent
	}

	if clientIP := ctx.Value("client_ip"); clientIP != nil {
		logData["client_ip"] = clientIP
	}

	al.logger.Warn("Permission revoked", logData)
}

// LogSuspiciousActivity registra atividade suspeita
func (al *AuditLogger) LogSuspiciousActivity(ctx context.Context, event SuspiciousActivityEvent) {
	al.mu.RLock()
	defer al.mu.RUnlock()

	if al.logger == nil {
		return
	}

	logData := map[string]interface{}{
		"event_type": "suspicious_activity",
		"user":       event.User,
		"activity":   event.Activity,
		"details":    event.Details,
		"severity":   event.Severity,
		"client_ip":  event.ClientIP,
		"user_agent": event.UserAgent,
		"timestamp":  event.Timestamp,
	}

	// Extrair informações do contexto se disponíveis
	if requestID := ctx.Value("request_id"); requestID != nil {
		logData["request_id"] = requestID
	}

	// Log com nível baseado na severidade
	switch event.Severity {
	case "critical":
		al.logger.Error("Suspicious activity detected", logData)
	case "high":
		al.logger.Warn("Suspicious activity detected", logData)
	default:
		al.logger.Info("Suspicious activity detected", logData)
	}
}

// LogGenericEvent registra um evento genérico
func (al *AuditLogger) LogGenericEvent(ctx context.Context, event AuditEvent) {
	al.mu.RLock()
	defer al.mu.RUnlock()

	if al.logger == nil {
		return
	}

	logData := map[string]interface{}{
		"event_type":  "generic_event",
		"id":          event.ID,
		"timestamp":   event.Timestamp,
		"user_id":     event.UserID,
		"action":      event.Action,
		"resource":    event.Resource,
		"result":      event.Result,
		"metadata":    event.Metadata,
		"request_id":  event.RequestID,
		"client_ip":   event.ClientIP,
		"user_agent":  event.UserAgent,
		"duration_ms": event.Duration.Milliseconds(),
	}

	al.logger.Info("Generic audit event", logData)
}

// GetLogs recupera logs de auditoria
func (al *AuditLogger) GetLogs(ctx context.Context, userID, resource string, limit int) ([]AuditEvent, error) {
	// Esta implementação não suporta recuperação de logs
	// Para uma implementação real, você precisaria de um armazenamento persistente
	return []AuditEvent{}, fmt.Errorf("log retrieval not supported by AuditLogger")
}

// Close fecha o serviço de auditoria gracefully
func (al *AuditLogger) Close() error {
	al.mu.Lock()
	defer al.mu.Unlock()

	// Nada para limpar nesta implementação simples
	return nil
}

// === AUDIT BUFFER ===

// AuditBuffer é um buffer para eventos de auditoria
type AuditBuffer struct {
	events    []AuditEvent
	mu        sync.RWMutex
	maxSize   int
	flushFunc func([]AuditEvent) error
}

// NewAuditBuffer cria um novo buffer de auditoria
func NewAuditBuffer(maxSize int, flushFunc func([]AuditEvent) error) *AuditBuffer {
	return &AuditBuffer{
		events:    make([]AuditEvent, 0, maxSize),
		maxSize:   maxSize,
		flushFunc: flushFunc,
	}
}

// Add adiciona um evento ao buffer
func (ab *AuditBuffer) Add(event AuditEvent) {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	ab.events = append(ab.events, event)

	// Flush se o buffer estiver cheio
	if len(ab.events) >= ab.maxSize {
		ab.flush()
	}
}

// Flush força o flush do buffer
func (ab *AuditBuffer) Flush() error {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	return ab.flush()
}

// flush executa o flush interno
func (ab *AuditBuffer) flush() error {
	if len(ab.events) == 0 {
		return nil
	}

	if ab.flushFunc != nil {
		if err := ab.flushFunc(ab.events); err != nil {
			return err
		}
	}

	ab.events = ab.events[:0] // Reset slice
	return nil
}

// === AUDIT STORE ===

// AuditStore é uma implementação de armazenamento de auditoria
type AuditStore struct {
	buffer *AuditBuffer
	logger Logger
}

// NewAuditStore cria um novo armazenamento de auditoria
func NewAuditStore(logger Logger) *AuditStore {
	store := &AuditStore{
		logger: logger,
	}

	// Criar buffer com função de flush
	store.buffer = NewAuditBuffer(100, store.flushToLogger)

	return store
}

// LogPermissionCheck registra uma verificação de permissão
func (as *AuditStore) LogPermissionCheck(ctx context.Context, event PermissionCheckEvent) {
	auditEvent := AuditEvent{
		ID:        generateEventID(),
		Timestamp: event.Timestamp,
		UserID:    extractUserID(event.User),
		Action:    "permission_check",
		Resource:  event.Object,
		Result:    event.Result,
		Metadata: map[string]interface{}{
			"relation": event.Relation,
			"duration": event.Duration.Milliseconds(),
		},
	}

	// Extrair informações do contexto
	if requestID := ctx.Value("request_id"); requestID != nil {
		auditEvent.RequestID = fmt.Sprintf("%v", requestID)
	}

	if userAgent := ctx.Value("user_agent"); userAgent != nil {
		auditEvent.UserAgent = fmt.Sprintf("%v", userAgent)
	}

	if clientIP := ctx.Value("client_ip"); clientIP != nil {
		auditEvent.ClientIP = fmt.Sprintf("%v", clientIP)
	}

	auditEvent.Duration = event.Duration

	as.buffer.Add(auditEvent)
}

// LogPermissionGrant registra uma concessão de permissão
func (as *AuditStore) LogPermissionGrant(ctx context.Context, event PermissionGrantEvent) {
	auditEvent := AuditEvent{
		ID:        generateEventID(),
		Timestamp: event.Timestamp,
		UserID:    extractUserID(event.Grantor),
		Action:    "permission_grant",
		Resource:  event.Object,
		Result:    "success",
		Metadata: map[string]interface{}{
			"grantee":  event.Grantee,
			"relation": event.Relation,
		},
	}

	// Extrair informações do contexto
	if requestID := ctx.Value("request_id"); requestID != nil {
		auditEvent.RequestID = fmt.Sprintf("%v", requestID)
	}

	if userAgent := ctx.Value("user_agent"); userAgent != nil {
		auditEvent.UserAgent = fmt.Sprintf("%v", userAgent)
	}

	if clientIP := ctx.Value("client_ip"); clientIP != nil {
		auditEvent.ClientIP = fmt.Sprintf("%v", clientIP)
	}

	as.buffer.Add(auditEvent)
}

// LogPermissionRevoke registra uma revogação de permissão
func (as *AuditStore) LogPermissionRevoke(ctx context.Context, event PermissionRevokeEvent) {
	auditEvent := AuditEvent{
		ID:        generateEventID(),
		Timestamp: event.Timestamp,
		UserID:    extractUserID(event.Revoker),
		Action:    "permission_revoke",
		Resource:  event.Object,
		Result:    "success",
		Metadata: map[string]interface{}{
			"revokee":  event.Revokee,
			"relation": event.Relation,
		},
	}

	// Extrair informações do contexto
	if requestID := ctx.Value("request_id"); requestID != nil {
		auditEvent.RequestID = fmt.Sprintf("%v", requestID)
	}

	if userAgent := ctx.Value("user_agent"); userAgent != nil {
		auditEvent.UserAgent = fmt.Sprintf("%v", userAgent)
	}

	if clientIP := ctx.Value("client_ip"); clientIP != nil {
		auditEvent.ClientIP = fmt.Sprintf("%v", clientIP)
	}

	as.buffer.Add(auditEvent)
}

// LogSuspiciousActivity registra atividade suspeita
func (as *AuditStore) LogSuspiciousActivity(ctx context.Context, event SuspiciousActivityEvent) {
	auditEvent := AuditEvent{
		ID:        generateEventID(),
		Timestamp: event.Timestamp,
		UserID:    extractUserID(event.User),
		Action:    "suspicious_activity",
		Resource:  "system",
		Result:    event.Severity,
		Metadata: map[string]interface{}{
			"activity": event.Activity,
			"details":  event.Details,
		},
		ClientIP:  event.ClientIP,
		UserAgent: event.UserAgent,
	}

	// Extrair informações do contexto
	if requestID := ctx.Value("request_id"); requestID != nil {
		auditEvent.RequestID = fmt.Sprintf("%v", requestID)
	}

	as.buffer.Add(auditEvent)
}

// flushToLogger envia eventos para o logger
func (as *AuditStore) flushToLogger(events []AuditEvent) error {
	if as.logger == nil {
		return nil
	}

	for _, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			continue
		}

		as.logger.Info("Audit event", map[string]interface{}{
			"audit_event": string(eventData),
		})
	}

	return nil
}

// === AUDIT ANALYZER ===

// AuditAnalyzer analisa eventos de auditoria para detectar padrões suspeitos
type AuditAnalyzer struct {
	logger Logger
	audit  AuditService
}

// NewAuditAnalyzer cria um novo analisador de auditoria
func NewAuditAnalyzer(logger Logger, audit AuditService) *AuditAnalyzer {
	return &AuditAnalyzer{
		logger: logger,
		audit:  audit,
	}
}

// AnalyzePermissionCheck analisa verificações de permissão
func (aa *AuditAnalyzer) AnalyzePermissionCheck(ctx context.Context, event PermissionCheckEvent) {
	// Detectar múltiplas tentativas de acesso negado
	if event.Result == "denied" {
		aa.detectRepeatedDenials(ctx, event)
	}

	// Detectar acessos incomuns
	aa.detectUnusualAccess(ctx, event)
}

// detectRepeatedDenials detecta tentativas repetidas de acesso negado
func (aa *AuditAnalyzer) detectRepeatedDenials(ctx context.Context, event PermissionCheckEvent) {
	// TODO: Implementar lógica para detectar tentativas repetidas
	// Por enquanto, apenas log
	if aa.logger != nil {
		aa.logger.Debug("Repeated denial detected", map[string]interface{}{
			"user":     event.User,
			"object":   event.Object,
			"relation": event.Relation,
		})
	}
}

// detectUnusualAccess detecta acessos incomuns
func (aa *AuditAnalyzer) detectUnusualAccess(ctx context.Context, event PermissionCheckEvent) {
	// TODO: Implementar lógica para detectar acessos incomuns
	// Por enquanto, apenas log
	if aa.logger != nil {
		aa.logger.Debug("Unusual access pattern", map[string]interface{}{
			"user":     event.User,
			"object":   event.Object,
			"relation": event.Relation,
		})
	}
}

// === UTILITY FUNCTIONS ===

// generateEventID gera um ID único para eventos de auditoria
func generateEventID() string {
	return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}

// extractUserID extrai o ID do usuário do formato OpenFGA
func extractUserID(user string) string {
	userID, err := ParseUser(user)
	if err != nil {
		return user // Retornar como está se não conseguir extrair
	}
	return userID
}

// === AUDIT METRICS ===

// AuditMetrics coleta métricas de auditoria
type AuditMetrics struct {
	mu                sync.RWMutex
	permissionChecks  int64
	permissionGrants  int64
	permissionRevokes int64
	suspiciousEvents  int64
}

// NewAuditMetrics cria um novo coletor de métricas de auditoria
func NewAuditMetrics() *AuditMetrics {
	return &AuditMetrics{}
}

// IncrementPermissionChecks incrementa o contador de verificações de permissão
func (am *AuditMetrics) IncrementPermissionChecks() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.permissionChecks++
}

// IncrementPermissionGrants incrementa o contador de concessões de permissão
func (am *AuditMetrics) IncrementPermissionGrants() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.permissionGrants++
}

// IncrementPermissionRevokes incrementa o contador de revogações de permissão
func (am *AuditMetrics) IncrementPermissionRevokes() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.permissionRevokes++
}

// IncrementSuspiciousEvents incrementa o contador de eventos suspeitos
func (am *AuditMetrics) IncrementSuspiciousEvents() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.suspiciousEvents++
}

// GetMetrics retorna as métricas atuais
func (am *AuditMetrics) GetMetrics() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	return map[string]interface{}{
		"permission_checks":  am.permissionChecks,
		"permission_grants":  am.permissionGrants,
		"permission_revokes": am.permissionRevokes,
		"suspicious_events":  am.suspiciousEvents,
	}
}
