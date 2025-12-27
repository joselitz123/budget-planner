package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// pgUUID converts a string to pgtype.UUID for database queries
func PgUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	if s != "" {
		_ = u.Scan(s)
	}
	return u
}

// pgUUIDPtr converts a string pointer to pgtype.UUID pointer
func PgUUIDPtr(s *string) pgtype.UUID {
	var u pgtype.UUID
	if s != nil && *s != "" {
		_ = u.Scan(*s)
	}
	return u
}

// pgDate converts a time.Time to pgtype.Date
func PgDate(t interface{}) pgtype.Date {
	var d pgtype.Date
	_ = d.Scan(t)
	return d
}

// numericToFloat64 converts pgtype.Numeric to float64
func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	// Use the underlying float representation
	f, _ := n.Float64Value()
	return f.Float64
}

// pgNumeric converts a float64 to pgtype.Numeric
func PgNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	// Convert float64 to string first, then scan into Numeric
	// This is more reliable than direct float scanning
	_ = n.Scan(strconv.FormatFloat(f, 'f', -1, 64))
	return n
}

// pgNumericPtr converts *float64 to pgtype.Numeric
func PgNumericPtr(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(*f, 'f', -1, 64))
	return n
}

// uuidToString converts pgtype.UUID to string
func UUIDToString(u pgtype.UUID) string {
	if u.Valid {
		// Format as standard UUID with hyphens: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
			u.Bytes[0], u.Bytes[1], u.Bytes[2], u.Bytes[3],
			u.Bytes[4], u.Bytes[5],
			u.Bytes[6], u.Bytes[7],
			u.Bytes[8], u.Bytes[9],
			u.Bytes[10], u.Bytes[11], u.Bytes[12], u.Bytes[13], u.Bytes[14], u.Bytes[15])
	}
	return ""
}

// pgText converts a string to pgtype.Text
func PgText(s string) pgtype.Text {
	var t pgtype.Text
	if s != "" {
		_ = t.Scan(s)
	}
	return t
}

// textToString converts pgtype.Text to string
func TextToString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

// textToStringPtr converts pgtype.Text to *string
func TextToStringPtr(t pgtype.Text) *string {
	if t.Valid {
		return &t.String
	}
	return nil
}

// numericToFloat64Ptr converts pgtype.Numeric to *float64
func NumericToFloat64Ptr(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, _ := n.Float64Value()
	return &f.Float64
}

// pgTextPtr converts *string to pgtype.Text
func PgTextPtr(s *string) pgtype.Text {
	var t pgtype.Text
	if s != nil {
		_ = t.Scan(*s)
	}
	return t
}

// dateToTime converts pgtype.Date to time.Time
func DateToTime(d pgtype.Date) time.Time {
	if d.Valid {
		return d.Time
	}
	return time.Time{}
}

// timestamptzToTime converts pgtype.Timestamptz to time.Time
func TimestamptzToTime(t pgtype.Timestamptz) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

// pgDatePtr converts *time.Time to pgtype.Date
func PgDatePtr(t *time.Time) pgtype.Date {
	var d pgtype.Date
	if t != nil {
		_ = d.Scan(*t)
	}
	return d
}

// pgBool converts bool to pgtype.Bool
func PgBool(b bool) pgtype.Bool {
	var bo pgtype.Bool
	_ = bo.Scan(b)
	return bo
}

// pgBoolPtr converts *bool to pgtype.Bool
func PgBoolPtr(b *bool) pgtype.Bool {
	var bo pgtype.Bool
	if b != nil {
		_ = bo.Scan(*b)
	}
	return bo
}

// pgInt4 converts int32 to pgtype.Int4
func PgInt4(i int32) pgtype.Int4 {
	var in pgtype.Int4
	_ = in.Scan(i)
	return in
}

// pgInt4Ptr converts *int32 to pgtype.Int4
func PgInt4Ptr(i *int32) pgtype.Int4 {
	var in pgtype.Int4
	if i != nil {
		_ = in.Scan(*i)
	}
	return in
}

// int4ToInt32 converts pgtype.Int4 to int32
func Int4ToInt32(i pgtype.Int4) *int32 {
	if i.Valid {
		return &i.Int32
	}
	return nil
}

// pgTimestamptz converts time.Time to pgtype.Timestamptz
func PgTimestamptz(t time.Time) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	_ = ts.Scan(t)
	return ts
}
