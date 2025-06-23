CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       surname VARCHAR(255) NOT NULL,
                       login VARCHAR(255) NOT NULL UNIQUE,
                       password VARCHAR(255) NOT NULL,
                       timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_login_password ON users(login, password);

CREATE TABLE folders (
                         id SERIAL PRIMARY KEY,
                         title VARCHAR(255) NOT NULL,
                         timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         user_id INTEGER NOT NULL,
                         FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_folders_user_id ON folders(user_id);

CREATE TABLE notes (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       content TEXT,
                       user_id INTEGER NOT NULL,
                       is_favorite BOOLEAN NOT NULL DEFAULT FALSE,
                       timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       tags TEXT[] DEFAULT '{}',
                       folder_id INTEGER,
                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                       FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
);

CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_notes_user_id_favorite ON notes(user_id, is_favorite) WHERE is_favorite = true;
CREATE INDEX idx_notes_tags ON notes USING gin(tags);