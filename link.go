package urlshortener

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
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
func (l *Link) Commit(c redis.Conn) bool {
	key := l.key()
	now := time.Now().Unix()
	_, err := c.Do("HMSET", key,
		"link", l.Link,
		"date_modified", now,
		"date_created", l.DateCreated,
		"visits", l.Visits)

	if err != nil {
		fmt.Printf("REDIS ERROR: %#v\n", err)
		return false
	}

	l.DateModified = now
	fmt.Printf("OK")
	return true
}

// Visit increments the counter for the number of visits
func (l *Link) Visit(c redis.Conn) bool {
	key := l.key()

	visits, err := c.Do("HINCRBY", key, "visits", 1)

	if err != nil {
		fmt.Printf("REDIS ERROR: %#v\n", err)
		return false
	}

	fmt.Printf("OK %s", visits)
	return true
}

// GetLink obtains a Link hash from Redis
func GetLink(c redis.Conn, scope string, user string, title string) (Link, error) {
	hash, err := redis.Values(c.Do("HGETALL", fmt.Sprintf("link:%s:%s:%s", scope, user, title)))

	if err != nil {
		fmt.Printf("REDIS ERROR: %#v\n", err)
		return Link{}, err
	}

	link := HashToLink(hash)

	link.User = user
	link.Scope = scope
	link.Title = title

	return link, nil
}

//HashToLink deserializes a redis hash into a Link
func HashToLink(v []interface{}) Link {
	var link string
	var dateModified int64
	var dateCreated int64
	var visits int64
	var err error

	for i := 0; i < len(v)-1; i = i + 2 {
		key := v[i]
		value := v[i+1]
		if key == "link" {
			link, err = redis.String(value, nil)
			if err != nil {
				link = ""
			}
		} else if key == "date_modified" {
			dateModified, err = redis.Int64(value, nil)
			if err != nil {
				dateModified = -1
			}
		} else if key == "date_created" {
			dateCreated, err = redis.Int64(value, nil)
			if err != nil {
				dateCreated = -1
			}
		} else if key == "visits" {
			visits, err = redis.Int64(value, nil)
			if err != nil {
				visits = -1
			}
		}
	}

	return Link{
		User:         "",
		Link:         link,
		Title:        "",
		DateModified: dateModified,
		DateCreated:  dateCreated,
		Scope:        "",
		Visits:       visits,
	}
}
