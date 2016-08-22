// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

// --------------------------------------------------

type Session struct {
	sessionName     string
	sessionID       string
	sessionFileName string
}

type SessionManager struct {
	sessionName string
	maxLifeTime int
}

// session manager is managed by session manager manager
// session manager manager is passed for context by default key
type sessionManagerManager map[string]*SessionManager

// --------------------------------------------------

// key for context's metadata
const sessionDefaultKey = "assdfqw82987ep98dofi22asdoif9j"

// session file directory. do not change it.
var tempDir = "/tmp"

// var tempDir, _ = ioutil.TempDir("", "mixo-sessions")

func generateSessionFileName(sessionName, sessionID string) string {
	buffer := bytes.NewBufferString(tempDir)
	buffer.WriteRune('/')
	buffer.WriteString(sessionName)
	buffer.WriteRune('_')
	buffer.WriteString(sessionID)
	return buffer.String()
}

// generate random session id
func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// --------------------------------------------------

func SessionStart(sessionName string, maxLifeTime int) Handler {
	sessionManager := &SessionManager{
		sessionName: sessionName,
		maxLifeTime: maxLifeTime,
	}
	// return handler for router
	return func(context *Context) {
		tmp, exists := context.GetMetadata(sessionDefaultKey)
		var smm sessionManagerManager
		if !exists {
			smm = sessionManagerManager{}
		} else {
			smm = tmp.(sessionManagerManager)
		}
		// register session manager
		smm[sessionName] = sessionManager
		context.SetMetadata(sessionDefaultKey, smm)
	}
}

func GetSessionManager(context *Context, sessionName string) *SessionManager {
	tmp, _ := context.GetMetadata(sessionDefaultKey)
	smm := tmp.(sessionManagerManager)
	sessionManager, _ := smm[sessionName]
	return sessionManager
}

// --------------------------------------------------

// generate new session and bake to cookie.
func (sessionManager *SessionManager) NewSession(context *Context) *Session {
	sessionName := sessionManager.sessionName
	sessionID := url.QueryEscape(generateSessionID())

	newSession := &Session{
		sessionName:     sessionName,
		sessionID:       sessionID,
		sessionFileName: generateSessionFileName(sessionName, sessionID),
	}
	// baking to cookie
	cookie := &http.Cookie{
		Name:     sessionManager.sessionName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   sessionManager.maxLifeTime,
	}
	http.SetCookie(context.ResponseWriter, cookie)
	return newSession
}

// load session files and return session object
func (sessionManager *SessionManager) GetSession(context *Context) *Session {
	request := context.Request
	sessionName := sessionManager.sessionName

	cookie, err := request.Cookie(sessionName)
	if err != nil {
		return nil
	}
	sessionID := cookie.Value
	sessionFileName := generateSessionFileName(sessionName, sessionID)

	return &Session{
		sessionName:     sessionName,
		sessionID:       sessionID,
		sessionFileName: sessionFileName,
	}
}

// --------------------------------------------------

// regular expression object to divide key and value of session file
var divideSessionFileLineRegexp = regexp.MustCompile("(.*)=(.*)")

func divideSessionFileLine(lineStr string) (string, string) {
	groups := divideSessionFileLineRegexp.FindSubmatch([]byte(lineStr))
	return string(groups[1]), string(groups[2])
}

// --------------------------------------------------

// set key and value to session file
func (session *Session) Set(key, value string) error {
	buffer := bytes.NewBufferString("")
	file, err := os.Open(session.sessionFileName)
	keyIsFound := false
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineStr := scanner.Text()
			k, _ := divideSessionFileLine(lineStr)

			if k == key {
				buffer.WriteString(key)
				buffer.WriteRune('=')
				buffer.WriteString(value)
				buffer.WriteRune('\n')
				keyIsFound = true
			} else {
				buffer.WriteString(lineStr)
				buffer.WriteRune('\n')
			}
		}
		file.Close()
	}
	// if key isn't found, add new line.
	if !keyIsFound {
		buffer.WriteString(key)
		buffer.WriteRune('=')
		buffer.WriteString(value)
		buffer.WriteRune('\n')
	}

	file, err = os.Create(session.sessionFileName)
	if err != nil {
		return err
	}
	file.Write(buffer.Bytes())
	file.Close()
	return nil
}

// get key and value from session file
func (session *Session) Get(key string) (string, error) {
	file, err := os.Open(session.sessionFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineStr := scanner.Text()
		k, v := divideSessionFileLine(lineStr)
		if k == key {
			return v, nil
		}
	}
	// if cannot find key, return empty string, but return error
	return "", nil
}
