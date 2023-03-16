package tgmux

import "sync"

type UserStateInterface interface {
	GetCurrentFunction() string
	SetCurrentFunction(function string)
	GetData() map[string]interface{}
	SetData(data map[string]interface{})
	UpdateData(key string, value interface{})
}

type UserStateManagerInterface interface { //TODO Redis
	GetUserState(userID int64) UserStateInterface
	SetUserStage(userID int64, function, stage string)
	ResetUserFunction(userID int64)
}

// UserStateManager is a struct responsible for managing user states for a Telegram bot.
type UserStateManager struct {
	userStates map[int64]*UserState
	mu         sync.RWMutex
}

// UserState is a struct that represents an individual user state.
type UserState struct {
	currentFunction string
	data            map[string]interface{}
	mu              sync.RWMutex
}

func NewUserState() *UserState {
	return &UserState{
		currentFunction: "",
		data:            make(map[string]interface{}),
	}
}

func NewUserStateManager() *UserStateManager {
	return &UserStateManager{
		userStates: make(map[int64]*UserState),
	}
}

// GetUserState retrieves the UserState for a given userID or creates a new one if it doesn't exist.
func (m *UserStateManager) GetUserState(userID int64) UserStateInterface {
	m.mu.RLock()
	state, exists := m.userStates[userID]
	m.mu.RUnlock()

	if !exists {
		m.mu.Lock()
		state = &UserState{
			currentFunction: "",
			data:            make(map[string]interface{}),
		}
		m.userStates[userID] = state
		m.mu.Unlock()
	}

	return state
}

// GetCurrentFunction safely retrieves the state(currentFunction) value.
func (u *UserState) GetCurrentFunction() string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.currentFunction
}

// SetCurrentFunction safely sets the currentFunction value.
func (u *UserState) SetCurrentFunction(function string) {
	u.mu.Lock()
	u.currentFunction = function
	u.mu.Unlock()
}

// GetData safely retrieves the data value.
func (u *UserState) GetData() map[string]interface{} {
	u.mu.RLock()
	defer u.mu.RUnlock()
	// Return a shallow copy of the map to avoid concurrent modification issues
	dataCopy := make(map[string]interface{}, len(u.data))
	for key, value := range u.data {
		dataCopy[key] = value
	}
	return dataCopy
}

// SetData safely sets the data value.
func (u *UserState) SetData(data map[string]interface{}) {
	u.mu.Lock()
	u.data = data
	u.mu.Unlock()
}

// UpdateData safely updates the data map with the provided key and value.
func (u *UserState) UpdateData(key string, value interface{}) {
	u.mu.Lock()
	u.data[key] = value
	u.mu.Unlock()
}

// SetUserStage sets the user's stage by updating their current function and stage.
func (m *UserStateManager) SetUserStage(userID int64, function, stage string) {
	state := m.GetUserState(userID)
	state.SetCurrentFunction(function)
}

// ResetUserFunction resets the user's current function and clears their data.
func (m *UserStateManager) ResetUserFunction(userID int64) {
	state := m.GetUserState(userID)
	state.SetCurrentFunction("")
	state.SetData(make(map[string]interface{}))
}
