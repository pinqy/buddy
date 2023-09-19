package dbmanager

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	io "github.com/pinqy/buddy/inputoutput"
)

func (dbc *DBClient) CreateExpense(req io.CreateExpenseRequest) (io.CreateExpenseResponse, error) {
	var resp io.CreateExpenseResponse

	// Validate required fields
	if req.Amount <= 0.00 {
		return resp, fmt.Errorf("CreateExpense: cannot create expense for negative or zero amount")
	}

	date, err := time.Parse(time.DateOnly, fmt.Sprintf("%d-%d-%d", req.Year, req.Month, req.Day))
	if err != nil {
		return resp, fmt.Errorf("CreateExpense: %v", err)
	}

	if !checkTagIdsValid(dbc, req.TagIDs) {
		return resp, fmt.Errorf("CreateExpense: one or more tagIds do not exist")
	}

	db := dbc.database

	// Try insert expense into db
	result, err := db.Exec("INSERT INTO expense (category_id, amount, date, location, notes) VALUES (?, ?, ?, ?, ?)", req.CategoryId, req.Amount, date, req.Location, req.Notes)
	if err != nil {
		return resp, fmt.Errorf("CreateExpense: %v", err)
	}

	// Try to fetch ID of expense
	id, err := result.LastInsertId()
	if err != nil {
		return resp, fmt.Errorf("CreateExpense: %v", err)
	}

	// Try creating all expense_tags
	didFailToAddTags := false
	for _, tagId := range req.TagIDs {
		_, err := db.Exec("INSERT INTO expense_tags (expense_id, tag_id) VALUES (?, ?)", id, tagId)
		if err != nil {
			log.Panicf("Failed to link TagID: %d to ExpenseID: %d", tagId, id)
			didFailToAddTags = true
		}
	}
	if didFailToAddTags {
		return resp, fmt.Errorf("CreateExpense: Failed to add one or more tags to expense object")
	}

	resp.ID = id
	return resp, nil
}

func (dbc *DBClient) GetExpenseById(req io.GetExpenseByIdRequest) (io.GetExpenseByIdResponse, error) {
	var resp io.GetExpenseByIdResponse

	// Check for valid ID
	if req.ID < 1 {
		return resp, fmt.Errorf("GetExpenseById: ID cannot be zero or negative")
	}

	db := dbc.database

	// Try get expense by ID
	var expense Expense
	row := db.QueryRow("SELECT * FROM category WHERE id = ?", req.ID)
	if err := row.Scan(&expense.ID, &expense.CategoryID, &expense.Amount, &expense.Date, &expense.Location, &expense.Notes); err != nil {
		if err == sql.ErrNoRows {
			return resp, fmt.Errorf("GetExpenseById: no tag with ID %d", req.ID)
		}
		return resp, fmt.Errorf("GetExpenseById: %v", err)
	}

	// Try get category name from ID
	var category io.GetCategoryByIdResponse
	if expense.CategoryID > 0 {
		cat, err := dbc.GetCategoryById(io.GetCategoryByIdRequest{ID: expense.CategoryID})
		if err != nil {
			return resp, fmt.Errorf("GetExpenseById: %v", err)
		}
		category = cat
	}

	tagNames, err := getTagNames(dbc, req.ID)
	if err != nil {
		return resp, nil
	}

	// Return category found
	resp.ID = expense.ID
	resp.CategoryName = category.Name
	resp.Day = int64(expense.Date.Day())
	resp.Month = int64(expense.Date.Month())
	resp.Year = int64(expense.Date.Year())
	resp.Location = expense.Location
	resp.Notes = expense.Notes
	resp.TagNames = tagNames
	return resp, nil
}

func (dbc *DBClient) GetTagsByExpenseId(req io.GetTagsByExpenseIdRequest) (io.GetTagsByExpenseIdResponse, error) {
	var resp io.GetTagsByExpenseIdResponse
	var tags []io.GetTagByIdResponse

	// Check for valid ID
	if req.ExpenseID < 1 {
		return resp, fmt.Errorf("GetTagsByExpenseId: ID cannot be zero or negative")
	}

	db := dbc.database

	// Try get tags by expense ID
	rows, err := db.Query("SELECT tag_id FROM expense_tags WHERE expense_id = ?", req.ExpenseID)
	if err != nil {
		return resp, fmt.Errorf("GetTagsByExpenseId: %v", err)
	}
	defer rows.Close()

	// Fetch tags for all tag IDs found
	for rows.Next() {
		var tagId int64
		if err := rows.Scan(&tagId); err != nil {
			return resp, fmt.Errorf("GetTagsByExpenseId: %v", err)
		}

		tag, err := dbc.GetTagById(io.GetTagByIdRequest{ID: tagId})
		if err != nil {
			return resp, fmt.Errorf("GetTagsByExpenseId: %v", err)
		}

		tags = append(tags, tag)
	}

	return resp, nil
}

func checkTagIdsValid(dbc *DBClient, tagIDs []int64) bool {
	for _, ID := range tagIDs {
		_, err := dbc.GetTagById(io.GetTagByIdRequest{ID: ID})
		if err != nil {
			return false
		}
	}

	return true
}

func getTagNames(dbc *DBClient, expenseId int64) ([]string, error) {
	var tagNames []string

	tags, err := dbc.GetTagsByExpenseId(io.GetTagsByExpenseIdRequest{ExpenseID: expenseId})
	if err != nil {
		return []string{}, fmt.Errorf("Failed to get tag names for expenseId: %d", expenseId)
	}

	for _, tag := range tags.Tags {
		tagNames = append(tagNames, tag.Name)
	}

	return tagNames, nil
}
