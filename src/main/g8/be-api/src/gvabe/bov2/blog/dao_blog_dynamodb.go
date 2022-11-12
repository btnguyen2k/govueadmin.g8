package blog

import (
	"log"
	"sort"

	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
	"main/src/gvabe/bov2/user"
)

// InitBlogCommentTableDynamodb is helper method to initialize AWS DynamoDB table to store blog comments.
//
// Available since template-v0.4.0
func InitBlogCommentTableDynamodb(adc *prom.AwsDynamodbConnect, tableName string) error {
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	if err := henge.InitDynamodbTables(adc, tableName, spec); err != nil {
		log.Printf("[WARN] error creating table %s (%s): %s\n", tableName, "DynamoDB", err)
		return err
	}

	// var colName, gsiName string

	// colName = CommentFieldOwnerId
	// gsiName = "gsi_" + colName
	// if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
	// 	[]prom.AwsDynamodbNameAndType{{Name: CommentFieldOwnerId, Type: prom.AwsAttrTypeString}},
	// 	[]prom.AwsDynamodbNameAndType{{Name: CommentFieldOwnerId, Type: prom.AwsKeyTypePartition}}); err != nil {
	// 	log.Printf("[WARN] error creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// } else if err := prom.AwsDynamodbWaitForGsiStatus(adc, tableName, gsiName, []string{"ACTIVE"}, 1*time.Second, 10*time.Second); err != nil {
	// 	log.Printf("[WARN] error waiting GSI for to be ACTIVE %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// }
	//
	// colName = CommentFieldPostId
	// gsiName = "gsi_" + colName
	// if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
	// 	[]prom.AwsDynamodbNameAndType{{Name: CommentFieldPostId, Type: prom.AwsAttrTypeString}},
	// 	[]prom.AwsDynamodbNameAndType{{Name: CommentFieldPostId, Type: prom.AwsKeyTypePartition}}); err != nil {
	// 	log.Printf("[WARN] error creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// } else if err := prom.AwsDynamodbWaitForGsiStatus(adc, tableName, gsiName, []string{"ACTIVE"}, 1*time.Second, 10*time.Second); err != nil {
	// 	log.Printf("[WARN] waiting GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// }

	return nil
}

// NewBlogCommentDaoDynamodb is helper method to create AWS DynamoDB-implementation of BlogCommentDao.
//
// Available since template-v0.3.0
func NewBlogCommentDaoDynamodb(adc *prom.AwsDynamodbConnect, tableName string) BlogCommentDao {
	dao := &BaseBlogCommentDaoImpl{}
	spec := &henge.DynamodbDaoSpec{}
	dao.UniversalDao = henge.NewUniversalDaoDynamodb(adc, tableName, spec)
	return dao
}

/*----------------------------------------------------------------------*/

// InitBlogPostTableDynamodb is helper method to initialize AWS DynamoDB table to store blog posts.
//
// Available since template-v0.4.0
func InitBlogPostTableDynamodb(adc *prom.AwsDynamodbConnect, tableName string) error {
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	if err := henge.InitDynamodbTables(adc, tableName, spec); err != nil {
		log.Printf("[WARN] creating table %s (%s): %s\n", tableName, "DynamoDB", err)
		return err
	}

	// var colName, gsiName string

	// colName = henge.FieldTimeCreated
	// gsiName = "gsi_" + colName
	// if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
	// 	[]prom.AwsDynamodbNameAndType{{Name: henge.FieldTimeCreated, Type: prom.AwsAttrTypeString}},
	// 	[]prom.AwsDynamodbNameAndType{{Name: PostFieldOwnerId, Type: prom.AwsKeyTypePartition}}); err != nil {
	// 	log.Printf("[WARN] error creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// } else if err := prom.AwsDynamodbWaitForGsiStatus(adc, tableName, gsiName, []string{"ACTIVE"}, 1*time.Second, 10*time.Second); err != nil {
	// 	log.Printf("[WARN] error waiting GSI for to be ACTIVE %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// }

	// colName = PostFieldOwnerId
	// gsiName = "gsi_" + colName
	// if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
	// 	[]prom.AwsDynamodbNameAndType{{Name: PostFieldOwnerId, Type: prom.AwsAttrTypeString}},
	// 	[]prom.AwsDynamodbNameAndType{{Name: PostFieldOwnerId, Type: prom.AwsKeyTypePartition}}); err != nil {
	// 	log.Printf("[WARN] error creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// } else if err := prom.AwsDynamodbWaitForGsiStatus(adc, tableName, gsiName, []string{"ACTIVE"}, 1*time.Second, 10*time.Second); err != nil {
	// 	log.Printf("[WARN] error waiting GSI for to be ACTIVE %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// }

	// colName = PostFieldIsPublic
	// gsiName = "gsi_" + colName
	// if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
	// 	[]prom.AwsDynamodbNameAndType{{Name: PostFieldIsPublic, Type: prom.AwsAttrTypeNumber}},
	// 	[]prom.AwsDynamodbNameAndType{{Name: PostFieldIsPublic, Type: prom.AwsKeyTypePartition}}); err != nil {
	// 	log.Printf("[WARN] error creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// } else if err := prom.AwsDynamodbWaitForGsiStatus(adc, tableName, gsiName, []string{"ACTIVE"}, 1*time.Second, 10*time.Second); err != nil {
	// 	log.Printf("[WARN] error waiting GSI for to be ACTIVE %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// }

	return nil
}

// NewBlogPostDaoDynamodb is helper method to create AWS DynamoDB-implementation of BlogPostDao.
//
// Available since template-v0.3.0
func NewBlogPostDaoDynamodb(adc *prom.AwsDynamodbConnect, tableName string) BlogPostDao {
	dao := &DynamodbBlogPostDaoImpl{&BaseBlogPostDaoImpl{}}
	spec := &henge.DynamodbDaoSpec{}
	udaoDynamodb := henge.NewUniversalDaoDynamodb(adc, tableName, spec)
	{
		// specific to DynamoDB: map search field(s) to GSI name(s).
		// udaoDynamodb.MapGsi("gsi_"+PostFieldOwnerId, PostFieldOwnerId)
		// udaoDynamodb.MapGsi("gsi_"+PostFieldIsPublic, PostFieldIsPublic)
		// udaoDynamodb.MapGsi("gsi_"+henge.FieldTimeCreated, henge.FieldTimeCreated)
	}
	dao.UniversalDao = udaoDynamodb
	return dao
}

type DynamodbBlogPostDaoImpl struct {
	*BaseBlogPostDaoImpl
}

// GetUserFeedN implements BlogPostDao.GetUserFeedN
func (dao *DynamodbBlogPostDaoImpl) GetUserFeedN(user *user.User, fromOffset, maxNumRows int) ([]*BlogPost, error) {
	filter := (&godal.FilterOptOr{}).
		Add(&godal.FilterOptFieldOpValue{FieldName: PostFieldOwnerId, Operator: godal.FilterOpEqual, Value: user.GetId()}).
		Add(&godal.FilterOptFieldOpValue{FieldName: PostFieldIsPublic, Operator: godal.FilterOpEqual, Value: 1})
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, nil)
	if err != nil {
		return nil, err
	}
	result := make([]*BlogPost, 0)
	for _, ubo := range uboList {
		app := NewBlogPostFromUbo(ubo)
		result = append(result, app)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].GetTimeCreated().After(result[j].GetTimeCreated())
	})
	return result, nil
}

// GetUserFeedAll implements BlogPostDao.GetUserFeedAll
func (dao *DynamodbBlogPostDaoImpl) GetUserFeedAll(user *user.User) ([]*BlogPost, error) {
	return dao.GetUserFeedN(user, 0, 0)
}

// GetUserPostsN implements BlogPostDao.GetUserPostsN
func (dao *DynamodbBlogPostDaoImpl) GetUserPostsN(user *user.User, fromOffset, maxNumRows int) ([]*BlogPost, error) {
	filter := &godal.FilterOptFieldOpValue{FieldName: PostFieldOwnerId, Operator: godal.FilterOpEqual, Value: user.GetId()}
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, nil)
	if err != nil {
		return nil, err
	}
	result := make([]*BlogPost, 0)
	for _, ubo := range uboList {
		app := NewBlogPostFromUbo(ubo)
		result = append(result, app)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].GetTimeCreated().After(result[j].GetTimeCreated())
	})
	return result, nil
}

// GetUserPostsAll implements BlogPostDao.GetUserPostsAll
func (dao *DynamodbBlogPostDaoImpl) GetUserPostsAll(user *user.User) ([]*BlogPost, error) {
	return dao.GetUserPostsN(user, 0, 0)
}

/*----------------------------------------------------------------------*/

// InitBlogVoteTableDynamodb is helper method to initialize AWS DynamoDB table to store blog votes.
//
// Available since template-v0.4.0
func InitBlogVoteTableDynamodb(adc *prom.AwsDynamodbConnect, tableName string) error {
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	if err := henge.InitDynamodbTables(adc, tableName, spec); err != nil {
		log.Printf("[WARN] creating table %s (%s): %s\n", tableName, "DynamoDB", err)
		return err
	}

	// var colName, gsiName string
	//
	// colName = VoteFieldOwnerId + "_" + VoteFieldTargetId
	// gsiName = "gsi_" + colName
	// if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
	// 	[]prom.AwsDynamodbNameAndType{{Name: VoteFieldOwnerId, Type: prom.AwsAttrTypeString}, {Name: VoteFieldTargetId, Type: prom.AwsAttrTypeString}},
	// 	[]prom.AwsDynamodbNameAndType{{Name: VoteFieldOwnerId, Type: prom.AwsKeyTypePartition}, {Name: VoteFieldTargetId, Type: prom.AwsKeyTypeSort}}); err != nil {
	// 	log.Printf("[WARN] error creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// } else if err := prom.AwsDynamodbWaitForGsiStatus(adc, tableName, gsiName, []string{"ACTIVE"}, 1*time.Second, 10*time.Second); err != nil {
	// 	log.Printf("[WARN] error waiting GSI for to be ACTIVE %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// }
	//
	// colName = VoteFieldTargetId + "_" + VoteFieldValue
	// gsiName = "gsi_" + colName
	// if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
	// 	[]prom.AwsDynamodbNameAndType{{Name: VoteFieldTargetId, Type: prom.AwsAttrTypeString}, {Name: VoteFieldValue, Type: prom.AwsAttrTypeNumber}},
	// 	[]prom.AwsDynamodbNameAndType{{Name: VoteFieldTargetId, Type: prom.AwsKeyTypePartition}, {Name: VoteFieldValue, Type: prom.AwsKeyTypeSort}}); err != nil {
	// 	log.Printf("[WARN] error creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// } else if err := prom.AwsDynamodbWaitForGsiStatus(adc, tableName, gsiName, []string{"ACTIVE"}, 1*time.Second, 10*time.Second); err != nil {
	// 	log.Printf("[WARN] error waiting GSI for to be ACTIVE %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	// 	return err
	// }

	return nil
}

// NewBlogVoteDaoDynamodb is helper method to create AWS DynamoDB-implementation of BlogVoteDao.
//
// Available since template-v0.3.0
func NewBlogVoteDaoDynamodb(adc *prom.AwsDynamodbConnect, tableName string) BlogVoteDao {
	dao := &BaseBlogVoteDaoImpl{}
	spec := &henge.DynamodbDaoSpec{UidxAttrs: [][]string{{VoteFieldOwnerId, VoteFieldTargetId}}}
	dao.UniversalDao = henge.NewUniversalDaoDynamodb(adc, tableName, spec)
	return dao
}
