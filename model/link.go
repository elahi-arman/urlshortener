package model

import (
	"fmt"
	"time"
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
func (l *Link) Commit(s ServerContext) error {
	key := l.key()
	l.DateModified = time.Now().Unix()
	s.Log.Debugf("%s Committing link with key %s", s.ID, key)
	return s.Linker.CommitLink(key, l)
}

// Visit increments the counter for the number of visits
func (l *Link) Visit(s ServerContext) error {
	l.Visits++
	key := l.key()
	s.Log.Debugf("%s Visited link with key %s", s.ID, key)
	return s.Linker.CommitLink(key, l)
}

// GetLink obtains a Link hash from Redis
func GetLink(s ServerContext, scope string, user string, title string) (*Link, error) {
	s.Log.Debugf("%s Retrieving link with matching %s %s %s", s.ID, scope, user, title)
	return s.Linker.GetLink(scope, user, title)
}

// SearchForLink looks in a user's personal scope and then the global scope for link
func SearchForLink(s ServerContext, scope string, title string) (*Link, error) {

	var (
		sLink *Link
		sErr  error
		gLink *Link
		gErr  error
	)

	sLink, sErr = GetLink(s, scope, "*", title)
	gLink, gErr = GetLink(s, "global", "*", title)

	if sErr == nil {
		s.Log.Debugf("%s Found link %s at %s scope, returning %#v", s.ID, title, scope, sLink)
		return sLink, nil
	}

	if _, ok := sErr.(NotFoundError); !ok {
		s.Log.Errorw(s.ID, "error", sErr)
		return nil, sErr
	}

	if gErr == nil {
		s.Log.Debugf("%s Found link %s at global scope, returning %#v", s.ID, title, sLink)
		return gLink, nil
	}

	s.Log.Errorw(s.ID, "error", gErr)
	return nil, gErr

}
