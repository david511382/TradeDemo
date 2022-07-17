package bootstrap

import (
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Var         Var         `yaml:"var"`
	Server      Server      `yaml:"server"`
	RedisConfig RedisConfig `yaml:"redis"`
}

type Var struct {
	UseDebug bool   `yaml:"use_debug"`
	LogDir   string `yaml:"log_dir"`
	TimeZone string `yaml:"time_zone"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (c *Server) Addr() string {
	if c.Port > 0 {
		return c.Host + ":" + strconv.Itoa(c.Port)
	}
	return c.Host
}

type Db struct {
	Server   `yaml:"server"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Param    string `yaml:"param"`
	// add s after protocol to enable https
	Protocol string `yaml:"protocol"`
}

// protocol://user:password@server:port/database?key=value
func (db *Db) ScanUrl(url string) error {
	if protocolStrs := strings.Split(url, "://"); len(protocolStrs) == 2 {
		db.Protocol = protocolStrs[0]
		url = protocolStrs[1]
	}
	if userPasswordStrs := strings.Split(url, "@"); len(userPasswordStrs) == 2 {
		if userStrs := strings.Split(userPasswordStrs[0], ":"); len(userStrs) == 2 {
			db.User = userStrs[0]
			db.Password = userStrs[1]
		}
		url = userPasswordStrs[1]
	}
	if paramStrs := strings.Split(url, "?"); len(paramStrs) == 2 {
		db.Param = paramStrs[1]
		url = paramStrs[0]
	}
	if dbStrs := strings.Split(url, "/"); len(dbStrs) == 2 {
		db.Database = dbStrs[1]
		url = dbStrs[0]
	}
	if portStrs := strings.Split(url, ":"); len(portStrs) == 2 {
		if port, err := strconv.Atoi(portStrs[1]); err != nil {
			return err
		} else {
			db.Port = port
		}
		url = portStrs[0]
	}
	db.Server.Host = url
	return nil
}

// protocol://user:password@server:port/database?key=value
func (db *Db) ParseToUrl() (url string) {
	if db.Protocol != "" {
		url = fmt.Sprintf(
			"%s://",
			db.Protocol,
		)
	}
	if db.User != "" || db.Password != "" {
		url = fmt.Sprintf(
			"%s%s:%s@",
			url,
			db.User,
			db.Password,
		)
	}
	url = fmt.Sprintf(
		"%s%s",
		url,
		db.Server.Addr(),
	)
	if db.Database != "" {
		url = fmt.Sprintf(
			"%s/%s",
			url,
			db.Database,
		)
	}
	if db.Param != "" {
		url = fmt.Sprintf(
			"%s?%s",
			url,
			db.Param,
		)
	}
	return
}

type RedisConfig struct {
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxOpenConns int `yaml:"max_open_conns"`
	MaxLifeHour  int `yaml:"max_lifehour"`

	RedisTradeKeyRoot string `yaml:"redis_trade_key_root"`

	Trade Db `yaml:"trade"`
}
