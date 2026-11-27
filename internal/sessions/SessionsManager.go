package sessions

import (
	"fmt"
	"sync"
	"time"

	"github.com/RandomChugokujin/CruxC2/internal/tasks"
)

type SessionManager interface {
	// bool return values are ok: true for successful operation
    Add(s *Session) bool
    Get(id string) (*Session, bool)
    List() []*Session
    Remove(id string) bool
	QueueTask(id string, task tasks.Task) bool
	UnqueueTask(id string) (*tasks.Task, bool)
	FinishTask(session_id, task_id string) (task tasks.Task, ok bool)
	// SetState(id string, state SessionState)
}

func NewManager() SessionManager {
    return &Manager{sessions: make(map[string]*Session)}
}

type Manager struct {
    mu       sync.RWMutex
    sessions map[string]*Session
}

func (m *Manager) Add(s *Session) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.sessions[s.AgentID]; found {
		return false
	}
	m.sessions[s.AgentID] = s
	return true
}

func (m *Manager) Get(id string) (s *Session, ok bool){
	m.mu.RLock()
	defer m.mu.RUnlock()
	if session, found := m.sessions[id]; found {
		return session, true
	}
	return nil, false
}

func (m *Manager) List() []*Session{
	sessions_slice := make([]*Session, 0)
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, session := range m.sessions {
		sessions_slice = append(sessions_slice, session)
	}
	return sessions_slice
}

func (m* Manager) Remove(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.sessions[id]; found {
		delete(m.sessions, id)
		return true
	}
	return false
}

func (m* Manager) QueueTask(id string, task tasks.Task) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, found := m.sessions[id]; found {
		session.TasksQueued = append(session.TasksQueued, task)
		return true
	}
	return false
}

func (m* Manager) UnqueueTask(id string) (t *tasks.Task, ok bool){
	m.mu.Lock()
	defer m.mu.Unlock()
	session, found := m.sessions[id];
	if !found {
		return nil, false
	}
	if len(session.TasksQueued) == 0 {
		return nil, true
	}
	// Poping the task from TasksQueued
	task := session.TasksQueued[0]
	session.TasksQueued = session.TasksQueued[1:]
	// Appending to TasksExecuting
	session.TasksExecuting = append(session.TasksExecuting, task)
	// Returning
	return &session.TasksExecuting[len(session.TasksExecuting)-1], true
}

func (m* Manager) FinishTask(session_id, task_id string) (task tasks.Task, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	task = tasks.Task{}
	session, found := m.sessions[session_id]
	if !found {
		return task, false
	}
	for index, cur_task := range session.TasksExecuting {
		if cur_task.TaskID == task_id {
			task = cur_task
			// Delete found task
			session.TasksExecuting = append(
				session.TasksExecuting[:index],
				session.TasksExecuting[index+1:]...)
			return task, true
		}
	}
	return task, false
}

type Session struct {
	Hostname 		string
	AgentID 		string
	IPAddr 			string
	Version 		string
	Platform		string
	Architecture 	string
	RegisterTime	time.Time
	LastPoll		time.Time
	TasksQueued		[]tasks.Task
	TasksExecuting	[]tasks.Task
}

func (s *Session) PrintSession() {
	fmt.Printf("Hostname:\t\t\t%s\n", s.Hostname)
	fmt.Printf("Agent ID:\t\t\t%s\n", s.AgentID)
	fmt.Printf("IP Address:\t\t\t%s\n", s.IPAddr)
	fmt.Printf("Version:\t\t\t%s\n", s.Version)
	fmt.Printf("Platform:\t\t\t%s\n", s.Platform)
	fmt.Printf("Architecture:\t\t\t%s\n", s.Architecture)
	fmt.Printf("Time of Registration:\t\t%s\n", s.RegisterTime.Format(time.ANSIC))
	fmt.Printf("Time of Last Poll:\t\t%s\n", s.RegisterTime.Format(time.ANSIC))
	fmt.Printf("Number of tasks queued:\t\t%d\n", len(s.TasksQueued))
	fmt.Printf("Number of tasks executing:\t%d\n\n", len(s.TasksExecuting))
}
