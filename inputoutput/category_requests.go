package inputoutput

type CategoryObject struct {
	ID          int64
	Name        string
	Description string
}

type CreateCategoryRequest struct {
	Name        string
	Description string
}

type CreateCategoryResponse struct {
	ID int64
}

type GetCategoryByIdRequest struct {
	ID int64
}

type GetCategoryByIdResponse struct {
	ID          int64
	Name        string
	Description string
}
