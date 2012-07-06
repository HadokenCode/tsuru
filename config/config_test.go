package config

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

var expected = map[interface{}]interface{}{
	"database": map[interface{}]interface{}{
		"host": "127.0.0.1",
		"port": 8080,
	},
	"auth": map[interface{}]interface{}{
		"salt": "xpto",
		"key":  "sometoken1234",
	},
	"xpto": "ble",
}

func (s *S) TearDownTest(c *C) {
	Configs = nil
}

func (s *S) TestConfig(c *C) {
	conf := `
database:
  host: 127.0.0.1
  port: 8080
auth:
  salt: xpto
  key: sometoken1234
xpto: ble
`
	err := ReadConfigBytes([]byte(conf))
	c.Assert(err, IsNil)
	c.Assert(Configs, DeepEquals, expected)
}

func (s *S) TestConfigFile(c *C) {
	configFile := "testdata/config.yml"
	err := ReadConfigFile(configFile)
	c.Assert(err, IsNil)
	c.Assert(Configs, DeepEquals, expected)
}

func (s *S) TestGetConfig(c *C) {
	configFile := "testdata/config.yml"
	err := ReadConfigFile(configFile)
	c.Assert(err, IsNil)
	value, err := Get("xpto")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "ble")
	value, err = Get("database:host")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "127.0.0.1")
}

func (s *S) TestGetConfigReturnErrorsIfTheKeyIsNotFound(c *C) {
	configFile := "testdata/config.yml"
	err := ReadConfigFile(configFile)
	c.Assert(err, IsNil)
	value, err := Get("xpta")
	c.Assert(value, IsNil)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "^key xpta not found$")
	value, err = Get("database:hhh")
	c.Assert(value, IsNil)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "^key database:hhh not found$")
}

func (s *S) TestGetString(c *C) {
	configFile := "testdata/config.yml"
	err := ReadConfigFile(configFile)
	c.Assert(err, IsNil)
	value, err := GetString("xpto")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "ble")
	value, err = GetString("database:host")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "127.0.0.1")
}

func (s *S) TestGetStringShouldReturnErrorIfTheKeyDoesNotRepresentAString(c *C) {
	configFile := "testdata/config.yml"
	err := ReadConfigFile(configFile)
	c.Assert(err, IsNil)
	value, err := GetString("database:port")
	c.Assert(value, Equals, "")
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "^key database:port has non-string value$")
}

func (s *S) TestGetStringShouldReturnErrorIfTheKeyDoesNotExist(c *C) {
	configFile := "testdata/config.yml"
	err := ReadConfigFile(configFile)
	c.Assert(err, IsNil)
	value, err := GetString("xpta")
	c.Assert(value, Equals, "")
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "^key xpta not found$")
}

func (s *S) TestSet(c *C) {
	err := ReadConfigFile("testdata/config.yml")
	c.Assert(err, IsNil)
	Set("xpto", "bla")
	value, err := GetString("xpto")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "bla")
}

func (s *S) TestSetChildren(c *C) {
	err := ReadConfigFile("testdata/config.yml")
	c.Assert(err, IsNil)
	Set("database:host", "database.com")
	value, err := GetString("database:host")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "database.com")
}

func (s *S) TestUnset(c *C) {
	err := ReadConfigFile("testdata/config.yml")
	c.Assert(err, IsNil)
	err = Unset("xpto")
	c.Assert(err, IsNil)
	_, err = Get("xpto")
	c.Assert(err, NotNil)
}

func (s *S) TestUnsetChildren(c *C) {
	err := ReadConfigFile("testdata/config.yml")
	c.Assert(err, IsNil)
	err = Unset("database:host")
	c.Assert(err, IsNil)
	_, err = Get("database:host")
	c.Assert(err, NotNil)
}

func (s *S) TestUnsetWithUndefinedKey(c *C) {
	err := ReadConfigFile("testdata/config.yml")
	c.Assert(err, IsNil)
	err = Unset("database:hoster")
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "^Key database:hoster not found$")
}
