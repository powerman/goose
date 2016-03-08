package goose

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"
)

func Status(db *sql.DB, dir string) error {
	// collect all migrations
	min := int64(0)
	max := int64((1 << 63) - 1)
	migrations, err := CollectMigrations(dir, min, max)
	if err != nil {
		return err
	}

	ms := migrationSorter(migrations)
	ms.Sort(0)

	// must ensure that the version table exists if we're running on a pristine DB
	if _, err := EnsureDBVersion(db); err != nil {
		return err
	}

	fmt.Println("goose: status")
	fmt.Println("    Applied At                  Migration")
	fmt.Println("    =======================================")
	for _, m := range ms {
		printMigrationStatus(db, m.Version, filepath.Base(m.Source))
	}

	return nil
}

func printMigrationStatus(db *sql.DB, version int64, script string) {
	var row MigrationRecord
	q := fmt.Sprintf("SELECT tstamp, is_applied FROM goose_db_version WHERE version_id=%d ORDER BY tstamp DESC LIMIT 1", version)
	e := db.QueryRow(q).Scan(&row.TStamp, &row.IsApplied)

	if e != nil && e != sql.ErrNoRows {
		log.Fatal(e)
	}

	var appliedAt string

	if row.IsApplied {
		appliedAt = row.TStamp.Format(time.ANSIC)
	} else {
		appliedAt = "Pending"
	}

	fmt.Printf("    %-24s -- %v\n", appliedAt, script)
}