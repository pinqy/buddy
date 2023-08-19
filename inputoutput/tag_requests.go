package inputoutput

type TagObject struct {
	ID   int64
	Name string
}

type CreateTagRequest struct {
	Name string
}

type CreateTagResponse struct {
	ID int64
}

type GetTagByIdRequest struct {
	ID int64
}

type GetTagByIdResponse struct {
	ID   int64
	Name string
}
