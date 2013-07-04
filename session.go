package mvc

import (
	"crypto/rand"
	"encoding/binary"
	"net/http"
	"sync"
	"time"
	"unicode/utf8"
)

var SESSION_IDENT = "gomvc_sessionid"

// Stores session information for an individual user
type Session struct {
	id          string
	values      map[string]interface{}
	lastVisited time.Time
}

// returns the item if found, nil if not
// Also returns a boolean indicated whether the item
// was found
func (s *Session) Get(key string) (interface{}, bool) {
	item, exists := s.values[key]
	return item, exists
}

// stores a value in the session
func (s *Session) Put(key string, value interface{}) {
	s.values[key] = value
}

type SessionManager struct {
	sessions       map[string]*Session
	sessionMutex   sync.RWMutex
	sessionTTL     time.Duration // how long do sessions remain active
	collectorDelay time.Duration // how often do we purge expired sessions
}

var sessionMgr *SessionManager
var mgrLock sync.Mutex

// Only to be used as a stand-alone session manager
// Stores and returns a singleton instance
func GetSessionManager() *SessionManager {
	if sessionMgr == nil {
		mgrLock.Lock()
		defer mgrLock.Unlock()
		// double checking nil state to avoid unneccessary locks
		if sessionMgr == nil {
			sessionMgr = &SessionManager{
				sessions:       make(map[string]*Session),
				sessionTTL:     1 * time.Minute,
				collectorDelay: 30 * time.Second,
			}
			sessionMgr.collectExpiredSessions()
		}
	}
	return sessionMgr
}

// Returns a new session manager
func NewSessionManager() *SessionManager {
	// TODO: Read from a config?
	sessionMgr = &SessionManager{
		sessions:       make(map[string]*Session),
		sessionTTL:     1 * time.Hour,
		collectorDelay: 30 * time.Second,
	}
	sessionMgr.collectExpiredSessions()
	return sessionMgr
}

// Gets the session associated with a particular request and handles
// maintaining the session in the response
func (sm *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) *Session {
	sessionId := sm.getSessionId(w, r)

	var session *Session
	sm.sessionMutex.RLock()
	session = sm.sessions[sessionId]
	sm.sessionMutex.RUnlock()
	if session == nil {
		session = NewSession(sessionId)
		sm.sessionMutex.Lock()
		sm.sessions[sessionId] = session
		sm.sessionMutex.Unlock()
	}
	session.lastVisited = time.Now()
	return session
}

// Retreives a session id and maintains the id cookie
func (sm *SessionManager) getSessionId(w http.ResponseWriter, r *http.Request) string {
	idCookie, err := r.Cookie(SESSION_IDENT)

	if err == nil {
		return idCookie.Value
	}

	id := StrongRandomString()

	c := &http.Cookie{Name: SESSION_IDENT, Value: id, Path: "/"}
	http.SetCookie(w, c)

	return id
}

// Creates a new session for the given id
func NewSession(sessionId string) *Session {
	session := new(Session)
	session.id = sessionId
	session.values = make(map[string]interface{})

	return session
}

// Starts a go function to garbage collect expired sessions
func (sm *SessionManager) collectExpiredSessions() {
	go func() {
		for {
			now := time.Now()
			sm.sessionMutex.Lock()
			for id, session := range sm.sessions {
				age := now.Sub(session.lastVisited)
				if age > sm.sessionTTL {
					delete(sm.sessions, id)
				}
			}
			sm.sessionMutex.Unlock()

			time.Sleep(sm.collectorDelay)
		}
	}()
}

// Creates a random string
func StrongRandomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return intToString(n)
}

var DIGITS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// translates an uint64 into a string.
func intToString(n uint64) string {
	if n == 0 {
		return "0"
	}

	var base uint64
	base = uint64(len(DIGITS))
	out := make([]byte, 0, 24)
	var remainder int

	for n > 0 {
		remainder = int(n % base)
		n /= base
		char, _ := utf8.DecodeRuneInString(DIGITS[remainder : remainder+1])
		// This will append in LIFO order, but that doesn't matter for our purposes
		out = append(out, byte(char))
	}
	return string(out)
}
