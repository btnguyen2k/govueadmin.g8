package blog

import (
	"testing"

	promdynamodb "github.com/btnguyen2k/prom/dynamodb"
	prommongo "github.com/btnguyen2k/prom/mongo"
	promsql "github.com/btnguyen2k/prom/sql"
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
	testAdc        *promdynamodb.AwsDynamodbConnect
	testMc         *prommongo.MongoConnect
	testSqlc       *promsql.SqlConnect
	testDaoComment BlogCommentDao
	testDaoPost    BlogPostDao
	testDaoVote    BlogVoteDao
)
