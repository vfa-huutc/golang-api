package utils

// CalculatePagination calculates the offset and limit values for pagination
// Parameters:
//   - page: The page number (1-based indexing). If less than 1, defaults to 1
//   - limit: Number of items per page. If less than 1, defaults to 10
//
// Returns:
//   - offset: The starting position for the current page ((page-1) * limit)
//   - limit: The number of items per page
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
