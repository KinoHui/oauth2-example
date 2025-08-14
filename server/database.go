package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// MemoryDB 内存数据库
type MemoryDB struct {
	clients map[string]*ClientInfo
	records map[int64]*AuthRecord
	mutex   sync.RWMutex
	nextID  int64
}

// NewMemoryDB 创建新的内存数据库
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		clients: make(map[string]*ClientInfo),
		records: make(map[int64]*AuthRecord),
		nextID:  1,
	}
}

// GenerateClientID 生成客户端ID
func (db *MemoryDB) GenerateClientID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateClientSecret 生成客户端密钥
func (db *MemoryDB) GenerateClientSecret() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CreateClient 创建客户端
func (db *MemoryDB) CreateClient(req *ClientRegistrationRequest) (*ClientInfo, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	clientID := db.GenerateClientID()
	clientSecret := db.GenerateClientSecret()

	client := &ClientInfo{
		ID:          clientID,
		Secret:      clientSecret,
		Name:        req.Name,
		RedirectURL: req.RedirectURL,
		GrantType:   req.GrantType,
		Scope:       req.Scope,
		AutoApprove: false, // 默认需要用户授权
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	db.clients[clientID] = client
	return client, nil
}

// GetClient 获取客户端信息
func (db *MemoryDB) GetClient(clientID string) (*ClientInfo, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	client, exists := db.clients[clientID]
	return client, exists
}

// CreateAuthRecord 创建权限申请记录
func (db *MemoryDB) CreateAuthRecord(clientID, userID, scope string, approved bool) (*AuthRecord, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	record := &AuthRecord{
		ID:        db.nextID,
		ClientID:  clientID,
		UserID:    userID,
		Scope:     scope,
		Approved:  approved,
		CreatedAt: time.Now(),
	}

	db.records[record.ID] = record
	db.nextID++

	return record, nil
}

// GetAuthRecord 获取权限申请记录
func (db *MemoryDB) GetAuthRecord(clientID, userID, scope string) (*AuthRecord, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, record := range db.records {
		if record.ClientID == clientID && record.UserID == userID && record.Scope == scope {
			return record, true
		}
	}
	return nil, false
}

// SetAutoApprove 设置客户端自动授权
func (db *MemoryDB) SetAutoApprove(clientID string, autoApprove bool) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	client, exists := db.clients[clientID]
	if !exists {
		return fmt.Errorf("client not found: %s", clientID)
	}

	client.AutoApprove = autoApprove
	client.UpdatedAt = time.Now()
	return nil
}

// GetAllClients 获取所有客户端
func (db *MemoryDB) GetAllClients() []*ClientInfo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	clients := make([]*ClientInfo, 0, len(db.clients))
	for _, client := range db.clients {
		clients = append(clients, client)
	}
	return clients
}
