package model

import (
	"fmt"

	"github.com/elahi-arman/urlshortener/config"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

//RedisLinker defines a Linker for interfacting with Redis
type RedisLinker struct {
	conn redis.Conn
}

//NewRedisLinker creates a Linker backed by Redis
func NewRedisLinker(cfg config.RedisConfig) (Linker, error) {
	c, err := redis.Dial("tcp", cfg.Address)
	if err != nil {
		log.Error(fmt.Sprintf("74533::Could not connect to redis %e", err))
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
		log.Error(fmt.Sprintf("95670::REDIS ERROR: %e", err))
		return err
	}

	log.Info(fmt.Sprintf("95671::Committed link %#v", l))
	return nil
}

//GetLink retrieves a link from the data store
func (r *RedisLinker) GetLink(scope string, user string, title string) (*Link, error) {

	hash, err := redis.Values(r.conn.Do("HGETALL", fmt.Sprintf("link:%s:%s:%s", scope, user, title)))

	if err != nil {
		log.Errorf("41231::REDIS ERROR: %e", err)
		return nil, err
	}

	link := bulkHashToLink(hash)
	link.User = user
	link.Scope = scope
	link.Title = title
	log.Debugf("41233::Filled in user, scope, title values")

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
			log.Errorf("60234::REDIS ERROR with cursor: %e", err)
			return l, err
		}

		hash, err := redis.Values(values[1], err)
		if err != nil {
			log.Errorf("60235::REDIS ERROR with latest hash: %#v %e", values[0], err)
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
			log.Errorf("48716::Could not convert key to String %e", err)
			continue
		}

		log.Debugf("48711::Trying to correlate %s %s", key, value)

		switch key {

		case "link":
			link, err = redis.String(value, nil)
			if err != nil {
				log.Errorf("48712::Couldn't convert link to string: %e", err)
				link = ""
			}

		case "date_modified":
			dateModified, err = redis.Int64(value, nil)
			if err != nil {
				log.Errorf("48713::Couldn't convert date_modified to int64: %e", err)
				dateModified = -1
			}

		case "date_created":
			dateCreated, err = redis.Int64(value, nil)
			if err != nil {
				log.Errorf("48714::Couldn't convert date_created to int64: %e", err)
				dateCreated = -1
			}

		case "visits":
			visits, err = redis.Int64(value, nil)
			if err != nil {
				log.Errorf("48713::Couldn't convert date_modified to int64: %e", err)
				visits = -1
			}

		default:
			log.Warn("48715::Could not correlate last key value pair")
		}
	}

	log.Infof("48714::Returning Link {link:%s, dateModified:%d, dateCreated:%d, visits:%d}", link, dateModified, dateCreated, visits)
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
