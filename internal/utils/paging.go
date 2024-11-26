package utils

// CalculatePagination takes the page number and limit as input and returns the offset and limit.
// Example: CalculatePagination(2, 10) returns (10, 10) for page 2 with 10 records per page.
func CalculatePagination(page int, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit if not provided or invalid
	}
	offset := (page - 1) * limit
	return offset, limit
}
