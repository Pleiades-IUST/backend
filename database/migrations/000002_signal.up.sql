CREATE TABLE IF NOT EXISTS signal (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    drive_id BIGINT REFERENCES drive(id),
    technology TEXT,
    strength INT,
    rsrp INT,
    rsrq INT
);