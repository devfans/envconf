package envconf

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Section is a based on map
type Section map[string]interface{}

// Config struct has map to contains secions and an attribute to indicate the current section
type Config struct {
	Path     string
	Sections map[string]Section
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

// Create new section in config
func NewSection(section string) Section {
	sec := make(Section)
	return sec
}

// Get config section with name
func (c *Config) GetSection(section string) Section {
	sec, ok := c.Sections[section]
	if !ok {
		sec = NewSection(section)
		c.Sections[section] = sec
	}
	return sec
}

// Sugar for section object to get env
func (sec *Section) Getenv(key interface{}) string {
	return _getEnv(_string(key))
}

// Sugar for config object to get env
func (c *Config) Getenv(key interface{}) string {
	return _getEnv(_string(key))
}

// Sugar for section object to set env
func (sec *Section) Setenv(key, value interface{}) {
	_setEnv(_string(key), _string(value))
}

// Sugar for config object to set env
func (c *Config) Setenv(key, value interface{}) {
	_setEnv(_string(key), _string(value))
}

// list keys section key without order
func (sec *Section) List() []string {
	if sec == nil { return nil }
	keys := make([]string, 0, len(*sec))
	for k := range *sec {
		keys = append(keys, k)
	}
	return keys
}

// get key from config
func (sec *Section) Getkey(key interface{}) string {
	var configValue string
	_configValue, ok := (*sec)[_string(key)]
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

// Get key from config
func (c *Config) Getkey(key interface{}) string {
	sec := c.GetSection(c.Section)
	return sec.Getkey(key)
}

// Add new key with or without value
func (sec *Section) Put(args ...interface{}) {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when adding a key!")
	}
	(*sec)[_string(args[0])] = _get(args, 1)
}

// Wrapper for Put
func (c *Config) Put(args ...interface{}) {
	sec := c.GetSection(c.Section)
	sec.Put(args...)
}

// Get config key, args pattern: envKey, configKey, defaultValue or just configKey
func (sec *Section) Get(args ...interface{}) string {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	if len(args) == 1 {
		return sec.GetConf(args[0])
	}

	envValue := _getEnv(_string(args[0]))
	if envValue != "" {
		return envValue
	}

	if len(args) > 1 {
		configValue, ok := (*sec)[_string(args[1])]
		if ok {
			return _string(configValue)
		}
	}

	if len(args) > 2 {
		return _string(args[2])
	}
	return ""
}

// Get config key, args pattern: configKey, envKey, defaultValue
func (sec *Section) Fetch(args ...interface{}) string {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	configValue, ok := (*sec)[_string(args[0])]
	if ok {
		return _string(configValue)
	}

	if len(args) > 1 {
		envValue := _getEnv(_string(args[1]))
		if envValue != "" {
			return envValue
		}
	}

	if len(args) > 2 {
		return _string(args[2])
	}
	return ""
}

// Get config key, args pattern: envKey, defaultValue
func (sec *Section) GetEnv(args ...interface{}) string {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	envValue := _getEnv(_string(args[0]))
	if envValue != "" {
		return envValue
	}

	if len(args) > 1 {
		return _string(args[1])
	}
	return ""
}

// Get config key, args pattern: confKey, defaultValue
func (sec *Section) GetConf(args ...interface{}) string {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	configValue, ok := (*sec)[_string(args[0])]
	if ok {
		return _string(configValue)
	}

	if len(args) > 1 {
		return _string(args[1])
	}
	return ""
}

// Get key values from config
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
func (c *Config) Get(args ...interface{}) string {
	sec := c.GetSection(c.Section)
	return sec.Get(args...)
}

// Get key values from config
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
func (c *Config) Fetch(args ...interface{}) string {
	sec := c.GetSection(c.Section)
	return sec.Fetch(args...)
}

// Get key values from config
//
// At least the key name should be provided
//
// Parameter sets: conf_key
//
// Parameter sets: conf_key, default_value
//
func (c *Config) GetConf(args ...interface{}) string {
	sec := c.GetSection(c.Section)
	return sec.GetConf(args...)
}

// Get key values from env
//
// At least the key name should be provided
//
// Parameter sets: env_key
//
// Parameter sets: env_key, default_value
//
func (c *Config) GetEnv(args ...interface{}) string {
	sec := c.GetSection(c.Section)
	return sec.GetEnv(args...)
}

func ParseValue(value string) string {
	if strings.HasPrefix(value, "\"") {
		tokens := strings.Split(value, "\"")
		if len(tokens) > 1 {
			return tokens[1]
		}
		return value[1:]
	} else {
		tokens := strings.SplitN(value, "//", 2)
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
		tokens := strings.Split(raw, "//")
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

// Save config file and default keys locally
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

// Create main Config instance with specified config file path
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
