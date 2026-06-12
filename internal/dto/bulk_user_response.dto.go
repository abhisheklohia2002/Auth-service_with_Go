package dto

type BulkUsersUploadResponse struct {
	TotalRows    int               `json:"totalRows"`
	SuccessCount int               `json:"successCount"`
	FailedCount  int               `json:"failedCount"`
	Errors       []BulkUsersError `json:"errors"`
}

type BulkUsersError struct {
	Row     int    `json:"row"`
	Email   string `json:"email,omitempty"`
	EmpID   string `json:"employeeId,omitempty"`
	Message string `json:"message"`
}
