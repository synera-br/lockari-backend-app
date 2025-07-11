package authorization

import (
	"sync"
	"time"
)

// MemoryCache é uma implementação de cache em memória
type MemoryCache struct {
	mu            sync.RWMutex
	data          map[string]*CacheEntry
	maxSize       int
	stats         *CacheStats
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
}

// CacheOptions define as opções para o cache
type CacheOptions struct {
	MaxSize         int
	CleanupInterval time.Duration
}

// NewMemoryCache cria um novo cache em memória
func NewMemoryCache(opts CacheOptions) *MemoryCache {
	if opts.MaxSize <= 0 {
		opts.MaxSize = 1000
	}

	if opts.CleanupInterval <= 0 {
		opts.CleanupInterval = 5 * time.Minute
	}

	cache := &MemoryCache{
		data:        make(map[string]*CacheEntry),
		maxSize:     opts.MaxSize,
		stats:       &CacheStats{MaxSize: int64(opts.MaxSize)},
		stopCleanup: make(chan struct{}),
	}

	// Iniciar limpeza automática
	cache.cleanupTicker = time.NewTicker(opts.CleanupInterval)
	go cache.cleanup()

	return cache
}

// Get recupera um valor do cache
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		c.stats.Misses++
		return nil, false
	}

	// Verificar se expirou
	if entry.IsExpired() {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.data, key)
		c.stats.Evictions++
		c.mu.Unlock()
		c.mu.RLock()
		c.stats.Misses++
		return nil, false
	}

	c.stats.Hits++
	c.updateHitRate()
	return entry.Value, true
}

// Set armazena um valor no cache
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Verificar se precisa fazer limpeza por tamanho
	if len(c.data) >= c.maxSize {
		c.evictLRU()
	}

	entry := &CacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}

	c.data[key] = entry
	c.stats.Size = int64(len(c.data))
}

// Delete remove um valor do cache
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.data[key]; exists {
		delete(c.data, key)
		c.stats.Size = int64(len(c.data))
	}
}

// Clear limpa todo o cache
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*CacheEntry)
	c.stats.Size = 0
	c.stats.Evictions++
}

// Stats retorna estatísticas do cache
func (c *MemoryCache) Stats() *CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Retornar cópia das estatísticas
	return &CacheStats{
		Hits:        c.stats.Hits,
		Misses:      c.stats.Misses,
		Evictions:   c.stats.Evictions,
		Size:        c.stats.Size,
		MaxSize:     c.stats.MaxSize,
		HitRate:     c.stats.HitRate,
		LastCleanup: c.stats.LastCleanup,
	}
}

// Close fecha o cache e para a limpeza automática
func (c *MemoryCache) Close() error {
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
	}

	close(c.stopCleanup)
	return nil
}

// cleanup remove entradas expiradas
func (c *MemoryCache) cleanup() {
	for {
		select {
		case <-c.cleanupTicker.C:
			c.removeExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// removeExpired remove entradas expiradas
func (c *MemoryCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	expiredKeys := make([]string, 0)

	for key, entry := range c.data {
		if now.After(entry.ExpiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(c.data, key)
		c.stats.Evictions++
	}

	c.stats.Size = int64(len(c.data))
	c.stats.LastCleanup = now
}

// evictLRU remove a entrada menos recentemente usada
func (c *MemoryCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.data {
		if oldestKey == "" || entry.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.CreatedAt
		}
	}

	if oldestKey != "" {
		delete(c.data, oldestKey)
		c.stats.Evictions++
	}
}

// updateHitRate atualiza a taxa de hit
func (c *MemoryCache) updateHitRate() {
	total := c.stats.Hits + c.stats.Misses
	if total > 0 {
		c.stats.HitRate = float64(c.stats.Hits) / float64(total)
	}
}

// === CACHE IMPLEMENTATIONS FOR DIFFERENT TYPES ===

// AuthorizationCache é uma implementação específica para cache de autorização
type AuthorizationCache struct {
	cache *MemoryCache
}

// NewAuthorizationCache cria um novo cache de autorização
func NewAuthorizationCache(opts CacheOptions) *AuthorizationCache {
	return &AuthorizationCache{
		cache: NewMemoryCache(opts),
	}
}

// GetPermission recupera uma permissão do cache
func (ac *AuthorizationCache) GetPermission(user, relation, object string) (bool, bool) {
	key := ac.makeKey(user, relation, object)
	value, found := ac.cache.Get(key)
	if !found {
		return false, false
	}

	if allowed, ok := value.(bool); ok {
		return allowed, true
	}

	return false, false
}

// SetPermission armazena uma permissão no cache
func (ac *AuthorizationCache) SetPermission(user, relation, object string, allowed bool, ttl time.Duration) {
	key := ac.makeKey(user, relation, object)
	ac.cache.Set(key, allowed, ttl)
}

// DeletePermission remove uma permissão do cache
func (ac *AuthorizationCache) DeletePermission(user, relation, object string) {
	key := ac.makeKey(user, relation, object)
	ac.cache.Delete(key)
}

// InvalidateUser invalida todas as permissões de um usuário
func (ac *AuthorizationCache) InvalidateUser(user string) {
	ac.cache.mu.Lock()
	defer ac.cache.mu.Unlock()

	keysToDelete := make([]string, 0)
	userPrefix := user + ":"

	for key := range ac.cache.data {
		if len(key) > len(userPrefix) && key[:len(userPrefix)] == userPrefix {
			keysToDelete = append(keysToDelete, key)
		}
	}

	for _, key := range keysToDelete {
		delete(ac.cache.data, key)
	}

	ac.cache.stats.Size = int64(len(ac.cache.data))
}

// InvalidateObject invalida todas as permissões de um objeto
func (ac *AuthorizationCache) InvalidateObject(object string) {
	ac.cache.mu.Lock()
	defer ac.cache.mu.Unlock()

	keysToDelete := make([]string, 0)
	objectSuffix := ":" + object

	for key := range ac.cache.data {
		if len(key) > len(objectSuffix) && key[len(key)-len(objectSuffix):] == objectSuffix {
			keysToDelete = append(keysToDelete, key)
		}
	}

	for _, key := range keysToDelete {
		delete(ac.cache.data, key)
	}

	ac.cache.stats.Size = int64(len(ac.cache.data))
}

// Clear limpa todo o cache
func (ac *AuthorizationCache) Clear() {
	ac.cache.Clear()
}

// Stats retorna estatísticas do cache
func (ac *AuthorizationCache) Stats() *CacheStats {
	return ac.cache.Stats()
}

// Close fecha o cache
func (ac *AuthorizationCache) Close() error {
	return ac.cache.Close()
}

// makeKey cria uma chave única para o cache
func (ac *AuthorizationCache) makeKey(user, relation, object string) string {
	return user + ":" + relation + ":" + object
}

// === CACHE WARMING ===

// CacheWarmer é responsável por pré-carregar o cache
type CacheWarmer struct {
	cache   *AuthorizationCache
	service *Service
}

// NewCacheWarmer cria um novo cache warmer
func NewCacheWarmer(cache *AuthorizationCache, service *Service) *CacheWarmer {
	return &CacheWarmer{
		cache:   cache,
		service: service,
	}
}

// WarmUserPermissions pré-carrega permissões de um usuário
func (cw *CacheWarmer) WarmUserPermissions(userID string, vaultIDs []string) error {
	// TODO: Implementar cache warming
	// Por enquanto, retornamos nil
	return nil
}

// WarmVaultPermissions pré-carrega permissões de um vault
func (cw *CacheWarmer) WarmVaultPermissions(vaultID string, userIDs []string) error {
	// TODO: Implementar cache warming
	// Por enquanto, retornamos nil
	return nil
}

// === CACHE METRICS ===

// CacheMetrics coleta métricas do cache
type CacheMetrics struct {
	cache *AuthorizationCache
}

// NewCacheMetrics cria um novo coletor de métricas
func NewCacheMetrics(cache *AuthorizationCache) *CacheMetrics {
	return &CacheMetrics{
		cache: cache,
	}
}

// CollectMetrics coleta métricas do cache
func (cm *CacheMetrics) CollectMetrics() map[string]interface{} {
	stats := cm.cache.Stats()

	return map[string]interface{}{
		"cache_hits":         stats.Hits,
		"cache_misses":       stats.Misses,
		"cache_evictions":    stats.Evictions,
		"cache_size":         stats.Size,
		"cache_max_size":     stats.MaxSize,
		"cache_hit_rate":     stats.HitRate,
		"cache_last_cleanup": stats.LastCleanup,
	}
}
