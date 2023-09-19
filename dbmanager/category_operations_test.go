package dbmanager

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	io "github.com/pinqy/buddy/inputoutput"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCategory(t *testing.T) {
	tests := []struct {
		name         string
		input        io.CreateCategoryRequest
		expectedResp io.CreateCategoryResponse
		expectedErr  error
	}{
		{
			name: "CreateCategory_InvalidName",
			input: io.CreateCategoryRequest{
				Name: "",
			},
			expectedResp: io.CreateCategoryResponse{},
			expectedErr:  fmt.Errorf("CreateCategory: cannot create category with no name"),
		},
		{
			name: "CreateCategory_Success",
			input: io.CreateCategoryRequest{
				Name:        "Test Name",
				Description: "Desc",
			},
			expectedResp: io.CreateCategoryResponse{
				ID: 1,
			},
			expectedErr: nil,
		},
		{
			name: "CreateCategory_ConnectionFailure",
			input: io.CreateCategoryRequest{
				Name:        "Test Name",
				Description: "Desc",
			},
			expectedResp: io.CreateCategoryResponse{},
			expectedErr:  fmt.Errorf("CreateCategory: %s", sql.ErrConnDone),
		},
		{
			name: "CreateCategory_LastInsertIdFailure",
			input: io.CreateCategoryRequest{
				Name:        "Test Name",
				Description: "Desc",
			},
			expectedResp: io.CreateCategoryResponse{},
			expectedErr:  fmt.Errorf("CreateCategory: %s", sql.ErrNoRows),
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	// mock for 2nd tc, successful category create
	mock.ExpectExec(`INSERT INTO category (name, description) VALUES (?, ?)`).
		WithArgs("Test Name", "Desc").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// mock for 3rd tc, connection failure
	mock.ExpectExec(`INSERT INTO category (name, description) VALUES (?, ?)`).
		WithArgs("Test Name", "Desc").
		WillReturnError(sql.ErrConnDone)

	// mock for 4th tc, last insert id failure
	mock.ExpectExec(`INSERT INTO category (name, description) VALUES (?, ?)`).
		WithArgs("Test Name", "Desc").
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

	for _, tc := range tests {
		dbc := DBClient{database: db}

		resp, err := dbc.CreateCategory(tc.input)

		assert.Equal(t, tc.expectedResp, resp)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func Test_GetCategoryById(t *testing.T) {
	tests := []struct {
		name         string
		input        io.GetCategoryByIdRequest
		expectedResp io.GetCategoryByIdResponse
		expectedErr  error
	}{
		{
			name: "GetCategoryById_InvalidId",
			input: io.GetCategoryByIdRequest{
				ID: 0,
			},
			expectedResp: io.GetCategoryByIdResponse{},
			expectedErr:  fmt.Errorf("GetCategoryById: ID cannot be zero or negative"),
		},
		{
			name: "GetCategoryById_Success",
			input: io.GetCategoryByIdRequest{
				ID: 1,
			},
			expectedResp: io.GetCategoryByIdResponse{
				ID:          1,
				Name:        "Test Name",
				Description: "Desc",
			},
			expectedErr: nil,
		},
		{
			name: "GetCategoryById_NoCategoryForId",
			input: io.GetCategoryByIdRequest{
				ID: 26,
			},
			expectedResp: io.GetCategoryByIdResponse{},
			expectedErr:  fmt.Errorf("GetCategoryById: no category with ID %d", 26),
		},
		{
			name: "GetCategoryById_ConnectionError",
			input: io.GetCategoryByIdRequest{
				ID: 1,
			},
			expectedResp: io.GetCategoryByIdResponse{},
			expectedErr:  fmt.Errorf("GetCategoryById: %s", sql.ErrConnDone),
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	// mock for 2nd tc, success
	mock.ExpectQuery(`SELECT * FROM category WHERE id = ?`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{`id`, `name`, `description`}).
			AddRow(int64(1), "Test Name", "Desc"))

	// mock for 3rd tc, no category for id
	mock.ExpectQuery(`SELECT * FROM category WHERE id = ?`).
		WithArgs(int64(26)).
		WillReturnRows(sqlmock.NewRows([]string{`id`, `name`, `description`}))

	// mock for 4th tc, connection failure
	mock.ExpectQuery(`SELECT * FROM category WHERE id = ?`).
		WithArgs(int64(1)).
		WillReturnError(sql.ErrConnDone)

	for _, tc := range tests {
		dbc := DBClient{database: db}

		resp, err := dbc.GetCategoryById(tc.input)

		assert.Equal(t, tc.expectedResp, resp)
		assert.Equal(t, tc.expectedErr, err)
	}
}
