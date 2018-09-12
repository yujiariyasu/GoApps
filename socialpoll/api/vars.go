package main

import (
	"net/http"
	"sync"
)

var (
	varsLock sync.RWMutex
	// 値のほうのマップにはリクエストのインスタンスに関連づけたデータが格納される
	vars map[*http.Request]map[string]interface{}
)

// マップvarsを生成
func OpenVars(r *http.Request) {
	varsLock.Lock()
	if vars == nil {
		vars = map[*http.Request]map[string]interface{}{}
	}
	vars[r] = map[string]interface{}{}
	varsLock.Unlock()
}

// マップvarsを削除
func CloseVars(r *http.Request) {
	varsLock.Lock()
	delete(vars, r)
	varsLock.Unlock()
}

func GetVar(r *http.Request, key string) interface{} {
	//RLockを用いると、書き込みが発生していない限り複数の読み出しを同時に行える
	varsLock.RLock()
	value := vars[r][key]
	varsLock.RUnlock()
	return value
}

func SetVar(r *http.Request, key string, value interface{}) {
	varsLock.Lock()
	vars[r][key] = value
	varLock.Unlock()
}
