package router

import (
	"github.com/fatih/set"
)

type (
	Watcher interface {
		OnBroadCast(BroadCastData)
	}

	Router interface {
		AddWatcher(Watcher)
		DelWatcher(Watcher)
		BroadCast(BroadCastData)
	}

	BroadCastData struct {
		Type  uint64 // 数据类型,含义由用户定义
		Data  []byte // 广播数据,含义由用户定义
		Users int    // 当前在线用户数
	}

	router struct {
		wathers *set.Set
	}
)

func NewRouter() Router {
	return &router{wathers: set.New()}
}

func (r *router) AddWatcher(w Watcher) {
	r.wathers.Add(w)
}

func (r *router) DelWatcher(w Watcher) {
	r.wathers.Remove(w)
}

func (r *router) BroadCast(what BroadCastData) {
	what.Users = r.wathers.Size()
	r.wathers.Each(func(x interface{}) bool {
		x.(Watcher).OnBroadCast(what)
		return true
	})
}
