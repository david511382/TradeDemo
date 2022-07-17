package test

import (
	"strconv"
	"strings"
	"testing"
	"zerologix-homework/bootstrap"
	rdsConn "zerologix-homework/src/repo/redis/conn"

	"github.com/google/uuid"
)

func GetTestSchemaName() string {
	testName := uuid.New().String()[:10]
	testName = strings.ReplaceAll(testName, "/", "a")
	testName = strings.ReplaceAll(testName, "-", "a")
	// for mysql schema name
	testName = strings.ReplaceAll(testName, "e", "a")
	testName = strings.ToLower(testName)
	return testName
}

func SetupTestCfg(
	t *testing.T,
	redisCfgRepos []IRedisCfgRepo,
) *bootstrap.Config {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		t.Fatal(errInfo.Error())
	}

	testName := GetTestSchemaName()

	for i, redisCfgRepo := range redisCfgRepos {
		c, _ := redisCfgRepo.Get(cfg)
		if connection, err := rdsConn.Connect(c); err != nil {
			t.Fatal(err.Error())
		} else {
			redisCfgRepo.Set(cfg, testName+strconv.Itoa(i))
			_, keyRoot := redisCfgRepo.Get(cfg)
			t.Cleanup(func() {
				dp := connection.Keys(keyRoot + "*")
				if err := dp.Err(); err != nil {
					t.Fatal(err.Error())
				}
				keys, err := dp.Result()
				if err != nil {
					t.Fatal(err.Error())
				}

				for _, key := range keys {
					dp := connection.Del(key)
					if err := dp.Err(); err != nil {
						t.Error(err.Error())
					}
				}
				_ = connection.Close()
			})
		}
	}

	return cfg
}
