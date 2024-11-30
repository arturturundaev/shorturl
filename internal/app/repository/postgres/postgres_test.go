package postgres

import (
	"errors"
	"fmt"
	"github.com/arturturundaev/shorturl/internal/app/entity"
	"github.com/arturturundaev/shorturl/internal/app/handler/batch"
	"github.com/arturturundaev/shorturl/internal/app/utils"
	"github.com/gin-gonic/gin"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

func TestPostgresRepository_Ping(t *testing.T) {

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	db, _, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "Success",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepository{db}

			if err = repo.Ping(ctx); !errors.Is(tt.wantErr, err) {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepository_FindByShortURL(t *testing.T) {

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		shortURL string
		want     *entity.ShortURLEntity
		wantErr  error
	}{
		{
			name:     "SUCCESS fund 1 row",
			shortURL: "SUCCESS fund 1 row",
			want: &entity.ShortURLEntity{
				ShortURL:      "SUCCESS fund 1 row",
				URL:           "SUCCESS fund 1 row",
				CorrelationID: "SUCCESS fund 1 row",
				AddedUserID:   "SUCCESS fund 1 row",
				IsDeleted:     false,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepository{DB: db}
			rows := sqlxmock.NewRows([]string{"url_short", "original_url", "correlation_id", "added_user_id", "is_deleted"}).
				AddRow(tt.want.ShortURL, tt.want.URL, tt.want.CorrelationID, tt.want.AddedUserID, tt.want.IsDeleted)
			mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("select url_full as original_url, url_short, is_deleted from %s where url_short = $1", TableName))).WithArgs(tt.shortURL).WillReturnRows(rows)
			got, er := repo.FindByShortURL(tt.shortURL)
			if !errors.Is(er, tt.wantErr) {
				t.Errorf("FindByShortURL() error = %v, wantErr %v", er, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByShortURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepository_Batch(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		args    []batch.ButchRequest
		want    []entity.ShortURLEntity
		wantErr bool
	}{
		{
			name: "Success count 2",
			args: []batch.ButchRequest{
				{
					CorrelationID: "1",
					OriginalURL:   "1",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "2",
				},
			},
			want: []entity.ShortURLEntity{
				{
					ShortURL:      utils.GenerateShortURL("1"),
					URL:           "1",
					CorrelationID: "1",
					IsDeleted:     false,
				},
				{
					ShortURL:      utils.GenerateShortURL("2"),
					URL:           "2",
					CorrelationID: "2",
					IsDeleted:     false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(`INSERT into url (.+)`).
				WithArgs(sqlxmock.AnyArg(), sqlxmock.AnyArg(), sqlxmock.AnyArg(), sqlxmock.AnyArg(), sqlxmock.AnyArg(), sqlxmock.AnyArg(), sqlxmock.AnyArg(), sqlxmock.AnyArg()).
				WillReturnResult(sqlxmock.NewResult(1, 1))
			mock.ExpectCommit()
			repo := &PostgresRepository{DB: db}
			got, er := repo.Batch(tt.args)
			if (er != nil) != tt.wantErr {
				t.Errorf("Batch() error = %v, wantErr %v", er, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Batch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepository_GetUrlsByUserID(t *testing.T) {
	db, mock, er := sqlxmock.Newx()
	if er != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", er)
	}
	defer db.Close()

	tests := []struct {
		name    string
		userID  string
		rows    []entity.ShortURLEntity
		want    []entity.ShortURLEntity
		wantErr bool
	}{
		{
			name:    "success",
			userID:  "success",
			rows:    []entity.ShortURLEntity{{ShortURL: "success", URL: "success"}},
			want:    []entity.ShortURLEntity{{ShortURL: "success", URL: "success", CorrelationID: "", AddedUserID: "", IsDeleted: false}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepository{
				DB: db,
			}
			rows := sqlxmock.NewRows([]string{"original_url", "url_short"})
			for _, row := range tt.rows {
				rows.AddRow(row.ShortURL, row.URL)
			}
			mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("select url_full as original_url, url_short from %s where added_user_id = $1", TableName))).WithArgs(tt.userID).WillReturnRows(rows)

			got, err := repo.GetUrlsByUserID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUrlsByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUrlsByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepository_Delete(t *testing.T) {
	db, mock, er := sqlxmock.Newx()
	if er != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", er)
	}
	defer db.Close()

	type args struct {
		shortURLs   []string
		addedUserID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{shortURLs: []string{"a", "b"}, addedUserID: "1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepository{
				DB: db,
			}

			mock.ExpectExec(`(.+)`).
				WithArgs("1", "a", "b").
				WillReturnResult(sqlxmock.NewResult(0, 1))

			if err := repo.Delete(tt.args.shortURLs, tt.args.addedUserID); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
