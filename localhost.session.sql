EXPLAIN ANALYZE
SELECT h.id,
    a.id,
    a.name,
    a.instrument_type,
    h.total_quantity,
    h.average_price,
    pd.curr_price,
    pd.prev_price,
    h.total_invested
FROM holdings h
    INNER JOIN user_assets ua ON h.user_asset_id = ua.id
    INNER JOIN assets a ON ua.asset_id = a.id
    INNER JOIN price_details pd ON a.id = pd.asset_id
WHERE ua.user_id = $1
ORDER BY h.id
LIMIT $2 OFFSET $3