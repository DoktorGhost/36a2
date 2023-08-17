CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    pub_date TEXT,
    source TEXT
);