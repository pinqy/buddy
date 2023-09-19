package inputoutput

type ExpenseObject struct {
	ID       int64
	Category CategoryObject
	Amount   float32
	Day      int64
	Month    int64
	Year     int64
	Location string
	Notes    string
	Tags     []TagObject
}

type CreateExpenseRequest struct {
	CategoryId int64
	Amount     float32
	Day        int64
	Month      int64
	Year       int64
	Location   string
	Notes      string
	TagIDs     []int64
}

type CreateExpenseResponse struct {
	ID int64
}

type GetExpenseByIdRequest struct {
	ID int64
}

type GetExpenseByIdResponse struct {
	ID           int64
	CategoryName string
	Amount       float32
	Day          int64
	Month        int64
	Year         int64
	Location     string
	Notes        string
	TagNames     []string
}

type GetTagsByExpenseIdRequest struct {
	ExpenseID int64
}

type GetTagsByExpenseIdResponse struct {
	Tags []TagObject
}
