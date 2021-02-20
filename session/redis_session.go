package session

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"sync"
)

const (
	SeeionFlagNone = iota
	SeeionFlagModify
)

type RedisSession struct {
	sessionId  string
	pool       *redis.Pool
	sessionMap map[string]interface{}
	rwlock     sync.RWMutex
	flag       int
}

//构造函数
func NewRedisSession(id string, pool *redis.Pool) *RedisSession {
	s := &RedisSession{
		sessionId:  id,
		sessionMap: make(map[string]interface{}, 16),
		pool:       pool,
		flag:       SeeionFlagNone,
	}
	return s
}

func (r RedisSession) Set(key string, value interface{}) (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	// 设置值
	r.sessionMap[key] = value
	// 标记记录
	r.flag = SeeionFlagModify
	return
}

func (r RedisSession) Save() (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	// 如果数据没变，不需要存
	if r.flag != SeeionFlagModify {
		return
	}
	data, err := json.Marshal(r.sessionMap)
	if err != nil {
		return
	}

	conn := r.pool.Get()
	_, err = conn.Do("SET", r.sessionId, string(data))
	if err != nil {
		return
	}
	return
}

func (r *RedisSession) Get(key string) (result interface{}, err error) {
	// 加锁
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	result, ok := r.sessionMap
	if !ok {
		err = errors.New("key not exists")

	}
	return
}
