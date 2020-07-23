package userdb

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// DefaultTimeout is the context timeout for DB requests
var DefaultTimeout = time.Second * 5

// Salt alters the bcrypt password hash (TODO: parameterize)
var Salt = "local_dev_secret"

var log = logrus.WithFields(logrus.Fields{"service": "web"})

type userDB struct {
	rdb *redis.Client
}

// NewConnection returns a UserDB interface
func NewConnection(addr, password string) UserDB {
	log.Debugf("Creating new redis connection to '%s'", addr)
	return &userDB{rdb: redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})}
}

func (u *userDB) Verify(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	if _, err := u.rdb.Ping(ctx).Result(); err != nil {
		log.Error("Failed ping to redis: ", err)
		return false
	}
	return true
}

func (u *userDB) Get(ctx context.Context, username string) (user User, err error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	raw, err := u.rdb.Get(ctx, "user:"+username).Bytes()
	if err == redis.Nil {
		log.Debugf("Redis returned no data for key '%s'", username)
		// TODO better error
		return
	} else if err != nil {
		log.Error("Error from redis: ", err)
		return
	}
	dec := gob.NewDecoder(bytes.NewReader(raw))
	err = dec.Decode(&user)
	if err != nil {
		log.Error("Error decoding response: ", err)
	}
	return
}

func (u *userDB) Save(ctx context.Context, user *User, ttl time.Duration) (err error) {
	log.Debugf("Updating '%s' (ttl: %s)", user.username, ttl)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(*user)
	if err != nil {
		log.Error("Error encoding user: ", err)
		return
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	err = u.rdb.Set(ctx, user.username, buf, ttl).Err()
	if err != nil {
		log.Error("Error from redis: ", err)
	}
	return
}

func (u *userDB) Login(ctx context.Context, username, password string) (b bool) {
	log.Debug("Logging in user ", username)
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	hash, err := u.rdb.Get(ctx, "login:"+username).Bytes()
	if err == redis.Nil {
		log.Warnf("User '%s' has no login", username)
		return
	} else if err != nil {
		log.Error("Error from redis: ", err)
		return
	}
	return bcrypt.CompareHashAndPassword(hash, []byte(password+Salt)) == nil
}

// TODO: Allow setting expiration?
func (u *userDB) Register(ctx context.Context, username, password string) (err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+Salt), bcrypt.DefaultCost)
	if err != nil {
		log.Print("failed to hash password: ", err)
		return
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	return u.rdb.Set(ctx, "login:"+username, hash, 0).Err()
}
