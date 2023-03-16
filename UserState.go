package tgmux

type UserState struct {
	CurrentFunction string
	Data            map[string]interface{}
}

func NewUserState() *UserState {
	return &UserState{Data: make(map[string]interface{})}
}

type UserStateManager struct {
	userStates map[int64]*UserState
}

func NewUserStateManager() *UserStateManager {
	return &UserStateManager{
		userStates: make(map[int64]*UserState),
	}
}

func (m *UserStateManager) GetUserState(userID int64) *UserState {
	if _, exists := m.userStates[userID]; !exists {
		m.userStates[userID] = &UserState{
			CurrentFunction: "",
			Data:            make(map[string]interface{}),
		}
	}
	return m.userStates[userID]
}
func (m *UserStateManager) SetUserStage(userID int64, function, stage string) {
	m.userStates[userID].CurrentFunction = function
}

func (m *UserStateManager) ResetUserFunction(userID int64) {
	m.userStates[userID].CurrentFunction = ""
	m.userStates[userID].Data = make(map[string]interface{})
}
