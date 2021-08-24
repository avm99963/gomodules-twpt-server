package db

import (
	"context"
	"database/sql"
	"fmt"

	pb "gomodules.avm99963.com/twpt-server/api_proto"
)

func GetFeatureByCodename(db *sql.DB, ctx context.Context, codename string) (*pb.Feature, error) {
	query := db.QueryRowContext(ctx, "SELECT feat_id, codename, feat_type FROM Feature WHERE codename = ?", codename)
	var f pb.Feature
	if err := query.Scan(&f.Id, &f.Codename, &f.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while querying feature by codename: %v.", err)
	}
	return &f, nil
}

func ListFeatures(db *sql.DB, ctx context.Context, withDeprecatedFeatures bool) ([]*pb.Feature, error) {
	var rows *sql.Rows
	var err error
	if withDeprecatedFeatures {
    rows, err = db.QueryContext(ctx, "SELECT feat_id, codename, feat_type FROM Feature")
	} else {
    rows, err = db.QueryContext(ctx, "SELECT feat_id, codename, feat_type FROM Feature WHERE feat_type <> ?", pb.Feature_TYPE_DEPRECATED)
	}
	if err != nil {
		return nil, fmt.Errorf("ListFeatures: ", err)
	}
	defer rows.Close()

	var features []*pb.Feature
	for rows.Next() {
		var f pb.Feature
		if err := rows.Scan(&f.Id, &f.Codename, &f.Type); err != nil {
			return nil, fmt.Errorf("ListFeatures: ", err)
		}
		features = append(features, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListFeatures: ", err)
	}
	return features, nil
}

func AddFeature(db *sql.DB, ctx context.Context, f *pb.Feature) error {
	_, err := db.ExecContext(ctx, "INSERT INTO Feature (codename, feat_type) VALUES (?, ?)", f.Codename, f.Type)
	return err
}

func UpdateFeature(db *sql.DB, ctx context.Context, id int32, f *pb.Feature) error {
	_, err := db.ExecContext(ctx, "UPDATE Feature SET feat_type = ? WHERE id = ?", f.Type, id)
	return err
}
