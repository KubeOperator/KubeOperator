package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/spf13/viper"
)

type SessionMgr struct {
	mCookieName  string
	mLock        sync.RWMutex
	mMaxLifeTime int64

	mSessions map[string]*Session
}

var GloablSessionMgr *SessionMgr = nil

func Init() {
	exp := viper.GetInt64("jwt.exp")
	GloablSessionMgr = NewSessionMgr(constant.CookieNameForSessionID, exp)
}

func NewSessionMgr(cookieName string, maxLifeTime int64) *SessionMgr {
	mgr := &SessionMgr{mCookieName: cookieName, mMaxLifeTime: maxLifeTime, mSessions: make(map[string]*Session)}

	go mgr.GC()

	return mgr
}

func (mgr *SessionMgr) StartSession(w http.ResponseWriter, r *http.Request) string {
	mgr.mLock.Lock()
	defer mgr.mLock.Unlock()

	newSessionID := url.QueryEscape(mgr.NewSessionID())

	var session *Session = &Session{
		mSessionID:        newSessionID,
		mLastTimeAccessed: time.Now(),
		mValues:           make(map[interface{}]interface{}),
	}
	mgr.mSessions[newSessionID] = session
	cookie := http.Cookie{Name: mgr.mCookieName, Value: newSessionID, Path: "/", HttpOnly: true, MaxAge: 0}
	http.SetCookie(w, &cookie)

	return newSessionID
}

func (mgr *SessionMgr) EndSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(mgr.mCookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		mgr.mLock.Lock()
		defer mgr.mLock.Unlock()

		delete(mgr.mSessions, cookie.Value)

		expiration := time.Now()
		cookie := http.Cookie{Name: mgr.mCookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (mgr *SessionMgr) EndSessionBy(sessionID string) {
	mgr.mLock.Lock()
	defer mgr.mLock.Unlock()

	delete(mgr.mSessions, sessionID)
}

func (mgr *SessionMgr) SetSessionVal(sessionID string, key interface{}, value interface{}) {
	mgr.mLock.Lock()
	defer mgr.mLock.Unlock()

	if session, ok := mgr.mSessions[sessionID]; ok {
		session.mValues[key] = value
	}
}

func (mgr *SessionMgr) GetSessionVal(sessionID string, key interface{}) (interface{}, bool) {
	mgr.mLock.RLock()
	defer mgr.mLock.RUnlock()

	if session, ok := mgr.mSessions[sessionID]; ok {
		if val, ok := session.mValues[key]; ok {
			return val, ok
		}
	}

	return nil, false
}

func (mgr *SessionMgr) GetSessionIDList() []string {
	mgr.mLock.RLock()
	defer mgr.mLock.RUnlock()

	sessionIDList := make([]string, 0)

	for k := range mgr.mSessions {
		sessionIDList = append(sessionIDList, k)
	}

	return sessionIDList
}

func (mgr *SessionMgr) CheckCookieValid(w http.ResponseWriter, r *http.Request) string {
	var cookie, err = r.Cookie(mgr.mCookieName)

	if cookie == nil ||
		err != nil {
		return ""
	}

	mgr.mLock.Lock()
	defer mgr.mLock.Unlock()

	sessionID := cookie.Value

	if session, ok := mgr.mSessions[sessionID]; ok {
		session.mLastTimeAccessed = time.Now()
		return sessionID
	}

	return ""
}

func (mgr *SessionMgr) GetLastAccessTime(sessionID string) time.Time {
	mgr.mLock.RLock()
	defer mgr.mLock.RUnlock()

	if session, ok := mgr.mSessions[sessionID]; ok {
		return session.mLastTimeAccessed
	}

	return time.Now()
}

func (mgr *SessionMgr) GC() {
	mgr.mLock.Lock()
	defer mgr.mLock.Unlock()

	for sessionID, session := range mgr.mSessions {
		if session.mLastTimeAccessed.Unix()+mgr.mMaxLifeTime < time.Now().Unix() {
			delete(mgr.mSessions, sessionID)
		}
	}

	time.AfterFunc(time.Duration(mgr.mMaxLifeTime)*time.Second, func() { mgr.GC() })
}

func (mgr *SessionMgr) NewSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		nano := time.Now().UnixNano() //微秒
		return strconv.FormatInt(nano, 10)
	}
	return base64.URLEncoding.EncodeToString(b)
}

type Session struct {
	mSessionID        string
	mLastTimeAccessed time.Time
	mValues           map[interface{}]interface{}
}
