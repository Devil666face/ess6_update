package drw6

import "sync"

type LoadState struct {
	m       sync.RWMutex
	isload  bool
	message string
	err     error
}

func (l *LoadState) Start()              { l.m.Lock(); defer l.m.Unlock(); l.isload, l.err = true, nil }
func (l *LoadState) Stop()               { l.m.Lock(); defer l.m.Unlock(); l.isload = false }
func (l *LoadState) SetMessage(m string) { l.m.Lock(); defer l.m.Unlock(); l.message = m }
func (l *LoadState) SetError(err error)  { l.m.Lock(); defer l.m.Unlock(); l.err = err }
func (l *LoadState) IsLoad() bool        { l.m.RLock(); defer l.m.RUnlock(); return l.isload }
func (l *LoadState) Message() string     { l.m.RLock(); defer l.m.RUnlock(); return l.message }
func (l *LoadState) Error() error        { l.m.RLock(); defer l.m.RUnlock(); return l.err }
