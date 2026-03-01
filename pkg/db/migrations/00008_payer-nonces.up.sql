CREATE TABLE nonce_table (
    nonce BIGINT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION fill_nonce_gap(pending_nonce BIGINT, num_elements INT)
    RETURNS INT AS $$
DECLARE
    inserted_rows INT;
BEGIN
    WITH nonces AS (
        -- Generate the required number of nonces
        SELECT generate_series(pending_nonce, pending_nonce + num_elements - 1) AS nonce
    )
    INSERT INTO nonce_table (nonce)
    SELECT nonce
    FROM nonces n
    WHERE NOT EXISTS (SELECT 1 FROM nonce_table nt WHERE nt.nonce = n.nonce) -- Skip existing ones
    ON CONFLICT DO NOTHING; -- Ensure no duplicates

    -- Capture the number of inserted rows
    GET DIAGNOSTICS inserted_rows = ROW_COUNT;

    RETURN inserted_rows;
END;
$$ LANGUAGE plpgsql;
