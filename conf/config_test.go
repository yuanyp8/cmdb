package conf_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuanyp8/cmdb/conf"
	"os"
	"testing"
)

var filepath string

func TestLoadConfigFromToml(t *testing.T) {
	should := assert.New(t)
	cfg, err := conf.LoadConfigFromToml(filepath)

	if should.NoError(err) {
		cfg.InitGlobal()
		should.Equal("cmdb", conf.C().App.Name)
	}
}

func TestTestLoadConfigFromTomlAndENV(t *testing.T) {
	should := assert.New(t)
	os.Setenv("MYSQL_DATABASE", "unit_test")

	cfg, err := conf.LoadConfigFromTomlAndENV(filepath)
	if should.NoError(err) {
		cfg.InitGlobal()
		should.Equal("unit_test", conf.C().MySQL.Database)
	}
}

// must to configuration the correct information
func TestGetDB(t *testing.T) {
	should := assert.New(t)
	cfg, err := conf.LoadConfigFromToml(filepath)

	if should.NoError(err) {
		cfg.InitGlobal()
		_, err := conf.C().MySQL.GetDB()
		should.NoError(err)
	}
}

func init() {
	filepath = "testdata/config-good.toml"
}
