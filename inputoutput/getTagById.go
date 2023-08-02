package inputoutput

type GetTagByIdRequest struct {
	ID int64
}

type GetTagByIdResponse struct {
	ID   int64
	Name string
}
