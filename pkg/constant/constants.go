package constant

import "errors"

const (
	PageNumQueryKey  = "pageNum"
	PageSizeQueryKey = "pageSize"

	BatchOperationUpdate = "update"
	BatchOperationCreate = "create"
	BatchOperationDelete = "delete"
)

var (
	NotSupportedBatchOperation = errors.New("not supported operation")
)
