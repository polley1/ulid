package ulid_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/polley1/ulid/v2"
	"gorm.io/gorm"
)

// Mock Dialector for testing GORM data types
type mockDialector struct {
	name string
	gorm.Dialector
}

func (m mockDialector) Name() string {
	return m.name
}

func TestNullableULID_Scan(t *testing.T) {
	id := ulid.MustNew(123, nil)
	validIDString := id.String()
	validIDBytes := id[:]

	tests := []struct {
		name    string
		input   interface{}
		want    ulid.NullableULID
		wantErr bool
	}{
		{
			name:  "nil input",
			input: nil,
			want:  ulid.NullableULID{Valid: false},
		},
		{
			name:  "valid string",
			input: validIDString,
			want:  ulid.NullableULID{ULID: id, Valid: true},
		},
		{
			name:  "valid bytes",
			input: validIDBytes,
			want:  ulid.NullableULID{ULID: id, Valid: true},
		},
		{
			name:  "valid bytes string representation",
			input: []byte(validIDString),
			want:  ulid.NullableULID{ULID: id, Valid: true},
		},
		{
			name:    "invalid string",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "invalid bytes",
			input:   []byte("invalid"),
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nu ulid.NullableULID
			err := nu.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullableULID.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if nu.Valid != tt.want.Valid {
					t.Errorf("NullableULID.Scan() valid = %v, want %v", nu.Valid, tt.want.Valid)
				}
				if nu.Valid && nu.ULID != tt.want.ULID {
					t.Errorf("NullableULID.Scan() ulid = %v, want %v", nu.ULID, tt.want.ULID)
				}
			}
		})
	}
}

func TestNullableULID_Value(t *testing.T) {
	id := ulid.MustNew(123, nil)

	tests := []struct {
		name    string
		nu      ulid.NullableULID
		want    driver.Value
		wantErr bool
	}{
		{
			name: "valid null ULID",
			nu:   ulid.NullableULID{Valid: false},
			want: nil,
		},
		{
			name: "valid non-null ULID",
			nu:   ulid.NullableULID{ULID: id, Valid: true},
			want: id.String(), // NullableULID.Value currently returns string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nu.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullableULID.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NullableULID.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableULID_GormDataType(t *testing.T) {
	var nu ulid.NullableULID
	if got := nu.GormDataType(); got != "uuid" {
		t.Errorf("NullableULID.GormDataType() = %v, want %v", got, "uuid")
	}
}

func TestNullableULID_GormDBDataType(t *testing.T) {
	tests := []struct {
		name      string
		dialector string
		want      string
	}{
		{"mysql", "mysql", "VARBINARY(16)"},
		{"postgres", "postgres", "uuid"},
		{"sqlite", "sqlite", "BLOB"},
		{"other", "sqlserver", "uuid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &gorm.DB{Config: &gorm.Config{Dialector: &mockDialector{name: tt.dialector}}}
			var nu ulid.NullableULID
			if got := nu.GormDBDataType(db, nil); got != tt.want {
				t.Errorf("NullableULID.GormDBDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanULID(t *testing.T) {
	id := ulid.MustNew(123, nil)
	validIDString := id.String()
	validIDBytes := id[:]

	tests := []struct {
		name    string
		input   interface{}
		want    ulid.ULID
		wantErr bool
	}{
		{
			name:  "nil input",
			input: nil,
			want:  ulid.ULID{},
		},
		{
			name:  "valid string",
			input: validIDString,
			want:  id,
		},
		{
			name:  "valid bytes",
			input: validIDBytes,
			want:  id,
		},
		{
			name:  "valid bytes string representation",
			input: []byte(validIDString),
			want:  id,
		},
		{
			name:    "invalid string",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "invalid bytes",
			input:   []byte("invalid"),
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ulid.ScanULID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanULID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ScanULID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueULID(t *testing.T) {
	id := ulid.MustNew(123, nil)

	tests := []struct {
		name    string
		u       ulid.ULID
		want    driver.Value
		wantErr bool
	}{
		{
			name: "zero ULID",
			u:    ulid.ULID{},
			want: nil,
		},
		{
			name: "non-zero ULID",
			u:    id,
			want: id[:], // ValueULID returns bytes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ulid.ValueULID(tt.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValueULID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			match := false
			if got == nil && tt.want == nil {
				match = true
			} else if gotBytes, ok := got.([]byte); ok {
				if wantBytes, ok := tt.want.([]byte); ok {
					if string(gotBytes) == string(wantBytes) {
						match = true
					}
				}
			}

			if !match {
				t.Errorf("ValueULID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestULIDPlaceholders(t *testing.T) {
	ids := []ulid.ULID{
		ulid.MustNew(1, nil),
		ulid.MustNew(2, nil),
		ulid.MustNew(3, nil),
	}

	tests := []struct {
		name   string
		dbType string
		want   string
	}{
		{"mysql", "mysql", "?,?,?"},
		{"postgres", "postgres", "$1,$2,$3"},
		{"sqlite", "sqlite", "?,?,?"},
		{"default", "other", "?,?,?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ulid.SetDBType(tt.dbType)
			args, placeholders := ulid.ULIDPlaceholders(ids)

			if len(args) != 3 {
				t.Errorf("ULIDPlaceholders() returned %d args, want 3", len(args))
			}

			if placeholders != tt.want {
				t.Errorf("ULIDPlaceholders() placeholders = %q, want %q", placeholders, tt.want)
			}
		})
	}
}

// Ensure NullableULID implements Scanner and Valuer
var _ sql.Scanner = (*ulid.NullableULID)(nil)
var _ driver.Valuer = (*ulid.NullableULID)(nil)
