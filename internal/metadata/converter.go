package metadata

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

func (m *articlePKMap) ToStringMap() map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue)
	for k, v := range *m {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			result[string(k)] = s
		}
	}
	return result
}
