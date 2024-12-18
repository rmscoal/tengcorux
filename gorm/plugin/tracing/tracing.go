package tracing

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"

	"github.com/rmscoal/tengcorux/tracer"
	"github.com/rmscoal/tengcorux/tracer/attribute"
	"gorm.io/gorm"
)

// SpanNameGenerator is a function where user can customize their span name.
type SpanNameGenerator func(callbackName string, tx *gorm.DB) string

var _defaultSpanNameGenerator SpanNameGenerator = func(
	callbackName string,
	tx *gorm.DB,
) string {
	return callbackName
}

// tracing implements gorm.Plugin becoming
// main struct for this package.
type tracing struct {
	provider tracer.Tracer

	// configs
	showSQLVariable   bool
	spanNameGenerator SpanNameGenerator
}

// NewPlugin returns a tracing instance that could be called during
// the registration of gorm Writer Plugin. For example:
//
//	db.Use(tracing.New(tracing.WithTracer(someTracer)))
//
// This will register the tracing within gorm's callbacks.
func NewPlugin(opts ...Option) gorm.Plugin {
	t := &tracing{
		provider:          tracer.GetGlobalTracer(),
		spanNameGenerator: _defaultSpanNameGenerator,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Name returns the plugin name as required by gorm.Plugin.
func (t *tracing) Name() string {
	return "gorm:tracing"
}

// Initialize registers all the tracing callbacks to the GORM instance.
// These callbacks are invoked at different stages of the ORM operations (Query, Create, Update, etc.)
// to provide observability by registering the before and after callbacks.
func (t *tracing) Initialize(db *gorm.DB) error {
	var errs error

	// QUERY
	errs = errors.Join(errs,
		db.Callback().Query().Before("gorm:query").Register("tracing:before:query",
			t.before("SQL SELECT")))
	errs = errors.Join(errs,
		db.Callback().Query().After("gorm:query").Register("tracing:after:query",
			t.after()))

	// CREATE
	errs = errors.Join(errs,
		db.Callback().Create().Before("gorm:create").Register("tracing:before:create",
			t.before("SQL INSERT")))
	errs = errors.Join(errs,
		db.Callback().Create().After("gorm:create").Register("tracing:after:create",
			t.after()))

	// UPDATE
	errs = errors.Join(errs,
		db.Callback().Update().Before("gorm:update").Register("tracing:before:update",
			t.before("SQL UPDATE")))
	errs = errors.Join(errs,
		db.Callback().Update().After("gorm:update").Register("tracing:after:update",
			t.after()))

	// DELETE
	errs = errors.Join(errs,
		db.Callback().Delete().Before("gorm:delete").Register("tracing:before:delete",
			t.before("SQL DELETE")))
	errs = errors.Join(errs,
		db.Callback().Delete().After("gorm:delete").Register("tracing:after:delete",
			t.after()))

	// ROW
	errs = errors.Join(errs,
		db.Callback().Row().Before("gorm:row").Register("tracing:before:row",
			t.before("SQL ROW")))
	errs = errors.Join(errs,
		db.Callback().Row().After("gorm:row").Register("tracing:after:row",
			t.after()))

	// RAW
	errs = errors.Join(errs,
		db.Callback().Raw().Before("gorm:raw").Register("tracing:before:raw",
			t.before("SQL RAW")))
	errs = errors.Join(errs,
		db.Callback().Raw().After("gorm:raw").Register("tracing:after:raw",
			t.after()))

	return errs
}

// operationNameKey defines the key for operation name in which the value will be passed
// through context during the before and after calls.
var operationNameKey struct{}

func (t *tracing) before(operationName string) func(*gorm.DB) {
	return func(tx *gorm.DB) {
		tx.Statement.Context, _ = t.provider.StartSpan(
			context.WithValue(tx.Statement.Context, operationNameKey,
				operationName),
			t.spanNameGenerator(operationName, tx),
			tracer.WithSpanType(tracer.SpanTypeLocal),
			tracer.WithSpanLayer(tracer.SpanLayerDatabase),
		)
	}
}

func (t *tracing) after() func(*gorm.DB) {
	return func(tx *gorm.DB) {
		span := t.provider.SpanFromContext(tx.Statement.Context)
		if span == nil {
			return
		}
		defer span.End()

		// Injects the following attributes:
		// 1. Query
		// 2. Table
		// 3. DB Name
		// 4. Dialect
		// 5. Record error if there are any
		var dbStmtAttr attribute.KeyValue
		if t.showSQLVariable {
			dbStmtAttr = attribute.DBStatement(tx.Dialector.Explain(tx.Statement.SQL.String(),
				tx.Statement.Vars...))

		} else {
			dbStmtAttr = attribute.DBStatement(tx.Statement.SQL.String())
		}
		span.SetAttributes(
			dbStmtAttr,
			attribute.DBTable(tx.Statement.Table),
			attribute.DBName(tx.Name()),
			attribute.DBSystem(mapDBSystem(tx.Dialector.Name())),
			attribute.DBOperation(
				strings.ReplaceAll(
					tx.Statement.Context.Value(operationNameKey).(string),
					"SQL ",
					"")),
			attribute.KeyValuePair("gorm.plugin.package",
				t.Name()+"@"+t.Version()),
		)

		switch {
		case tx.Error == nil,
			tx.Error == io.EOF,
			errors.Is(tx.Error, gorm.ErrRecordNotFound),
			errors.Is(tx.Error, driver.ErrSkip),
			errors.Is(tx.Error, sql.ErrNoRows):
			// We ignore these errors
		default:
			span.RecordError(tx.Error)
		}
	}
}

const (
	MySQL      = "MySQL"
	MsSQL      = "Microsoft SQL Server"
	PostgreSQL = "PostgreSQL"
	SQLite     = "SQLite"
	SQLServer  = "SQL Server"
)

// mapDBSystem maps a know dialect that are often used with GORM. The provided
// dialects are "mysql", "mssql", "postgres", "postgresql", "sqlite", "sqlserver".
// Unknown dialects will return an empty string.
func mapDBSystem(name string) string {
	switch name {
	case "mysql":
		return MySQL
	case "mssql":
		return MsSQL
	case "postgres", "postgresql", "pgx":
		return PostgreSQL
	case "sqlite":
		return SQLite
	case "sqlserver":
		return SQLServer
	default:
		return ""
	}
}
