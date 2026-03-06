CREATE OR REPLACE FUNCTION insert_staged_originator_envelope_batch_v2(
    p_topics          bytea[],
    p_payer_envelopes bytea[]
)
    RETURNS TABLE (
                      id              bigint,
                      originator_time timestamp,
                      topic           bytea,
                      payer_envelope  bytea
                  )
    LANGUAGE plpgsql
AS $$
BEGIN
    IF array_length(p_topics, 1) IS DISTINCT FROM array_length(p_payer_envelopes, 1) THEN
        RAISE EXCEPTION
            'p_topics and p_payer_envelopes must have the same length';
    END IF;

    -- Ensures that generated sequence IDs follow batch input order.
    -- Lock is held until the end of the enclosing transaction.
    PERFORM pg_advisory_xact_lock(hashtext('staged_originator_envelopes_sequence'));

    RETURN QUERY
        WITH input AS (
            SELECT
                i,
                p_topics[i]          AS topic,
                p_payer_envelopes[i] AS payer_envelope
            FROM generate_subscripts(p_topics, 1) AS g(i)
        ),
             inserted AS (
                 INSERT INTO staged_originator_envelopes (
                                                          topic,
                                                          payer_envelope
                     )
                     SELECT
                         input.topic,
                         input.payer_envelope
                     FROM input
                     ORDER BY input.i
                     ON CONFLICT DO NOTHING
                     RETURNING
                         staged_originator_envelopes.id,
                         staged_originator_envelopes.originator_time,
                         staged_originator_envelopes.topic,
                         staged_originator_envelopes.payer_envelope
             )
        SELECT
            inserted.id,
            inserted.originator_time,
            inserted.topic,
            inserted.payer_envelope
        FROM inserted;
END;
$$;