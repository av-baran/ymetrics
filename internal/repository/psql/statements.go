package psql

const (
	createTableStatement = `
CREATE TABLE IF NOT EXISTS metrics (
	"id" VARCHAR(256) PRIMARY KEY,
	"value"  DOUBLE PRECISION,
	"type" TEXT,
	"delta" BIGINT);`

	setStatement = `
INSERT INTO metrics (id, type, value, delta) VALUES %s
	ON CONFLICT (id) DO UPDATE SET
		value = EXCLUDED.value,
		type  = EXCLUDED.type,
		delta = CASE WHEN metrics.delta IS NOT NULL OR EXCLUDED.delta IS NOT NULL
			THEN coalesce(metrics.delta, 0) + coalesce(EXCLUDED.delta, 0)
			ELSE NULL END;`

	getStatement    = `SELECT id, type, value, delta FROM metrics WHERE id=$1 AND type=$2 ORDER BY id;`
	getAllStatement = `SELECT id, type, value, delta FROM metrics ORDER BY id;`
)
