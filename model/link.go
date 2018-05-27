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
func (l *Link) Commit(linker Linker) error {
	key := l.key()
	l.DateModified = time.Now().Unix()
	return linker.CommitLink(key, l)
}

// Visit increments the counter for the number of visits
func (l *Link) Visit(linker Linker) error {
	l.Visits++
	key := l.key()
	return linker.CommitLink(key, l)
}

// GetLink obtains a Link hash from Redis
func GetLink(linker Linker, scope string, user string, title string) (*Link, error) {
	return linker.GetLink(scope, user, title)
}
