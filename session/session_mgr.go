package session

//  SessionMgr 管理者
type SessionMgr interface {
	Init(addr string, optings ...string) (err error)
	CreateSession() (session Session, err error)
	Get(sessionId string) (session Session, err error)
}
