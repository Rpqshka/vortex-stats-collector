CREATE TABLE IF NOT EXISTS order_book
(
    id UUID DEFAULT generateUUIDv4(),
    exchange String,
    pair String,
    asks Nested(
           price Float64,
           base_qty Float64
    ),
    bids Nested(
           price Float64,
           base_qty Float64
    )
) ENGINE = MergeTree()
ORDER BY id;
