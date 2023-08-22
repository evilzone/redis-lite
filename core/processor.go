package core

import (
	"bytes"
	"errors"
	"fmt"
	"redis-lite/storage"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RequestProcessor interface {
	Process(request Request) (Response, error)
}

type CommandProcessor struct {
	Cache *storage.InMemoryStorage
}

func (cp *CommandProcessor) Process(request Request) (Response, error) {
	switch request.Command {
	case CMDGet:
		return cp.ProcessGet(request)
	case CMDSet:
		return cp.ProcessSet(request)
	case CMDDel:
		return cp.ProcessDel(request)
	case CMDPing:
		return cp.ProcessPing()
	case CMDExpire:
		return cp.ProcessExpiry(request)
	case CMDTtl:
		return cp.ProcessTTL(request)
	case CMDKeys:
		return cp.ProcessKeys(request)
	default:
		return Response{}, ErrorInvalidCommand
	}
}

func (cp *CommandProcessor) ProcessKeys(request Request) (Response, error) {

	pattern := request.Params[0]
	keys := cp.Cache.Keys()
	fmt.Println(keys)

	if len(keys) == 0 {
		return Response{Success: false, Value: []byte("empty list or set")}, nil
	}

	matchedKeys := make([]string, 0)

	pattern = strings.ReplaceAll(pattern, "*", ".*")
	pattern = strings.ReplaceAll(pattern, "?", ".?")

	for _, key := range keys {
		matched, err := regexp.MatchString(pattern, key)
		fmt.Println("pattern %s ", pattern, " key %s ", key,
			" matched %v ", matched, " err %v ", err)

		if matched {
			matchedKeys = append(matchedKeys, key)
		}
	}

	var buffer bytes.Buffer

	for index := 0; index < len(matchedKeys); index++ {
		str := fmt.Sprintf("%d) %s\n", index+1, matchedKeys[index])
		buffer.WriteString(str)
	}

	return Response{Success: true, Value: buffer.Bytes()}, nil
}

func (cp *CommandProcessor) ProcessGet(request Request) (Response, error) {

	item, err := cp.Cache.Get(request.Params[0])

	if err != nil {
		return Response{Value: []byte("(nil)")}, nil
	}
	return Response{Success: true, Value: item.Value}, nil
}

func (cp *CommandProcessor) ProcessSet(request Request) (Response, error) {

	key := request.Params[0]
	value := request.Params[1]

	if len(request.Params) > 2 {
		ttl, err := strconv.Atoi(request.Params[2])
		if err != nil {
			return Response{}, err
		}
		cp.Cache.Set(key, []byte(value), time.Duration(ttl)*time.Second)

		go cp.ExpireKey(key, time.Duration(ttl)*time.Second)
	} else {
		cp.Cache.Set(key, []byte(value), 0)
	}
	return Response{Success: true}, nil
}

func (cp *CommandProcessor) ProcessDel(request Request) (Response, error) {
	cp.Cache.Delete(request.Params)
	return Response{Success: true,
		Value: []byte(strconv.FormatInt(int64(len(request.Params)), 10))}, nil
}

func (cp *CommandProcessor) ProcessPing() (Response, error) {
	return Response{Success: true, Value: []byte("PONG")}, nil
}

// ProcessExpiry this is the basic implementation without the options {'NX','XX', 'GT','LT'}
func (cp *CommandProcessor) ProcessExpiry(request Request) (Response, error) {
	key := request.Params[0]
	ttl, err1 := strconv.Atoi(request.Params[1])

	// err conditions: if ttl not a valid integer or key not present in the cache
	if err1 != nil {
		return Response{Success: true, Value: []byte("0")}, nil
	}

	item, err2 := cp.Cache.Get(key)

	if err2 != nil {
		return Response{Success: false, Value: []byte("0")}, nil
	}

	cp.Cache.Set(key, item.Value, time.Duration(ttl)*time.Second)
	go cp.ExpireKey(key, time.Duration(ttl)*time.Second)

	return Response{Success: false, Value: []byte("1")}, nil
}

func (cp *CommandProcessor) ProcessTTL(request Request) (Response, error) {
	key := request.Params[0]

	item, err := cp.Cache.Get(key)

	// if key doesn't exist return -2
	if errors.Is(err, storage.ErrKeyNotFound) {
		return Response{Success: false, Value: []byte("-2")}, nil
	}

	// if key exists with no expiry or if the key is expired return -1
	if item.ExpireAt == -1 || errors.Is(err, storage.ErrKeyExpired) {
		return Response{Success: false, Value: []byte("-1")}, nil
	}

	return Response{Success: true,
		Value: []byte(strconv.FormatInt((item.ExpireAt-time.Now().UnixNano())/1e9, 10))}, nil
}

func (cp *CommandProcessor) ExpireKey(key string, ttl time.Duration) {
	<-time.After(ttl)
	cp.Cache.Delete([]string{key})
}
