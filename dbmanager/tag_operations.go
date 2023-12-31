package dbmanager

import (
	"database/sql"
	"fmt"

	io "github.com/pinqy/buddy/inputoutput"
)

func (dbc *DBClient) CreateTag(req io.CreateTagRequest) (io.CreateTagResponse, error) {
	var resp io.CreateTagResponse

	// Check for valid tag name
	if len(req.Name) == 0 {
		return resp, fmt.Errorf("CreateTag: cannot create tag with no name")
	}

	db := dbc.database

	// Try insert tag into db
	result, err := db.Exec("INSERT INTO tag (name) VALUES (?)", req.Name)
	if err != nil {
		return resp, fmt.Errorf("CreateTag: %v", err)
	}

	// Try to fetch ID of inserted tag
	id, err := result.LastInsertId()
	if err != nil {
		return resp, fmt.Errorf("CreateTag: %v", err)
	}

	// Return ID of created tag
	resp.ID = id
	return resp, nil
}

func (dbc *DBClient) GetTagById(req io.GetTagByIdRequest) (io.GetTagByIdResponse, error) {
	var resp io.GetTagByIdResponse

	// Check for valid ID
	if req.ID < 1 {
		return resp, fmt.Errorf("GetTagById: ID cannot be zero or negative")
	}

	db := dbc.database

	// Try get tag by ID
	var tag Tag
	row := db.QueryRow("SELECT * FROM tag WHERE id = ?", req.ID)
	if err := row.Scan(&tag.ID, &tag.Name); err != nil {
		if err == sql.ErrNoRows {
			return resp, fmt.Errorf("GetTagById: no tag with ID %d", req.ID)
		}
		return resp, fmt.Errorf("GetTagById: %v", err)
	}

	// Return tag found
	resp.ID = tag.ID
	resp.Name = tag.Name
	return resp, nil
}
