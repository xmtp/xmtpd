-- pgFormatter-ignore
CREATE FUNCTION select_gateway_envelopes(cursor_node_ids INT[], cursor_sequence_ids BIGINT[], topics BYTEA[], originator_node_ids INT[], row_limit INT)
	RETURNS SETOF gateway_envelopes
	AS $$
DECLARE
	num_topics INT := COALESCE(ARRAY_LENGTH(topics, 1), 0);
	num_originators INT := COALESCE(ARRAY_LENGTH(originator_node_ids, 1), 0);
BEGIN
	RETURN QUERY
    WITH cursors AS (
		SELECT
			UNNEST(cursor_node_ids) AS cursor_node_id,
			UNNEST(cursor_sequence_ids) AS cursor_sequence_id
    )
	SELECT
		gateway_envelopes.*
	FROM
		gateway_envelopes
	-- Assumption: There is only one cursor per node ID. Caller must verify this
	LEFT JOIN cursors ON gateway_envelopes.originator_node_id = cursors.cursor_node_id
    WHERE (num_topics = 0 OR topic = ANY (topics))
		AND (num_originators = 0 OR originator_node_id = ANY (originator_node_ids))
		AND originator_sequence_id > COALESCE(cursor_sequence_id, 0)
	ORDER BY
		-- Assumption: envelopes are inserted in sequence_id order per originator, therefore
		-- gateway_time preserves sequence_id order
		gateway_time,
		originator_node_id,
		originator_sequence_id ASC
	LIMIT NULLIF(row_limit, 0);
END;
$$
LANGUAGE plpgsql;

