package blog

import (
	"testing"

	"github.com/btnguyen2k/prom"
)

type TestSetupOrTeardownFunc func(t *testing.T, testName string)

func setupTest(t *testing.T, testName string, extraSetupFunc, extraTeardownFunc TestSetupOrTeardownFunc) func(t *testing.T) {
	if extraSetupFunc != nil {
		extraSetupFunc(t, testName)
	}
	return func(t *testing.T) {
		if extraTeardownFunc != nil {
			extraTeardownFunc(t, testName)
		}
	}
}

var (
	testAdc        *prom.AwsDynamodbConnect
	testMc         *prom.MongoConnect
	testSqlc       *prom.SqlConnect
	testDaoComment BlogCommentDao
	testDaoPost    BlogPostDao
	testDaoVote    BlogVoteDao
)

const (
	testTable = "table_temp"
)
