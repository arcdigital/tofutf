package logs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/leg100/otf"
	"github.com/leg100/otf/sql"
	"github.com/leg100/otf/sql/pggen"
)

// pgdb is a logs database on postgres
type pgdb struct {
	otf.DB // provides access to generated SQL queries
}

// put persists a chunk of logs to the DB and returns the chunk updated with a
// unique identifier

// put persists data to the DB and returns a unique identifier for the chunk
func (db *pgdb) put(ctx context.Context, opts otf.PutChunkOptions) (string, error) {
	if len(opts.Data) == 0 {
		return "", fmt.Errorf("refusing to persist empty chunk")
	}
	id, err := db.InsertLogChunk(ctx, pggen.InsertLogChunkParams{
		RunID:  sql.String(opts.RunID),
		Phase:  sql.String(string(opts.Phase)),
		Chunk:  opts.Data,
		Offset: opts.Offset,
	})
	if err != nil {
		return "", sql.Error(err)
	}
	return strconv.Itoa(id), nil
}

// GetByID implements pubsub.Getter
func (db *pgdb) GetByID(ctx context.Context, chunkID string) (any, error) {
	id, err := strconv.Atoi(chunkID)
	if err != nil {
		return nil, err
	}
	chunk, err := db.FindLogChunkByID(ctx, id)
	if err != nil {
		return nil, sql.Error(err)
	}
	return otf.Chunk{
		ID:     chunkID,
		RunID:  chunk.RunID.String,
		Phase:  otf.PhaseType(chunk.Phase.String),
		Data:   chunk.Chunk,
		Offset: chunk.Offset,
	}, nil
}
