package dbmanager

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	io "github.com/pinqy/buddy/inputoutput"
)

func Test_CreateTag(t *testing.T) {
	tests := []struct {
		name         string
		input        io.CreateTagRequest
		expectedResp io.CreateTagResponse
		expectedErr  error
	}{
		{
			name: "CreateTag_Success",
			input: io.CreateTagRequest{
				Name: "Test Name",
			},
			expectedResp: io.CreateTagResponse{
				ID: 1,
			},
			expectedErr: nil,
		},
		{
			name: "CreateTag_NoName",
			input: io.CreateTagRequest{
				Name: "",
			},
			expectedResp: io.CreateTagResponse{},
			expectedErr:  fmt.Errorf("CreateTag: cannot create tag with no name"),
		},
		{
			name: "CreateTag_ConnectionFailure",
			input: io.CreateTagRequest{
				Name: "Test Name",
			},
			expectedResp: io.CreateTagResponse{},
			expectedErr:  fmt.Errorf("CreateTag: %s", sql.ErrConnDone),
		},
		{
			name: "CreateTag_LastInsertIdFailure",
			input: io.CreateTagRequest{
				Name: "Test Name",
			},
			expectedResp: io.CreateTagResponse{},
			expectedErr:  fmt.Errorf("CreateTag: %s", sql.ErrNoRows),
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	// mock for 1st tc, successful tag creation
	mock.ExpectExec(`INSERT INTO tag (name) VALUES (?)`).
		WithArgs("Test Name").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// mock for 3rd tc, fail on query
	mock.ExpectExec(`INSERT INTO tag (name) VALUES (?)`).
		WithArgs("Test Name").
		WillReturnError(sql.ErrConnDone)

	mock.ExpectExec(`INSERT INTO tag (name) VALUES (?)`).
		WithArgs("Test Name").
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

	for _, tc := range tests {
		dbc := DBClient{database: db}

		resp, err := dbc.CreateTag(tc.input)

		assert.Equal(t, tc.expectedResp, resp)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func Test_GetTagById(t *testing.T) {
	tests := []struct {
		name         string
		input        io.GetTagByIdRequest
		expectedResp io.GetTagByIdResponse
		expectedErr  error
	}{
		{
			name: "GetTagById_InvalidId",
			input: io.GetTagByIdRequest{
				ID: 0,
			},
			expectedResp: io.GetTagByIdResponse{},
			expectedErr:  fmt.Errorf("GetTagById: ID cannot be zero or negative"),
		},
		{
			name: "GetTagById_Success",
			input: io.GetTagByIdRequest{
				ID: 1,
			},
			expectedResp: io.GetTagByIdResponse{
				ID:   1,
				Name: "Test Name",
			},
			expectedErr: nil,
		},
		{
			name: "GetTagById_NoTagForId",
			input: io.GetTagByIdRequest{
				ID: 26,
			},
			expectedResp: io.GetTagByIdResponse{},
			expectedErr:  fmt.Errorf("GetTagById: no tag with ID 26"),
		},
		{
			name: "GetTagById_ConnectionError",
			input: io.GetTagByIdRequest{
				ID: 1,
			},
			expectedResp: io.GetTagByIdResponse{},
			expectedErr:  fmt.Errorf("GetTagById: %s", sql.ErrConnDone),
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock", err)
	}
	defer db.Close()

	// mock for 2nd tc, successful query
	mock.ExpectQuery(`SELECT * FROM tag WHERE id = ?`).
		WithArgs(int64(1)).
		WillReturnRows(
			sqlmock.NewRows([]string{`id`, `name`}).AddRow(int64(1), "Test Name"))

	// mock for 3rd tc, no tag found
	mock.ExpectQuery(`SELECT * FROM tag WHERE id = ?`).
		WithArgs(int64(26)).
		WillReturnRows(sqlmock.NewRows([]string{`id`, `name`}))

	// mock for 4th tc, query failed
	mock.ExpectQuery(`SELECT * FROM tag WHERE id = ?`).
		WithArgs(int64(1)).
		WillReturnError(sql.ErrConnDone)

	for _, tc := range tests {
		dbc := DBClient{database: db}

		resp, err := dbc.GetTagById(tc.input)

		assert.Equal(t, tc.expectedResp, resp)
		assert.Equal(t, tc.expectedErr, err)
	}
}
