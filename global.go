package gotenv

import (
	"os"
	"strings"

	"github.com/spf13/cast"
)

var (
	defaultEnv = make(map[string]interface{})
)

// SetDefault sets the default value for this key.
// SetDefault is case-insensitive for a key.
// Default only used when no value is provided by the user via ENV.
func SetDefault(key string, value interface{}) {
	// If alias passed in, then set the proper default
	key = strings.ToLower(key)
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, ".")
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(defaultEnv, path[0:len(path)-1])
	// set innermost value
	deepestMap[lastKey] = value
}

// toCaseInsensitiveValue checks if the value is a  map;
// if so, create a copy and lower-case the keys recursively.
func toCaseInsensitiveValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		value = copyAndInsensitiviseMap(cast.ToStringMap(v))
	case map[string]interface{}:
		value = copyAndInsensitiviseMap(v)
	}

	return value
}

// copyAndInsensitiviseMap behaves like insensitiviseMap, but creates a copy of
// any map it makes case insensitive.
func copyAndInsensitiviseMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})

	for key, val := range m {
		lkey := strings.ToLower(key)
		switch v := val.(type) {
		case map[interface{}]interface{}:
			nm[lkey] = copyAndInsensitiviseMap(cast.ToStringMap(v))
		case map[string]interface{}:
			nm[lkey] = copyAndInsensitiviseMap(v)
		default:
			nm[lkey] = v
		}
	}

	return nm
}

// deepSearch scans deep maps, following the key indexes listed in the
// sequence "path".
// The last value is expected to be another map, and is returned.
//
// In case intermediate keys do not exist, or map to a non-map value,
// a new map is created and inserted, and the search continues from there:
// the initial map "m" may be modified!
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}

// Get returns the value associated with the key as an interface.
func Get(key string) interface{} {
	key = strings.ToLower(key)
	val := os.Getenv(key)
	if val == "" {
		if v, ok := defaultEnv[key]; ok {
			return v
		}
	}

	return val
}

// GetString returns the value associated with the key as a string.
func GetString(key string) string {
	return cast.ToString(Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool {
	return cast.ToBool(Get(key))
}

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int {
	return cast.ToInt(Get(key))
}

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 {
	return cast.ToInt32(Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 {
	return cast.ToInt64(Get(key))
}

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint {
	return cast.ToUint(Get(key))
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 {
	return cast.ToUint32(Get(key))
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 {
	return cast.ToUint64(Get(key))
}
