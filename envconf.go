package envconf

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Value define the general value type
type Value string

// Int converts Value to int64
func (v Value) Int() int64 {
	value, _ := strconv.ParseInt(string(v), 10, 64)
	return value
}

// Int converts Value to uint64
func (v Value) Uint() uint64 {
	value, _ := strconv.ParseUint(string(v), 10, 64)
	return value
}

// Int converts Value to bool
func (v Value) Bool() bool {
	value, _ := strconv.ParseBool(string(v))
	return value
}

// Int converts Value to string
func (v Value) String() string {
	return string(v)
}

// Float converts Value to float64
func (v Value) Float() float64 {
	value, _ := strconv.ParseFloat(string(v), 64)
	return value
}


// Section is a based on map, not thread safe
type Section map[string]interface{}

// Config struct has map to contains secions and an attribute to indicate the current section
type Config struct {
	Path     string
	Sections map[string]Section
	sync.Mutex
	Section  string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func _getEnv(key string) string {
	return os.Getenv(key)
}

func _setEnv(key, value string) {
	err := os.Setenv(key, value)
	checkError(err)
}

func _string(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func _get(array []interface{}, index int) interface{} {
	if len(array) > index {
		return array[index]
	} else {
		return ""
	}
}

// NewSection will create new section in config
func NewSection(section string) Section {
	sec := make(Section)
	return sec
}

// GetSection gets config section with name
func (c *Config) GetSection(section string) Section {
	c.Lock()
	sec, ok := c.Sections[section]
	if !ok {
		sec = NewSection(section)
		c.Sections[section] = sec
	}
	c.Unlock()
	return sec
}

// Getenv is a sugar for section object to get env
func (sec Section) Getenv(key interface{}) string {
	return _getEnv(_string(key))
}

// Getenv is a sugar for config object to get env
func (c *Config) Getenv(key interface{}) string {
	return _getEnv(_string(key))
}

// Setenv is a sugar for section object to set env
func (sec Section) Setenv(key, value interface{}) {
	_setEnv(_string(key), _string(value))
}

// Setenv is a sugar for config object to set env
func (c *Config) Setenv(key, value interface{}) {
	_setEnv(_string(key), _string(value))
}

// List will list keys section key without order
func (sec Section) List() []string {
	if sec == nil { return nil }
	keys := make([]string, 0, len(sec))
	for k := range sec {
		keys = append(keys, k)
	}
	return keys
}

// Getkey will get key from config
func (sec Section) Getkey(key interface{}) string {
	var configValue string
	_configValue, ok := sec[_string(key)]
	if ok {
		configValue = _string(_configValue)
	}
	return configValue
}

// List keys from current section
func (c *Config) List() []string {
	sec := c.GetSection(c.Section)
	return sec.List()
}

// Getkey will key from config
func (c *Config) Getkey(key interface{}) string {
	sec := c.GetSection(c.Section)
	return sec.Getkey(key)
}

// Put will add new key with or without value
func (sec Section) Put(args ...interface{}) {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when adding a key!")
	}
	sec[_string(args[0])] = _get(args, 1)
}

// Put is a Wrapper for Put
func (c *Config) Put(args ...interface{}) {
	sec := c.GetSection(c.Section)
	sec.Put(args...)
}

// Get will get config key, args pattern: envKey, configKey, defaultValue or just configKey
func (sec Section) Get(args ...interface{}) Value {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	if len(args) == 1 {
		return sec.GetConf(args[0])
	}

	envValue := _getEnv(_string(args[0]))
	if envValue != "" {
		return Value(envValue)
	}

	if len(args) > 1 {
		configValue, ok := sec[_string(args[1])]
		if ok {
			return Value(_string(configValue))
		}
	}

	if len(args) > 2 {
		return Value(_string(args[2]))
	}
	return ""
}

// Fetch will get config key, args pattern: configKey, envKey, defaultValue
func (sec Section) Fetch(args ...interface{}) Value {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	configValue, ok := sec[_string(args[0])]
	if ok {
		return Value(_string(configValue))
	}

	if len(args) > 1 {
		envValue := _getEnv(_string(args[1]))
		if envValue != "" {
			return Value(envValue)
		}
	}

	if len(args) > 2 {
		return Value(_string(args[2]))
	}
	return ""
}

// GetEnv will get config key, args pattern: envKey, defaultValue
func (sec Section) GetEnv(args ...interface{}) Value {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	envValue := _getEnv(_string(args[0]))
	if envValue != "" {
		return Value(envValue)
	}

	if len(args) > 1 {
		return Value(_string(args[1]))
	}
	return ""
}

// String parse config value
//
// args set: (name)
// args set: (name, defaultValue)
func (sec Section) String(args... interface{}) string {
	return sec.GetConf(args...).String()
}

// Int parse config value as int64
//
// args set: (name)
// args set: (name, defaultValue)
func (sec Section) Int(args... interface{}) int64 {
	return sec.GetConf(args...).Int()
}

// Uint parse config value as uint64
//
// args set: (name)
// args set: (name, defaultValue)
func (sec Section) Uint(args... interface{}) uint64 {
	return sec.GetConf(args...).Uint()
}

// Bool parse config value as bool
//
// args set: (name)
// args set: (name, defaultValue)
func (sec Section) Bool(args... interface{}) bool {
	return sec.GetConf(args...).Bool()
}

// Float parse config value as float64
//
// args set: (name)
// args set: (name, defaultValue)
func (sec Section) Float(args... interface{}) float64 {
	return sec.GetConf(args...).Float()
}

// String parse config value as string
//
// args set: (name)
// args set: (name, defaultValue)
func (c *Config) String(args... interface{}) string {
	return c.GetSection(c.Section).String(args...)}

// Int parse env value as int64
//
// args set: (name)
// args set: (name, defaultValue)
func (c *Config) Int(args... interface{}) int64 {
	return c.GetSection(c.Section).Int(args...)
}

// Uint parse env value as uint64
//
// args set: (name)
// args set: (name, defaultValue)
func (c *Config) Uint(args... interface{}) uint64 {
	return c.GetSection(c.Section).Uint(args...)
}

// Float parse env value as float64
//
// args set: (name)
// args set: (name, defaultValue)
func (c *Config) Float(args... interface{}) float64 {
	return c.GetSection(c.Section).Float(args...)
}

// Bool parse env value as bool
//
// args set: (name)
// args set: (name, defaultValue)
func (c *Config) Bool(args... interface{}) bool {
	return c.GetSection(c.Section).Bool(args...)
}

// GetConf will get config key, args pattern: confKey, defaultValue
func (sec Section) GetConf(args ...interface{}) Value {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	configValue, ok := sec[_string(args[0])]
	if ok {
		return Value(_string(configValue))
	}

	if len(args) > 1 {
		return Value(_string(args[1]))
	}
	return ""
}

// Get will get key values from config
//
// At least the key name should be provided
//
// Parameter sets: conf_key
//
// Parameter sets: env_key, conf_key
//
// Parameter sets: env_key, conf_key, default_value
//
// When env_key is provided it will try to fetch env variable first,
// if it's empty, it will try to get it from config
func (c *Config) Get(args ...interface{}) Value {
	sec := c.GetSection(c.Section)
	return sec.Get(args...)
}

// Feth will get key values from config
//
// At least the key name should be provided
//
// Parameter sets: conf_key
//
// Parameter sets: conf_key, env_key
//
// Parameter sets: conf_key, env_key, default_value
//
// When env_key is provided it will try to fetch env variable only
// if the value of conf_key is empty
func (c *Config) Fetch(args ...interface{}) Value {
	sec := c.GetSection(c.Section)
	return sec.Fetch(args...)
}

// GetConf will get key values from config
//
// At least the key name should be provided
//
// Parameter sets: conf_key
//
// Parameter sets: conf_key, default_value
//
func (c *Config) GetConf(args ...interface{}) Value {
	sec := c.GetSection(c.Section)
	return sec.GetConf(args...)
}

// GetEnv will get key values from env
//
// At least the key name should be provided
//
// Parameter sets: env_key
//
// Parameter sets: env_key, default_value
//
func (c *Config) GetEnv(args ...interface{}) Value {
	sec := c.GetSection(c.Section)
	return sec.GetEnv(args...)
}

// ParseValue will parse config from raw string
func ParseValue(value string) string {
	if strings.HasPrefix(value, "\"") {
		tokens := strings.Split(value, "\"")
		if len(tokens) > 1 {
			return tokens[1]
		}
		return value[1:]
	} else {
		tokens := strings.SplitN(value, "#", 2)
		if len(tokens) > 0 {
			return strings.TrimSpace(tokens[0])
		}
	}
	return ""
}

// Parse config file
func (c *Config) parse() {
	configDir, err := filepath.Abs(c.Path)
	if err != nil {
		log.Printf("Failed to locate config file %s, skipped it\n", c.Path)
		return
	}

	f, err := os.Open(configDir)
	if err != nil {
		log.Printf("Failed to load config file %s, skipped it\n", configDir)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var line, key, value string
	sec := c.GetSection("main")
	re := regexp.MustCompile(`\[(.+)\]`)
	for scanner.Scan() {
		raw := scanner.Text()
		// Split line to allow comments
		tokens := strings.Split(raw, "#")
		line = tokens[0]

		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			sec = c.GetSection(strings.TrimSpace(matches[1]))
			continue
		}
		args := strings.SplitN(line, "=", 2)
		if len(args) < 1 {
			continue
		}
		key = strings.TrimSpace(args[0])
		if key == "" {
			continue
		}
		args = strings.SplitN(raw, "=", 2)
		if len(args) > 1 {
			value = ParseValue(strings.TrimSpace(args[1]))
		} else {
			value = ""
		}
		sec.Put(key, value)
	}
	err = scanner.Err()
	checkError(err)
}

// Save saves config file and default keys locally
func (c *Config) Save() {
	configDir, err := filepath.Abs(c.Path)
	checkError(err)

	// make sure dir exists
	dir := filepath.Dir(configDir)
	if _, err = os.Stat(dir); err != nil {
		os.MkdirAll(dir, os.ModePerm)
	}

	// backup old dir
	if _, err = os.Stat(configDir); err == nil {
		err = os.Rename(configDir, configDir+string(time.Now().Format(time.RFC3339)))
		if err != nil {
			log.Fatalf("Failed to rename old config file %v.\n", configDir)
		}
	}

	f, err := os.OpenFile(configDir, os.O_RDWR|os.O_CREATE, 0664)
	checkError(err)
	defer f.Close()

	firstLine := true
	_secTemplate := "[%v]\n"
	c.Lock()
	defer c.Unlock()
	for section, config := range c.Sections {
		secTemplate := _secTemplate
		if !firstLine {
			secTemplate = "\n" + secTemplate
		} else {
			firstLine = false
		}
		_, err = f.WriteString(fmt.Sprintf(secTemplate, section))
		checkError(err)
		for key, value := range config {
			_, err = f.WriteString(fmt.Sprintf("%v = %v\n", key, value))
			checkError(err)
		}
	}
	log.Printf("Saved config file as %s", c.Path)
}

// NewConfig creates main Config instance with specified config file path
func NewConfig(paths ...string) *Config {
	c := &Config{}
	c.Sections = make(map[string]Section)
	c.Section = "main"
	if len(paths) > 0 {
		path := paths[0]
		c.Path = strings.Replace(path, "~", _getEnv("HOME"), 1)
		c.parse()
	}
	return c
}

// NewEmptyConfig creates an empty Config instance
func NewEmptyConfig() *Config {
	c := &Config{}
	c.Sections = make(map[string]Section)
	c.Section = "main"
	return c
}
