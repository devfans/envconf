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

type Section map[string]interface{}

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

// create new section in config
func NewSection(section string) Section {
	sec := make(Section)
	return sec
}

// get config section
func (c *Config) GetSection(section string) Section {
	sec, ok := c.Sections[section]
	if !ok {
		sec = NewSection(section)
		c.Sections[section] = sec
	}
	return sec
}

// sugar for section object to get env
func (sec *Section) Getenv(key interface{}) string {
	return _getEnv(_string(key))
}

// sugar for config object to get env
func (c *Config) Getenv(key interface{}) string {
	return _getEnv(_string(key))
}

// sugar for section object to set env
func (sec *Section) Setenv(key, value interface{}) {
	_setEnv(_string(key), _string(value))
}

// sugar for config object to set env
func (c *Config) Setenv(key, value interface{}) {
	_setEnv(_string(key), _string(value))
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

// get key from config
func (c *Config) Getkey(key interface{}) string {
	sec := c.GetSection(c.Section)
	return sec.Getkey(key)
}

// add new key with or without value
func (sec *Section) Put(args ...interface{}) {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when adding a key!")
	}
	(*sec)[_string(args[0])] = _get(args, 1)
}

// wrapper for Put
func (c *Config) Put(args ...interface{}) {
	sec := c.GetSection(c.Section)
	sec.Put(args...)
}

// get config key, args pattern: envKey, configKey or configKey
func (sec *Section) Get(args ...interface{}) string {
	if len(args) == 0 {
		log.Fatalln("Please at least specify key name when getting a key value!")
	}

	var configValue string
	key := _string(args[0])
	_configValue, ok := (*sec)[key]
	if ok {
		configValue = _string(_configValue)
	}

	if len(args) > 1 {
		envValue := _getEnv(key)
		if envValue != "" {
			return envValue
		}
	}
	return configValue
}

// wrapper for Get
func (c *Config) Get(args ...interface{}) string {
	sec := c.GetSection(c.Section)
	return sec.Get(args...)
}

// parse config file
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
		line = scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			sec = c.GetSection(matches[1])
			continue
		}
		args := strings.Split(line, "=")
    if len(args) < 2 {
      continue
    }
		key = strings.TrimSpace(args[0])
	  value = strings.TrimSpace(args[1])
		sec.Put(key, value)
	}
	err = scanner.Err()
	checkError(err)
}

// save config file and default keys
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

// create main config
func NewConfig(path string) *Config {
	c := &Config{}
	c.Path = strings.Replace(path, "~", _getEnv("HOME"), 1)
	c.Sections = make(map[string]Section)
	c.Section = "main"
	c.parse()
	return c
}
