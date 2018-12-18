package session

import (
	"../utils"
)

type sessionData struct {
	Username string
}

type Session struct {
	data map[string]*sessionData
}

func NewSession() *Session {
	s := new(Session)

	s.data = make(map[string]*sessionData)

	return s
}

func (s *Session) Init(username string) string {
	sessionId := utils.GenerateId()
	data := &sessionData{Username: username}
	s.data[sessionId] = data

	return sessionId
}

func (s *Session) Get(sessionId string) string { //for check working and contains coockie
	data := s.data[sessionId]
	if data == nil {
		return ""
	}
	return data.Username
}
