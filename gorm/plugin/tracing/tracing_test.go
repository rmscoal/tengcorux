package tracing

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/rmscoal/tengcorux/tracer/attribute"
	"github.com/rmscoal/tengcorux/tracer/tracetest"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type userModel struct {
	ID          int
	Name        string
	PhoneNumber string
}

func (userModel) TableName() string {
	return "users"
}

type testCases struct {
	name          string
	extraOptions  []Option
	queryFunc     func(ctx context.Context, db *gorm.DB) error
	expectedSpans func(t *testing.T, spans []*tracetest.ReadOnlySpan)
}

func TestTracing(t *testing.T) {
	tests := []testCases{
		{
			name:         "QUERY",
			extraOptions: []Option{WithSQLVariables()},
			queryFunc: func(ctx context.Context, db *gorm.DB) error {
				var user userModel
				return db.WithContext(ctx).
					Table("users").
					Select("id, name").
					Where("phone_number = ?", 1234).
					Take(&user).Error
			},
			expectedSpans: func(t *testing.T, spans []*tracetest.ReadOnlySpan) {
				if len(spans) != 1 {
					t.Fatalf("want 1 span, got %d", len(spans))
				}

				firstSpan := spans[0]
				for _, attr := range firstSpan.Attributes {
					val := attr.Value.(string)
					switch attr.Key {
					case attribute.DBSystemKey:
						if val != SQLite {
							t.Fatalf("want %q, got %q", SQLite, val)
						}
					case attribute.DBStatementKey:
						if !strings.Contains(val,
							"SELECT id, name FROM `users`") {
							t.Fatalf("unexpected query statement, got %s", val)
						}
						if !strings.Contains(val, "phone_number = 1234") {
							t.Fatalf("WithSQLVariables option was passed, expected vars to be included, got %s",
								val)
						}
					case attribute.DBTableKey:
						if val != "users" {
							t.Fatalf("expected a table users but got %s", val)
						}
					case attribute.DBOperationKey:
						if val != "SELECT" {
							t.Fatalf("expected a SELECT operation but got %s",
								val)
						}
					}
				}
			},
		},
		{
			name: "CREATE",
			queryFunc: func(ctx context.Context, db *gorm.DB) error {
				user := userModel{ID: 1, Name: "test"}
				return db.WithContext(ctx).Create(&user).Error
			},
			expectedSpans: func(t *testing.T, spans []*tracetest.ReadOnlySpan) {
				if len(spans) != 1 {
					t.Fatalf("want 1 span, got %d", len(spans))
				}

				firstSpan := spans[0]
				for _, attr := range firstSpan.Attributes {
					val := attr.Value.(string)
					switch attr.Key {
					case attribute.DBSystemKey:
						if val != SQLite {
							t.Fatalf("want %q, got %q", SQLite, val)
						}
					case attribute.DBStatementKey:
						if !strings.Contains(val, "INSERT INTO `users`") {
							t.Fatalf("unexpected query statement, got %s", val)
						}
					case attribute.DBTableKey:
						if val != "users" {
							t.Fatalf("expected a table users but got %s", val)
						}
					case attribute.DBOperationKey:
						if val != "INSERT" {
							t.Fatalf("expected a INSERT operation but got %s",
								val)
						}
					}
				}
			},
		},
		{
			name: "UPDATE",
			queryFunc: func(ctx context.Context, db *gorm.DB) error {
				user := userModel{ID: 1, Name: "test"}
				return db.WithContext(ctx).Updates(&user).Error
			},
			expectedSpans: func(t *testing.T, spans []*tracetest.ReadOnlySpan) {
				if len(spans) != 1 {
					t.Fatalf("want 1 span, got %d", len(spans))
				}

				firstSpan := spans[0]
				for _, attr := range firstSpan.Attributes {
					val := attr.Value.(string)
					switch attr.Key {
					case attribute.DBSystemKey:
						if val != SQLite {
							t.Fatalf("want %q, got %q", SQLite, val)
						}
					case attribute.DBStatementKey:
						if !strings.Contains(val, "UPDATE `users` SET") {
							t.Fatalf("unexpected query statement, got %s", val)
						}
						if strings.Contains(val, "WHERE `id` = 1") {
							t.Fatalf("WithSQLVariables option was not passed, expected vars to not be included, got %s",
								val)
						}
					case attribute.DBTableKey:
						if val != "users" {
							t.Fatalf("expected a table users but got %s", val)
						}
					case attribute.DBOperationKey:
						if val != "UPDATE" {
							t.Fatalf("expected a UPDATE operation but got %s",
								val)
						}
					}
				}
			},
		},
		{
			name: "DELETE",
			queryFunc: func(ctx context.Context, db *gorm.DB) error {
				user := userModel{ID: 1, Name: "test"}
				return db.WithContext(ctx).Delete(&user).Error
			},
			expectedSpans: func(t *testing.T, spans []*tracetest.ReadOnlySpan) {
				if len(spans) != 1 {
					t.Fatalf("want 1 span, got %d", len(spans))
				}

				firstSpan := spans[0]
				for _, attr := range firstSpan.Attributes {
					val := attr.Value.(string)
					switch attr.Key {
					case attribute.DBSystemKey:
						if val != SQLite {
							t.Fatalf("want %q, got %q", SQLite, val)
						}
					case attribute.DBStatementKey:
						if !strings.Contains(val, "DELETE FROM `users`") {
							t.Fatalf("unexpected query statement, got %s", val)
						}
					case attribute.DBTableKey:
						if val != "users" {
							t.Fatalf("expected a table users but got %s", val)
						}
					case attribute.DBOperationKey:
						if val != "DELETE" {
							t.Fatalf("expected a DELETE operation but got %s",
								val)
						}
					}
				}
			},
		},
		{
			name: "ROW",
			queryFunc: func(ctx context.Context, db *gorm.DB) error {
				var num int
				return db.WithContext(ctx).Raw("SELECT 12").Scan(&num).Error
			},
			expectedSpans: func(t *testing.T, spans []*tracetest.ReadOnlySpan) {
				if len(spans) != 1 {
					t.Fatalf("want 1 span, got %d", len(spans))
				}

				firstSpan := spans[0]
				for _, attr := range firstSpan.Attributes {
					val := attr.Value.(string)
					switch attr.Key {
					case attribute.DBSystemKey:
						if val != SQLite {
							t.Fatalf("want %q, got %q", SQLite, val)
						}
					case attribute.DBStatementKey:
						if !strings.Contains(val, "SELECT 12") {
							t.Fatalf("unexpected query statement, got %s", val)
						}
					case attribute.DBOperationKey:
						if val != "ROW" {
							t.Fatalf("expected a ROW operation but got %s", val)
						}
					}
				}
			},
		},
		{
			name: "RAW",
			queryFunc: func(ctx context.Context, db *gorm.DB) error {
				return db.WithContext(ctx).Exec("CREATE TABLE foo (id int)").Error
			},
			expectedSpans: func(t *testing.T, spans []*tracetest.ReadOnlySpan) {
				if len(spans) != 1 {
					t.Fatalf("want 1 span, got %d", len(spans))
				}

				firstSpan := spans[0]
				for _, attr := range firstSpan.Attributes {
					val := attr.Value.(string)
					switch attr.Key {
					case attribute.DBSystemKey:
						if val != SQLite {
							t.Fatalf("want %q, got %q", SQLite, val)
						}
					case attribute.DBStatementKey:
						if !strings.Contains(val, "CREATE TABLE foo (id int)") {
							t.Fatalf("unexpected query statement, got %s", val)
						}
					case attribute.DBOperationKey:
						if val != "RAW" {
							t.Fatalf("expected a RAW operation but got %s", val)
						}
					}
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setting up DB, tracing, registration of plugin and check
			// whether the plugin has been registered successfully.
			db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"),
				&gorm.Config{})
			if err != nil {
				t.Fatalf("failed to connect database: %v", err)
			}
			t.Cleanup(func() { // Close DB during cleanup
				sqlDB, err := db.DB()
				if err != nil {
					t.Fatalf("failed to connect database: %v", err)
				}
				_ = sqlDB.Close()
			})

			err = db.AutoMigrate(&userModel{})
			if err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			db = db.Debug() // Set to debug mode
			tracer := tracetest.NewTracer()
			test.extraOptions = append(test.extraOptions, WithTracer(tracer))
			err = db.Use(NewPlugin(test.extraOptions...))
			if err != nil {
				t.Fatalf("failed to register tracing gorm middleware: %v", err)
			}

			plugin := db.Config.Plugins["gorm:tracing"]
			if plugin == nil {
				t.Fatalf("failed to find tracing gorm middleware in gorm plugins")
			}

			// Start the test loop
			err = test.queryFunc(context.Background(), db)
			checkGormError(t, err)

			endedSpans := tracer.Recorder().EndedSpans()
			test.expectedSpans(t, endedSpans)
		})
	}
}

func TestMapDBSystem(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "mysql",
			want:  MySQL,
		},
		{
			input: "mssql",
			want:  MsSQL,
		},
		{
			input: "postgres",
			want:  PostgreSQL,
		},
		{
			input: "postgresql",
			want:  PostgreSQL,
		},
		{
			input: "pgx",
			want:  PostgreSQL,
		},
		{
			input: "sqlite",
			want:  SQLite,
		},
		{
			input: "sqlserver",
			want:  SQLServer,
		},
		{
			input: "someUnknownDialect",
			want:  "",
		},
	}

	for _, test := range tests {
		got := mapDBSystem(test.input)
		if got != test.want {
			t.Errorf("expected %s but got %s", test.want, got)
			t.FailNow()
		}
	}
}

func TestTracingOption(t *testing.T) {
	t.Run("WithTracer", func(t *testing.T) {
		plugin := NewPlugin(WithTracer(nil))
		tr := plugin.(*tracing)
		if tr.provider != nil {
			t.Fatalf("expected nil tracing provider")
		}
	})

	t.Run("WithSQLVariables", func(t *testing.T) {
		plugin := NewPlugin(WithSQLVariables())
		tr := plugin.(*tracing)
		if tr.showSQLVariable == false {
			t.Fatalf("expected showSQLVariable to be true")
		}
	})

	t.Run("WithSpanNameGenerator", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"),
			&gorm.Config{})
		if err != nil {
			t.Fatalf("failed to connect database: %v", err)
		}
		defer func() { // Close DB during cleanup
			sqlDB, err := db.DB()
			if err != nil {
				t.Fatalf("failed to connect database: %v", err)
			}
			_ = sqlDB.Close()
		}()

		db = db.Debug() // Set to debug mode
		tracer := tracetest.NewTracer()
		err = db.Use(NewPlugin(
			WithSpanNameGenerator(func(
				callbackName string, tx *gorm.DB,
			) string {
				return "hello_world"
			}),
			WithTracer(tracer),
		))

		var count int
		err = db.WithContext(context.TODO()).Raw("SELECT 1").Scan(&count).Error
		if err != nil {
			t.Fatalf("failed to query database: %v", err)
		}
		spans := tracer.Recorder().EndedSpans()
		if len(spans) != 1 {
			t.Fatalf("want 1 span, got %d", len(spans))
		}
		if spans[0].Name != "hello_world" {
			t.Errorf("want %q, got %q", "hello_world", spans[0].Name)
		}
	})
}

func checkGormError(t *testing.T, err error) {
	switch {
	case err == nil,
		err == io.EOF,
		errors.Is(err, gorm.ErrRecordNotFound),
		errors.Is(err, driver.ErrSkip),
		errors.Is(err, sql.ErrNoRows):
	default:
		t.Fatalf("failed to execute query: %v", err)
	}
}
