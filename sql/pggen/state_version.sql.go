// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertStateVersionSQL = `INSERT INTO state_versions (
    state_version_id,
    created_at,
    serial,
    state,
    workspace_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);`

type InsertStateVersionParams struct {
	ID          pgtype.Text
	CreatedAt   pgtype.Timestamptz
	Serial      int
	State       []byte
	WorkspaceID pgtype.Text
}

// InsertStateVersion implements Querier.InsertStateVersion.
func (q *DBQuerier) InsertStateVersion(ctx context.Context, params InsertStateVersionParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertStateVersion")
	cmdTag, err := q.conn.Exec(ctx, insertStateVersionSQL, params.ID, params.CreatedAt, params.Serial, params.State, params.WorkspaceID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertStateVersion: %w", err)
	}
	return cmdTag, err
}

// InsertStateVersionBatch implements Querier.InsertStateVersionBatch.
func (q *DBQuerier) InsertStateVersionBatch(batch genericBatch, params InsertStateVersionParams) {
	batch.Queue(insertStateVersionSQL, params.ID, params.CreatedAt, params.Serial, params.State, params.WorkspaceID)
}

// InsertStateVersionScan implements Querier.InsertStateVersionScan.
func (q *DBQuerier) InsertStateVersionScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertStateVersionBatch: %w", err)
	}
	return cmdTag, err
}

const findStateVersionsByWorkspaceNameSQL = `SELECT
    state_versions.*,
    array_remove(array_agg(state_version_outputs), NULL) AS state_version_outputs
FROM state_versions
JOIN workspaces USING (workspace_id)
LEFT JOIN state_version_outputs USING (state_version_id)
WHERE workspaces.name               = $1
AND   workspaces.organization_name  = $2
GROUP BY state_versions.state_version_id
ORDER BY created_at DESC
LIMIT $3
OFFSET $4
;`

type FindStateVersionsByWorkspaceNameParams struct {
	WorkspaceName    pgtype.Text
	OrganizationName pgtype.Text
	Limit            int
	Offset           int
}

type FindStateVersionsByWorkspaceNameRow struct {
	StateVersionID      pgtype.Text           `json:"state_version_id"`
	CreatedAt           pgtype.Timestamptz    `json:"created_at"`
	Serial              int                   `json:"serial"`
	State               []byte                `json:"state"`
	WorkspaceID         pgtype.Text           `json:"workspace_id"`
	StateVersionOutputs []StateVersionOutputs `json:"state_version_outputs"`
}

// FindStateVersionsByWorkspaceName implements Querier.FindStateVersionsByWorkspaceName.
func (q *DBQuerier) FindStateVersionsByWorkspaceName(ctx context.Context, params FindStateVersionsByWorkspaceNameParams) ([]FindStateVersionsByWorkspaceNameRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindStateVersionsByWorkspaceName")
	rows, err := q.conn.Query(ctx, findStateVersionsByWorkspaceNameSQL, params.WorkspaceName, params.OrganizationName, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindStateVersionsByWorkspaceName: %w", err)
	}
	defer rows.Close()
	items := []FindStateVersionsByWorkspaceNameRow{}
	stateVersionOutputsArray := q.types.newStateVersionOutputsArray()
	for rows.Next() {
		var item FindStateVersionsByWorkspaceNameRow
		if err := rows.Scan(&item.StateVersionID, &item.CreatedAt, &item.Serial, &item.State, &item.WorkspaceID, stateVersionOutputsArray); err != nil {
			return nil, fmt.Errorf("scan FindStateVersionsByWorkspaceName row: %w", err)
		}
		if err := stateVersionOutputsArray.AssignTo(&item.StateVersionOutputs); err != nil {
			return nil, fmt.Errorf("assign FindStateVersionsByWorkspaceName row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindStateVersionsByWorkspaceName rows: %w", err)
	}
	return items, err
}

// FindStateVersionsByWorkspaceNameBatch implements Querier.FindStateVersionsByWorkspaceNameBatch.
func (q *DBQuerier) FindStateVersionsByWorkspaceNameBatch(batch genericBatch, params FindStateVersionsByWorkspaceNameParams) {
	batch.Queue(findStateVersionsByWorkspaceNameSQL, params.WorkspaceName, params.OrganizationName, params.Limit, params.Offset)
}

// FindStateVersionsByWorkspaceNameScan implements Querier.FindStateVersionsByWorkspaceNameScan.
func (q *DBQuerier) FindStateVersionsByWorkspaceNameScan(results pgx.BatchResults) ([]FindStateVersionsByWorkspaceNameRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindStateVersionsByWorkspaceNameBatch: %w", err)
	}
	defer rows.Close()
	items := []FindStateVersionsByWorkspaceNameRow{}
	stateVersionOutputsArray := q.types.newStateVersionOutputsArray()
	for rows.Next() {
		var item FindStateVersionsByWorkspaceNameRow
		if err := rows.Scan(&item.StateVersionID, &item.CreatedAt, &item.Serial, &item.State, &item.WorkspaceID, stateVersionOutputsArray); err != nil {
			return nil, fmt.Errorf("scan FindStateVersionsByWorkspaceNameBatch row: %w", err)
		}
		if err := stateVersionOutputsArray.AssignTo(&item.StateVersionOutputs); err != nil {
			return nil, fmt.Errorf("assign FindStateVersionsByWorkspaceName row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindStateVersionsByWorkspaceNameBatch rows: %w", err)
	}
	return items, err
}

const countStateVersionsByWorkspaceNameSQL = `SELECT count(*)
FROM state_versions
JOIN workspaces USING (workspace_id)
WHERE workspaces.name                 = $1
AND   workspaces.organization_name    = $2
;`

// CountStateVersionsByWorkspaceName implements Querier.CountStateVersionsByWorkspaceName.
func (q *DBQuerier) CountStateVersionsByWorkspaceName(ctx context.Context, workspaceName pgtype.Text, organizationName pgtype.Text) (*int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountStateVersionsByWorkspaceName")
	row := q.conn.QueryRow(ctx, countStateVersionsByWorkspaceNameSQL, workspaceName, organizationName)
	var item int
	if err := row.Scan(&item); err != nil {
		return &item, fmt.Errorf("query CountStateVersionsByWorkspaceName: %w", err)
	}
	return &item, nil
}

// CountStateVersionsByWorkspaceNameBatch implements Querier.CountStateVersionsByWorkspaceNameBatch.
func (q *DBQuerier) CountStateVersionsByWorkspaceNameBatch(batch genericBatch, workspaceName pgtype.Text, organizationName pgtype.Text) {
	batch.Queue(countStateVersionsByWorkspaceNameSQL, workspaceName, organizationName)
}

// CountStateVersionsByWorkspaceNameScan implements Querier.CountStateVersionsByWorkspaceNameScan.
func (q *DBQuerier) CountStateVersionsByWorkspaceNameScan(results pgx.BatchResults) (*int, error) {
	row := results.QueryRow()
	var item int
	if err := row.Scan(&item); err != nil {
		return &item, fmt.Errorf("scan CountStateVersionsByWorkspaceNameBatch row: %w", err)
	}
	return &item, nil
}

const findStateVersionByIDSQL = `SELECT
    state_versions.*,
    array_remove(array_agg(state_version_outputs), NULL) AS state_version_outputs
FROM state_versions
LEFT JOIN state_version_outputs USING (state_version_id)
WHERE state_versions.state_version_id = $1
GROUP BY state_versions.state_version_id
;`

type FindStateVersionByIDRow struct {
	StateVersionID      pgtype.Text           `json:"state_version_id"`
	CreatedAt           pgtype.Timestamptz    `json:"created_at"`
	Serial              int                   `json:"serial"`
	State               []byte                `json:"state"`
	WorkspaceID         pgtype.Text           `json:"workspace_id"`
	StateVersionOutputs []StateVersionOutputs `json:"state_version_outputs"`
}

// FindStateVersionByID implements Querier.FindStateVersionByID.
func (q *DBQuerier) FindStateVersionByID(ctx context.Context, id pgtype.Text) (FindStateVersionByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindStateVersionByID")
	row := q.conn.QueryRow(ctx, findStateVersionByIDSQL, id)
	var item FindStateVersionByIDRow
	stateVersionOutputsArray := q.types.newStateVersionOutputsArray()
	if err := row.Scan(&item.StateVersionID, &item.CreatedAt, &item.Serial, &item.State, &item.WorkspaceID, stateVersionOutputsArray); err != nil {
		return item, fmt.Errorf("query FindStateVersionByID: %w", err)
	}
	if err := stateVersionOutputsArray.AssignTo(&item.StateVersionOutputs); err != nil {
		return item, fmt.Errorf("assign FindStateVersionByID row: %w", err)
	}
	return item, nil
}

// FindStateVersionByIDBatch implements Querier.FindStateVersionByIDBatch.
func (q *DBQuerier) FindStateVersionByIDBatch(batch genericBatch, id pgtype.Text) {
	batch.Queue(findStateVersionByIDSQL, id)
}

// FindStateVersionByIDScan implements Querier.FindStateVersionByIDScan.
func (q *DBQuerier) FindStateVersionByIDScan(results pgx.BatchResults) (FindStateVersionByIDRow, error) {
	row := results.QueryRow()
	var item FindStateVersionByIDRow
	stateVersionOutputsArray := q.types.newStateVersionOutputsArray()
	if err := row.Scan(&item.StateVersionID, &item.CreatedAt, &item.Serial, &item.State, &item.WorkspaceID, stateVersionOutputsArray); err != nil {
		return item, fmt.Errorf("scan FindStateVersionByIDBatch row: %w", err)
	}
	if err := stateVersionOutputsArray.AssignTo(&item.StateVersionOutputs); err != nil {
		return item, fmt.Errorf("assign FindStateVersionByID row: %w", err)
	}
	return item, nil
}

const findCurrentStateVersionByWorkspaceIDSQL = `SELECT
    sv.*,
    array_remove(array_agg(svo), NULL) AS state_version_outputs
FROM state_versions sv
LEFT JOIN state_version_outputs svo USING (state_version_id)
JOIN workspaces w ON w.current_state_version_id = sv.state_version_id
WHERE w.workspace_id = $1
GROUP BY sv.state_version_id
;`

type FindCurrentStateVersionByWorkspaceIDRow struct {
	StateVersionID      pgtype.Text           `json:"state_version_id"`
	CreatedAt           pgtype.Timestamptz    `json:"created_at"`
	Serial              int                   `json:"serial"`
	State               []byte                `json:"state"`
	WorkspaceID         pgtype.Text           `json:"workspace_id"`
	StateVersionOutputs []StateVersionOutputs `json:"state_version_outputs"`
}

// FindCurrentStateVersionByWorkspaceID implements Querier.FindCurrentStateVersionByWorkspaceID.
func (q *DBQuerier) FindCurrentStateVersionByWorkspaceID(ctx context.Context, workspaceID pgtype.Text) (FindCurrentStateVersionByWorkspaceIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindCurrentStateVersionByWorkspaceID")
	row := q.conn.QueryRow(ctx, findCurrentStateVersionByWorkspaceIDSQL, workspaceID)
	var item FindCurrentStateVersionByWorkspaceIDRow
	stateVersionOutputsArray := q.types.newStateVersionOutputsArray()
	if err := row.Scan(&item.StateVersionID, &item.CreatedAt, &item.Serial, &item.State, &item.WorkspaceID, stateVersionOutputsArray); err != nil {
		return item, fmt.Errorf("query FindCurrentStateVersionByWorkspaceID: %w", err)
	}
	if err := stateVersionOutputsArray.AssignTo(&item.StateVersionOutputs); err != nil {
		return item, fmt.Errorf("assign FindCurrentStateVersionByWorkspaceID row: %w", err)
	}
	return item, nil
}

// FindCurrentStateVersionByWorkspaceIDBatch implements Querier.FindCurrentStateVersionByWorkspaceIDBatch.
func (q *DBQuerier) FindCurrentStateVersionByWorkspaceIDBatch(batch genericBatch, workspaceID pgtype.Text) {
	batch.Queue(findCurrentStateVersionByWorkspaceIDSQL, workspaceID)
}

// FindCurrentStateVersionByWorkspaceIDScan implements Querier.FindCurrentStateVersionByWorkspaceIDScan.
func (q *DBQuerier) FindCurrentStateVersionByWorkspaceIDScan(results pgx.BatchResults) (FindCurrentStateVersionByWorkspaceIDRow, error) {
	row := results.QueryRow()
	var item FindCurrentStateVersionByWorkspaceIDRow
	stateVersionOutputsArray := q.types.newStateVersionOutputsArray()
	if err := row.Scan(&item.StateVersionID, &item.CreatedAt, &item.Serial, &item.State, &item.WorkspaceID, stateVersionOutputsArray); err != nil {
		return item, fmt.Errorf("scan FindCurrentStateVersionByWorkspaceIDBatch row: %w", err)
	}
	if err := stateVersionOutputsArray.AssignTo(&item.StateVersionOutputs); err != nil {
		return item, fmt.Errorf("assign FindCurrentStateVersionByWorkspaceID row: %w", err)
	}
	return item, nil
}

const findStateVersionStateByIDSQL = `SELECT state
FROM state_versions
WHERE state_version_id = $1
;`

// FindStateVersionStateByID implements Querier.FindStateVersionStateByID.
func (q *DBQuerier) FindStateVersionStateByID(ctx context.Context, id pgtype.Text) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindStateVersionStateByID")
	row := q.conn.QueryRow(ctx, findStateVersionStateByIDSQL, id)
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindStateVersionStateByID: %w", err)
	}
	return item, nil
}

// FindStateVersionStateByIDBatch implements Querier.FindStateVersionStateByIDBatch.
func (q *DBQuerier) FindStateVersionStateByIDBatch(batch genericBatch, id pgtype.Text) {
	batch.Queue(findStateVersionStateByIDSQL, id)
}

// FindStateVersionStateByIDScan implements Querier.FindStateVersionStateByIDScan.
func (q *DBQuerier) FindStateVersionStateByIDScan(results pgx.BatchResults) ([]byte, error) {
	row := results.QueryRow()
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindStateVersionStateByIDBatch row: %w", err)
	}
	return item, nil
}

const deleteStateVersionByIDSQL = `DELETE
FROM state_versions
WHERE state_version_id = $1
RETURNING state_version_id
;`

// DeleteStateVersionByID implements Querier.DeleteStateVersionByID.
func (q *DBQuerier) DeleteStateVersionByID(ctx context.Context, stateVersionID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteStateVersionByID")
	row := q.conn.QueryRow(ctx, deleteStateVersionByIDSQL, stateVersionID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteStateVersionByID: %w", err)
	}
	return item, nil
}

// DeleteStateVersionByIDBatch implements Querier.DeleteStateVersionByIDBatch.
func (q *DBQuerier) DeleteStateVersionByIDBatch(batch genericBatch, stateVersionID pgtype.Text) {
	batch.Queue(deleteStateVersionByIDSQL, stateVersionID)
}

// DeleteStateVersionByIDScan implements Querier.DeleteStateVersionByIDScan.
func (q *DBQuerier) DeleteStateVersionByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteStateVersionByIDBatch row: %w", err)
	}
	return item, nil
}
