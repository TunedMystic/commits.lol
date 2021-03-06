CREATE TABLE IF NOT EXISTS config_badword (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    text VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS config_groupterm (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    text VARCHAR(50) NOT NULL,
    groupname VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS config_searchterm (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    text VARCHAR(50) NOT NULL,
    rank INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS git_user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source INTEGER NOT NULL,
    username VARCHAR(50) NOT NULL,
    url VARCHAR(100) UNIQUE NOT NULL,
    avatar_url VARCHAR(200) NOT NULL
);

CREATE TABLE IF NOT EXISTS git_repo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source INTEGER NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(100) NOT NULL,
    url VARCHAR(200) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS git_commit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    repo_id INTEGER NOT NULL,
    message VARCHAR(500) NOT NULL,
    message_censored VARCHAR(500) NOT NULL DEFAULT '',
    sha VARCHAR(40) NOT NULL,
    url VARCHAR(200) UNIQUE NOT NULL,
    date DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    valid BOOL NOT NULL DEFAULT TRUE,
    groupname VARCHAR(50) NOT NULL,
    color_bg VARCHAR(10) NOT NULL,
    color_fg VARCHAR(10) NOT NULL,
    FOREIGN KEY(repo_id) REFERENCES git_repo(id),
    FOREIGN KEY(author_id) REFERENCES git_user(id)
);
