package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/semver/v3"

	pb "gomodules.avm99963.com/twpt-server/api_proto"
)

func isValidSemVersion(v string) bool {
	_, err := semver.NewVersion(v)
	return err == nil
}

func fillRawKillSwitch(db *sql.DB, ctx context.Context, k *pb.KillSwitch, featureId int32) error {
	f, err := GetFeatureById(db, ctx, featureId)
	if err != nil {
		return fmt.Errorf("fillRawKillSwitch: %v", err)
	}
	if f == nil {
		return fmt.Errorf("fillRawKillSwitch: linked feature doesn't exist.")
	}
	k.Feature = f

	rows, err := db.QueryContext(ctx, "SELECT browser FROM KillSwitch2Browser WHERE kswitch_id = ?", k.Id)
	if err != nil {
		return fmt.Errorf("fillRawKillSwitch: %v", err)
	}
	defer rows.Close()

	var browsers []pb.Environment_Browser
	for rows.Next() {
		var b pb.Environment_Browser
		if err := rows.Scan(&b); err != nil {
			return fmt.Errorf("fillRawKillSwitch: %v", err)
		}
		browsers = append(browsers, b)
	}
	k.Browsers = browsers

	return nil
}

func GetKillSwitchById(db *sql.DB, ctx context.Context, id int32) (*pb.KillSwitch, error) {
	query := db.QueryRowContext(ctx, "SELECT kswitch_id, feat_id, min_version, max_version, active FROM KillSwitch WHERE kswitch_id = ?", id)
	var featureId int32
	var k pb.KillSwitch
	if err := query.Scan(&k.Id, &featureId, &k.MinVersion, &k.MaxVersion, &k.Active); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetKillSwitchById: %v.", err)
	}

	err := fillRawKillSwitch(db, ctx, &k, featureId)
	if err != nil {
		return nil, fmt.Errorf("GetKillSwitchById: $v.", err)
	}

	return &k, nil
}

func ListKillSwitches(db *sql.DB, ctx context.Context) ([]*pb.KillSwitch, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.QueryContext(ctx, "SELECT kswitch_id, feat_id, min_version, max_version, active FROM KillSwitch")
	if err != nil {
		return nil, fmt.Errorf("ListKillSwitches: ", err)
	}
	defer rows.Close()

	var killSwitches []*pb.KillSwitch
	for rows.Next() {
		var featureId int32
		var k pb.KillSwitch
		if err := rows.Scan(&k.Id, &featureId, &k.MinVersion, &k.MaxVersion, &k.Active); err != nil {
			return nil, fmt.Errorf("ListKillSwitches: %v.", err)
		}
		err := fillRawKillSwitch(db, ctx, &k, featureId)
		if err != nil {
			return nil, fmt.Errorf("ListKillSwitches: $v.", err)
		}
		killSwitches = append(killSwitches, &k)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListKillSwitches: ", err)
	}
	return killSwitches, nil
}

func EnableKillSwitch(db *sql.DB, ctx context.Context, k *pb.KillSwitch) error {
	if k.MinVersion != "" && !isValidSemVersion(k.MinVersion) {
		return fmt.Errorf("min_version is not a valid semantic version.")
	}
	if k.MaxVersion != "" && !isValidSemVersion(k.MaxVersion) {
		return fmt.Errorf("max_version is not a valid semantic version.")
	}
	if k.MinVersion != "" && k.MaxVersion != "" {
		minVersion, _ := semver.NewVersion(k.MinVersion)
		maxVersion, _ := semver.NewVersion(k.MaxVersion)
		if minVersion.GreaterThan(maxVersion) {
			return fmt.Errorf("min_version must be less than max_version.")
		}
	}
	for _, b := range k.GetBrowsers() {
		if b == pb.Environment_BROWSER_UNKNOWN {
			return fmt.Errorf("browsers cannot contain BROWSER_UNKNOWN.")
		}
	}

	k.Active = true

	f, err := GetFeatureById(db, ctx, k.Feature.Id)
	if err != nil {
		return fmt.Errorf("EnableKillSwitch: %v", err)
	}
	if f == nil {
		return fmt.Errorf("EnableKillSwitch: this feature doesn't exist.")
	}
	k.Feature = f

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, "INSERT INTO KillSwitch (feat_id, min_version, max_version, active) VALUES (?, ?, ?, ?)", k.Feature.Id, k.MinVersion, k.MaxVersion, k.Active)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	k.Id = int32(id)

	for _, b := range k.GetBrowsers() {
		if _, err := tx.ExecContext(ctx, "INSERT INTO KillSwitch2Browser (kswitch_id, browser) VALUES (?, ?)", k.Id, b); err != nil {
			tx.Rollback()
			return err
		}
	}

	logEntry := &pb.KillSwitchAuditLogEntry{
		Description: &pb.KillSwitchAuditLogEntry_KillSwitchEnabled_{
			&pb.KillSwitchAuditLogEntry_KillSwitchEnabled{
				KillSwitch: k,
			},
		},
	}
	if err := AddKillSwitchAuditLogEntry(tx, ctx, logEntry); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DisableKillSwitch(db *sql.DB, ctx context.Context, id int32) error {
	oldKillSwitch, err := GetKillSwitchById(db, ctx, id)
	if err != nil {
		return err
	}
	if oldKillSwitch == nil {
		return fmt.Errorf("Such kill switch doesn't exist")
	}
	if oldKillSwitch.GetActive() != true {
		return fmt.Errorf("The kill switch is not active")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	newKillSwitch := *oldKillSwitch
	newKillSwitch.Active = false

	if _, err := tx.ExecContext(ctx, "UPDATE KillSwitch SET active = ? WHERE kswitch_id = ?", newKillSwitch.Active, id); err != nil {
		tx.Rollback()
		return err
	}

	logEntry := &pb.KillSwitchAuditLogEntry{
		Description: &pb.KillSwitchAuditLogEntry_KillSwitchDisabled_{
			&pb.KillSwitchAuditLogEntry_KillSwitchDisabled{
				Transformation: &pb.KillSwitchTransformation{
					Old: oldKillSwitch,
					New: &newKillSwitch,
				},
			},
		},
	}
	if err := AddKillSwitchAuditLogEntry(tx, ctx, logEntry); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
