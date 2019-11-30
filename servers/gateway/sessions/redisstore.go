package sessions

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	Client          *redis.Client
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	result := RedisStore{client, sessionDuration}
	return &result
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	//TODO: marshal the `sessionState` to JSON and save it in the redis database,
	//using `sid.getRedisKey()` for the key.
	//return any errors that occur along the way.

	sessionJSON, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}

	key := sid.getRedisKey()
	setErr := (*rs).Client.Set(key, sessionJSON, (*rs).SessionDuration).Err()

	return setErr
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	//TODO: get the previously-saved session state data from redis,
	//unmarshal it back into the `sessionState` parameter
	//and reset the expiry time, so that it doesn't get deleted until
	//the SessionDuration has elapsed.

	//for extra-credit using the Pipeline feature of the redis
	//package to do both the get and the reset of the expiry time
	//in just one network round trip!

	key := sid.getRedisKey()

	pipe := (*rs).Client.Pipeline()
	search := pipe.Get(key)
	pipe.Expire(key, (*rs).SessionDuration)
	_, err := pipe.Exec()

	if err != nil && err != redis.Nil {
		return err
	}

	vals, searchErr := search.Result()

	if searchErr != nil {
		if searchErr == redis.Nil {
			return ErrStateNotFound
		} else {
			return searchErr
		}
	}

	unmarshalErr := json.Unmarshal([]byte(vals), &sessionState)
	if unmarshalErr != nil {
		return errors.New("Something wrong for unmarshaling the data")
	}

	return unmarshalErr
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	//TODO: delete the data stored in redis for the provided SessionID
	key := sid.getRedisKey()

	err := (*rs).Client.Del(key).Err()

	return err
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}

func (rs *RedisStore) IncrementFailCount(email string) error {
	emailID := GetEmailID(email)
	count, getErr := rs.GetEmailFailLogIn(email)
	if getErr != nil {
		(*rs).Client.Set(emailID, 0, 300*time.Second)
	}
	setErr := (*rs).Client.Set(emailID, count+1, 5*time.Minute).Err()
	return setErr
}

func (rs *RedisStore) RemoveFailRecord(email string) error {
	key := GetEmailID(email)
	err := (*rs).Client.Del(key).Err()
	return err
}

func (rs *RedisStore) GetEmailFailLogIn(email string) (int, error) {
	emailID := GetEmailID(email)
	getResult, err := (*rs).Client.Get(emailID).Result()
	if err != nil {
		return 0, err
	}
	count, parseErr := strconv.Atoi(getResult)
	return count, parseErr
}

func GetEmailID(email string) string {
	return "fail:" + email
}
