package db

import (
	"context"
	"database/sql"
	"fmt"

	pb "gomodules.avm99963.com/twpt-server/api_proto"
)

func GetAuthorizedUserById(db *sql.DB, ctx context.Context, id int32) (*pb.KillSwitchAuthorizedUser, error) {
	query := db.QueryRowContext(ctx, "SELECT user_id, google_uid, email, access_level FROM KillSwitchAuthorizedUser WHERE user_id = ?", id)
	var u pb.KillSwitchAuthorizedUser
	if err := query.Scan(&u.Id, &u.GoogleUid, &u.Email, &u.AccessLevel); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetAuthorizedUserById: %v.", err)
	}
	return &u, nil
}

func AddAuthorizedUser(db *sql.DB, ctx context.Context, u *pb.KillSwitchAuthorizedUser) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := db.ExecContext(ctx, "INSERT INTO KillSwitchAuthorizedUser (google_uid, email, access_level) VALUES (?, ?, ?)", u.GoogleUid, u.Email, u.AccessLevel)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	u.Id = int32(id)

	logEntry := &pb.KillSwitchAuditLogEntry{
		Description: &pb.KillSwitchAuditLogEntry_AuthorizedUserAdded_{
			&pb.KillSwitchAuditLogEntry_AuthorizedUserAdded{
				User: u,
			},
		},
	}
	if err := AddKillSwitchAuditLogEntry(tx, ctx, logEntry); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdateAuthorizedUser(db *sql.DB, ctx context.Context, id int32, newUser *pb.KillSwitchAuthorizedUser) error {
	oldUser, err := GetAuthorizedUserById(db, ctx, id)
	if err != nil {
		return err
	}
	if oldUser == nil {
		return fmt.Errorf("Such user doesn't exist")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "UPDATE KillSwitchAuthorizedUser SET google_uid = ?, email = ?, access_level = ? WHERE user_id = ?", newUser.GoogleUid, newUser.Email, newUser.AccessLevel, id); err != nil {
		tx.Rollback()
		return err
	}

	newUser.Id = id

	logEntry := &pb.KillSwitchAuditLogEntry{
		Description: &pb.KillSwitchAuditLogEntry_AuthorizedUserUpdated_{
			&pb.KillSwitchAuditLogEntry_AuthorizedUserUpdated{
				Transformation: &pb.AuthorizedUserTransformation{
					Old: oldUser,
					New: newUser,
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

func DeleteAuthorizedUser(db *sql.DB, ctx context.Context, id int32) error {
	u, err := GetAuthorizedUserById(db, ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return fmt.Errorf("Such user doesn't exist")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM KillSwitchAuthorizedUser WHERE user_id = ?", id); err != nil {
		tx.Rollback()
		return err
	}

	logEntry := &pb.KillSwitchAuditLogEntry{
		Description: &pb.KillSwitchAuditLogEntry_AuthorizedUserDeleted_{
			&pb.KillSwitchAuditLogEntry_AuthorizedUserDeleted{
				OldUser: u,
			},
		},
	}
	if err := AddKillSwitchAuditLogEntry(tx, ctx, logEntry); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func ListAuthorizedUsers(db *sql.DB, ctx context.Context) ([]*pb.KillSwitchAuthorizedUser, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.QueryContext(ctx, "SELECT user_id, google_uid, email, access_level FROM KillSwitchAuthorizedUser")
	if err != nil {
		return nil, fmt.Errorf("ListAuthorizedUsers: %v", err)
	}
	defer rows.Close()

	var users []*pb.KillSwitchAuthorizedUser
	for rows.Next() {
		var u pb.KillSwitchAuthorizedUser
		if err := rows.Scan(&u.Id, &u.GoogleUid, &u.Email, &u.AccessLevel); err != nil {
			return nil, fmt.Errorf("ListAuthorizedUsers: %v", err)
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListAuthorizedUsers: %v", err)
	}
	return users, nil
}
