package statemachine

import (
	"sync"

	"github.com/looplab/fsm"
)

type Manager struct {
	data map[int64]*fsm.FSM
	mtx  sync.RWMutex
}

func (m *Manager) Set(key int64, value *fsm.FSM) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.data[key] = value
}

func (m *Manager) Get(key int64) (*fsm.FSM, bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	value, ok := m.data[key]
	return value, ok
}

func (m *Manager) Delete(key int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.data, key)
}

func NewManager() *Manager {
	return &Manager{
		data: make(map[int64]*fsm.FSM),
	}
}
