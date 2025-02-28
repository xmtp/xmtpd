CREATE TABLE nonce_table (
    nonce BIGINT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION fill_nonce_gap(pending_nonce BIGINT, num_elements INT)
    RETURNS VOID AS $$
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
END;
$$ LANGUAGE plpgsql;
