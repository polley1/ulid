// +build cockroach

package ulid_test

import (
	"testing"
	"time"

	"github.com/polley1/ulid/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestULIDModel is a simple model to verify ULID behavior with UUID columns
type TestULIDModel struct {
	ID        ulid.ULID `gorm:"type:uuid;primaryKey"`
	Name      string
	CreatedAt time.Time
}

func TestVerifyPackageLink(t *testing.T) {
	var u ulid.ULID
	// Pass something that triggers the print (byte slice)
	// it will likely fail parsing, but should print
	u.Scan([]byte("1234567890123456"))
}

func TestCockroachULID(t *testing.T) {
	// Connection details from config.yaml
	dsn := "host=192.168.200.73 port=26257 user=root password=JacobQ98 dbname=latix sslmode=disable"

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate the test model
	// This will create the table "test_ulid_models" with a 'uuid' column for ID
	if err := db.AutoMigrate(&TestULIDModel{}); err != nil {
		t.Fatalf("Failed to AutoMigrate: %v", err)
	}

	// Clean up previous runs
	db.Exec("DELETE FROM test_ulid_models")

	// Generate a new ULID
	id := ulid.Make()
	t.Logf("Generated ULID: %s", id.String())

	// Create a record
	record := TestULIDModel{
		ID:   id,
		Name: "Test Record",
	}

	// Save to DB
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}
	t.Log("Successfully verified creation")

	// Retrieve from DB
	var loaded TestULIDModel
	if err := db.First(&loaded, "id = ?", id).Error; err != nil {
		t.Fatalf("Failed to retrieve record: %v", err)
	}

	// Verify ID matches
	if loaded.ID != id {
		t.Errorf("Loaded ID mismatch. Expected %s, got %s", id.String(), loaded.ID.String())
	} else {
		t.Log("Successfully verified load and match")
	}

	// Verify Name matches
	if loaded.Name != "Test Record" {
		t.Errorf("Loaded Name mismatch. Expected 'Test Record', got '%s'", loaded.Name)
	}

	// Clean up
	// db.Migrator().DropTable(&TestULIDModel{})
}
