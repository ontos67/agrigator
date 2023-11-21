DROP TABLE IF EXISTS articles;
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    a_url VARCHAR(255) NOT NULL,
    publisher TEXT NOT NULL,
    author TEXT NOT NULL
);
