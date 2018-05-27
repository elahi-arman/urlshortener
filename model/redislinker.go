package model

import (
	"fmt"

	"github.com/elahi-arman/urlshortener/config"
	"github.com/gomodule/redigo/redis"
)

//RedisLinker defines a Linker for interfacting with Redis
type RedisLinker struct {
	conn redis.Conn
}

//NewRedisLinker creates a Linker backed by Redis
func NewRedisLinker(cfg config.RedisConfig) (Linker, error) {
	c, err := redis.Dial("tcp", cfg.Address)
	if err != nil {
		return nil, err
	}

	return &RedisLinker{
		conn: c,
	}, nil
}

//CommitLink updates the given link (or creates a new one) at the given key
func (r *RedisLinker) CommitLink(key string, l *Link) error {
	_, err := r.conn.Do("HMSET", key,
		"link", l.Link,
		"date_modified", l.DateModified,
		"date_created", l.DateCreated,
		"visits", l.Visits)

	if err != nil {
		return err
	}

	return nil
}

//GetLink retrieves a link from the data store
func (r *RedisLinker) GetLink(scope string, user string, title string) (*Link, error) {

	hash, err := redis.Values(r.conn.Do("HGETALL", fmt.Sprintf("link:%s:%s:%s", scope, user, title)))

	if err != nil {
		return nil, err
	}

	link := bulkHashToLink(hash)
	link.User = user
	link.Scope = scope
	link.Title = title

	return link, nil

}

//GetLinksInScope retrieves all links in the given scope
func (r *RedisLinker) GetLinksInScope(scope string) ([]Link, error) {
	var l []Link
	cursor := 0
	for {
		values, err := redis.Values(r.conn.Do(fmt.Sprintf("HSCAN link:%s:*:* %d", scope, cursor)))

		cursor, err := redis.Int64(values[0], err)
		if err != nil {
			return l, err
		}

		hash, err := redis.Values(values[1], err)
		if err != nil {
			return l, err
		}

		link := bulkHashToLink(hash)
		link.Scope = scope

		l = append(l, *link)
		if cursor == 0 {
			break
		}
	}

	return l, nil
}

//GetLinksByUser retrieves all links (from all scopes) for the given user
func (r *RedisLinker) GetLinksByUser(user string) ([]Link, error) {
	return []Link{}, nil
}

//GetLinksInScopeByUser retrieves all links inthe given scope from the given user
func (r *RedisLinker) GetLinksInScopeByUser(scope string, user string) ([]Link, error) {
	return []Link{}, nil
}

//bulkHashToLink deserializes a redis hash into a Link
// because HSCAN returns array of elements and we don't want to hedge bets
// on the order of array elements, this function does the heavy lifting
func bulkHashToLink(v []interface{}) *Link {
	var link string
	var dateModified int64
	var dateCreated int64
	var visits int64

	for i := 0; i < len(v)-1; i = i + 2 {
		key, err := redis.String(v[i], nil)
		value := v[i+1]
		if err != nil {
			continue
		}

		switch key {

		case "link":
			link, err = redis.String(value, nil)
			if err != nil {
				link = ""
			}

		case "date_modified":
			dateModified, err = redis.Int64(value, nil)
			if err != nil {
				dateModified = -1
			}

		case "date_created":
			dateCreated, err = redis.Int64(value, nil)
			if err != nil {
				dateCreated = -1
			}

		case "visits":
			visits, err = redis.Int64(value, nil)
			if err != nil {
				visits = -1
			}

		default:
		}
	}

	return &Link{
		User:         "",
		Link:         link,
		Title:        "",
		DateModified: dateModified,
		DateCreated:  dateCreated,
		Scope:        "",
		Visits:       visits,
	}
}
