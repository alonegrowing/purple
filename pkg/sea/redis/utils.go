package redis

import (
	"github.com/garyburd/redigo/redis"
)

func Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

func ByteSlices(reply interface{}, err error) ([][]byte, error) {
	return redis.ByteSlices(reply, err)
}

func Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

func Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

func Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

func Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

func Int64Map(result interface{}, err error) (map[string]int64, error) {
	return redis.Int64Map(result, err)
}

/*
* IntMap is a helper that converts an array
* of strings (alternating key, value) into a
* map[string]int. The HGETALL commands return
* replies in this format. Requires an even
* number of values in result.
 */
func IntMap(result interface{}, err error) (map[string]int, error) {
	return redis.IntMap(result, err)
}

func Ints(reply interface{}, err error) ([]int, error) {
	return redis.Ints(reply, err)
}

/*
* Positions is a helper that
* converts an array of positions (lat, long) into a [][2]float64.
* The GEOPOS command returns replies in this format.
 */
func Positions(result interface{}, err error) ([]*[2]float64, error) {
	return redis.Positions(result, err)
}

/*
* String is a helper that converts
* a command reply to a string.
* If err is not equal to nil,
* then String returns "", err.
* Otherwise String converts
* the reply to a string as follows:
 */
func String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

/*
* StringMap is a helper that converts an
* array of strings (alternating key, value)
* into a map[string]string. The HGETALL
* and CONFIG GET commands return replies
* in this format. Requires an even number
* of values in result.
 */
func StringMap(result interface{}, err error) (map[string]string, error) {
	return redis.StringMap(result, err)
}

func Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

func Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

/*
* Values is a helper that converts an
* array command reply to a []interface{}.
* If err is not equal to nil, then
* Values returns nil, err. Otherwise,
* Values converts the reply as follows:
*
* Reply type      Result
* array           reply, nil
* nil             nil, ErrNil
* other           nil, error
 */
func Values(reply interface{}, err error) ([]interface{}, error) {
	return redis.Values(reply, err)
}
