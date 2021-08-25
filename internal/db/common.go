package db

import (
	"context"
	"database/sql"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gomodules.avm99963.com/twpt-server/api_proto"
)

func AddKillSwitchAuditLogEntry(tx *sql.Tx, ctx context.Context, logEntry *pb.KillSwitchAuditLogEntry) error {
	logEntry.Timestamp = timestamppb.Now()

	logEntryBytes, err := proto.Marshal(logEntry)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO KillSwitchAuditLog (data) VALUES (?)", logEntryBytes); err != nil {
		return err
	}

	return nil
}
