package util

func TransformTransactionResult(transactionType string, data interface{}) map[string]interface{} {

	result := make(map[string]interface{})
	result[transactionType] = data

	return result

}
