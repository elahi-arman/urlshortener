package urlshortener

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
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
		log.Error(fmt.Sprintf("95670::REDIS ERROR: %e", err))
		return false
	}

	l.DateModified = now
	log.Info(fmt.Sprintf("95671::Committed link %#v", l))
	return true
}

// Visit increments the counter for the number of visits
func (l *Link) Visit(c redis.Conn) bool {
	log.Debug(fmt.Sprintf("14230::Visiting link %s", l.Title))
	key := l.key()
	log.Debug(fmt.Sprintf("14231::Generated key %s", key))
	visits, err := c.Do("HINCRBY", key, "visits", 1)

	if err != nil {
		log.Error(fmt.Sprintf("14232::REDIS ERROR %#v", err))
		return false
	}

	l.Visits, err = redis.Int64(visits, nil)
	if err != nil {
		log.Error(fmt.Sprintf("14234::Couldn't serialize visits to an int64"))
		return false
	}
	log.Debug(fmt.Sprintf("14233::Incremented number of visits: %d", l.Visits))

	return true
}

// GetLink obtains a Link hash from Redis
func GetLink(c redis.Conn, scope string, user string, title string) (Link, error) {
	log.Info(fmt.Sprintf("41230::Sending HGETALL to Redis for link:%s:%s:%s", scope, user, title))
	hash, err := redis.Values(c.Do("HGETALL", fmt.Sprintf("link:%s:%s:%s", scope, user, title)))

	if err != nil {
		log.Error(fmt.Sprintf("41231::REDIS ERROR: %e", err))
		return Link{}, err
	}

	link := HashToLink(hash)
	link.User = user
	link.Scope = scope
	link.Title = title
	log.Debug(fmt.Sprintf("41233::Filled in user, scope, title values"))

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
		key, err := redis.String(v[i], err)
		if err != nil {
			log.Error(fmt.Sprintf("48716::Could not convert key to String %e", err))
			continue
		}
		value := v[i+1]
		log.Debug(fmt.Sprintf("48711::Trying to correlate %s %s", key, value))
		if key == "link" {
			link, err = redis.String(value, err)
			if err != nil {
				log.Error(fmt.Sprintf("48712::Couldn't convert link to string: %e", err))
				link = ""
			}
		} else if key == "date_modified" {
			dateModified, err = redis.Int64(value, err)
			if err != nil {
				log.Error(fmt.Sprintf("48713::Couldn't convert date_modified to int64: %e", err))
				dateModified = -1
			}
		} else if key == "date_created" {
			dateCreated, err = redis.Int64(value, err)
			if err != nil {
				log.Error(fmt.Sprintf("48714::Couldn't convert date_created to int64: %e", err))
				dateCreated = -1
			}
		} else if key == "visits" {
			visits, err = redis.Int64(value, err)
			if err != nil {
				log.Error(fmt.Sprintf("48713::Couldn't convert date_modified to int64: %e", err))
				visits = -1
			}
		} else {
			log.Warn("48715::Could not correlate last key value pair")
		}
	}

	log.Info(fmt.Sprintf("48714::Returning Link {link:%s, dateModified:%d, dateCreated:%d, visits:%d}", link, dateModified, dateCreated, visits))
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
