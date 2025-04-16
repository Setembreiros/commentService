CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    postId VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índice para búsquedas por username
CREATE INDEX idx_comments_username ON comments(username);

-- Índice para búsquedas por postId
CREATE INDEX idx_comments_postid ON comments(postId);

-- Índice compuesto para búsquedas que combinan username y postId
CREATE INDEX idx_comments_username_postid ON comments(username, postId);