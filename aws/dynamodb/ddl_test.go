package dynamodb_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bradmccoydev/self-service-sdk/aws/dynamodb"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

// Test CreateTable
func TestCreateTable(t *testing.T) {

	// Setup conf test data
	var emptyConf dynamodb.TableConf
	confNoMode := dynamodb.TableConf{TableName: "test", BillingMode: "", ReadCapacityUnits: 0, WriteCapacityUnits: 0}
	validConf := dynamodb.TableConf{TableName: "test", BillingMode: "PAY_PER_REQUEST", ReadCapacityUnits: 0, WriteCapacityUnits: 0}

	// Setup attribute test data
	var emptyInput []dynamodb.TableAttributes
	attribNoName := []dynamodb.TableAttributes{{Name: "", Type: "", IsKey: false, KeyType: ""}}
	attribNoType := []dynamodb.TableAttributes{{Name: "fred", Type: "", IsKey: false, KeyType: ""}}
	attribNoKey := []dynamodb.TableAttributes{{Name: "fred", Type: "S", IsKey: false, KeyType: ""}}
	attribKeyNoType := []dynamodb.TableAttributes{{Name: "fred", Type: "S", IsKey: true, KeyType: ""}}
	attribWithKey := []dynamodb.TableAttributes{{Name: "fred", Type: "S", IsKey: true, KeyType: "HASH"}}

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		conf      dynamodb.TableConf
		attribs   []dynamodb.TableAttributes
		expectErr bool
	}{
		{"No inputs", false, emptyConf, emptyInput, true},
		{"Just session", true, emptyConf, emptyInput, true},
		{"Session & table name", true, confNoMode, emptyInput, true},
		{"Session, table name & billing mode", true, validConf, emptyInput, true},
		{"Session, table name & attribute without a name", true, validConf, attribNoName, true},
		{"Session, table name & attribute without a type", true, validConf, attribNoType, true},
		{"Session, table name & attribute without a key", true, validConf, attribNoKey, true},
		{"Session, table name & key attribute without a type", true, validConf, attribKeyNoType, true},
		{"Session, table name & full key attribute", true, validConf, attribWithKey, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			var sess *session.Session
			if test.validSess {
				sess = internal.CreateAwsSession(true)
			} else {
				sess = internal.CreateAwsSession(false)
			}
			err := dynamodb.CreateTable(sess, test.conf, test.attribs)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
