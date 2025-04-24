CREATE TABLE commentservice.comments (
    id BIGSERIAL PRIMARY KEY,
    postId VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índice para búsquedas por username
CREATE INDEX idx_comments_username ON commentservice.comments(username);

-- Índice para búsquedas por postId
CREATE INDEX idx_comments_postid ON commentservice.comments(postId);

-- Índice compuesto para búsquedas que combinan username y postId
CREATE INDEX idx_comments_username_postid ON commentservice.comments(username, postId);