package model

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Link represents a Link in the system
type Link struct {
	User         string `json:"user"`
	Link         string `json:"link"`
	Title        string `json:"title"`
	DateModified int64  `json:"date_modified"`
	DateCreated  int64  `json:"date_created"`
	Scope        string `json:"scope"`
	Visits       int64  `json:"visits"`
}

// key gets the hash key for the given Link
func (l *Link) key() string {
	return fmt.Sprintf("link:%s:%s:%s", l.Scope, l.User, l.Title)
}

// Commit saves this Link to a Redis connection
func (l *Link) Commit(s ServerContext, log *zap.SugaredLogger) error {
	key := l.key()
	l.DateModified = time.Now().Unix()
	return s.linker.CommitLink(key, l)
}

// Visit increments the counter for the number of visits
func (l *Link) Visit(s ServerContext, log *zap.SugaredLogger) error {
	l.Visits++
	key := l.key()
	return s.linker.CommitLink(key, l)
}

// GetLink obtains a Link hash from Redis
func GetLink(s ServerContext, scope string, user string, title string) (*Link, error) {
	return s.linker.GetLink(scope, user, title)
}

// SearchForLink looks in a user's personal scope and then the global scope for link
func SearchForLink(s ServerContext, user string, title string) (*Link, error) {

	var (
		pLink *Link
		pErr  error
		gLink *Link
		gErr  error
	)

	// TODO: this should be optimized to buffer the calls
	pLink, pErr = GetLink(s, "personal", user, title)
	gLink, gErr = GetLink(s, "global", user, title)

	if pErr == nil {
		return pLink, nil
	}

	if _, ok := pErr.(NotFoundError); !ok {
		return nil, pErr
	}

	if gErr == nil {
		return gLink, nil
	}

	return nil, gErr

}
