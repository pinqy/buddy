package dbmanager

import (
	"database/sql"
	"fmt"

	io "github.com/pinqy/buddy/inputoutput"
)

func (dbc *DBClient) CreateCategory(req io.CreateCategoryRequest) (io.CreateCategoryResponse, error) {
	var resp io.CreateCategoryResponse

	// Check for valid category arguments
	if len(req.Name) < 1 {
		return resp, fmt.Errorf("CreateCategory: cannot create category with no name")
	}

	db := dbc.database

	// Try insert category into db
	result, err := db.Exec("INSERT INTO category (name, description) VALUES (?, ?)", req.Name, req.Description)
	if err != nil {
		return resp, fmt.Errorf("CreateCategory: %v", err)
	}

	// Try to fetch ID of inserted category
	id, err := result.LastInsertId()
	if err != nil {
		return resp, fmt.Errorf("CreateCategory: %v", err)
	}

	// Return ID of created category
	resp.ID = id
	return resp, nil
}

func (dbc *DBClient) GetCategoryById(req io.GetCategoryByIdRequest) (io.GetCategoryByIdResponse, error) {
	var resp io.GetCategoryByIdResponse

	// Check for valid ID
	if req.ID < 1 {
		return resp, fmt.Errorf("GetCategoryById: ID cannot be zero or negative")
	}

	db := dbc.database

	// Try get category by ID
	var category Category
	row := db.QueryRow("SELECT * FROM category WHERE id = ?", req.ID)
	if err := row.Scan(&category.ID, &category.Name, &category.Description); err != nil {
		if err == sql.ErrNoRows {
			return resp, fmt.Errorf("GetCategoryById: no category with ID %d", req.ID)
		}
		return resp, fmt.Errorf("GetCategoryById: %v", err)
	}

	// Return category found
	resp.ID = category.ID
	resp.Name = category.Name
	resp.Description = category.Description
	return resp, nil
}
