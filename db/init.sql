CREATE DATABASE IF NOT EXISTS `github-network`;
USE `github-network`;

-- DANGER!!!
-- DROP TABLE IF EXISTS pull_requests, issues, users;

CREATE TABLE IF NOT EXISTS pull_requests (
    id                            INT          PRIMARY KEY,
    owner                         VARCHAR(100) NOT NULL,
    repo                          VARCHAR(100) NOT NULL,
    number                        INT          NOT NULL,
    state                         VARCHAR(20)  NOT NULL,
    title                         TEXT,
    user_login                    VARCHAR(100) NOT NULL,
    body                          TEXT,
    active_lock_reason            TEXT,
    created_at                    VARCHAR(100) NOT NULL,
    closed_at                     VARCHAR(100),
    merged_at                     VARCHAR(100),
    assignees_login_csv           TEXT,
    requested_reviewers_login_csv TEXT,
    head_ref                      VARCHAR(100),
    base_ref                      VARCHAR(100),
    author_association            VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS issues (
    id                  INT PRIMARY KEY,
    owner               VARCHAR(100) NOT NULL,
    repo                VARCHAR(100) NOT NULL,
    number              INT          NOT NULL,
    state               VARCHAR(20)  NOT NULL,
    title               TEXT,
    body                TEXT,
    user_login          VARCHAR(100) NOT NULL,
    assignees_login_csv TEXT,
    active_lock_reason  TEXT,
    comments            INT,
    pull_request_url    VARCHAR(200),
    closed_at           VARCHAR(100),
    created_at          VARCHAR(100) NOT NULL,
    closed_by           VARCHAR(100),
    author_association  VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS users (
    id           INT          PRIMARY KEY,
    login        VARCHAR(100) NOT NULL,
    type         VARCHAR(20)  NOT NULL,
    company      VARCHAR(100) NOT NULL,
    location     VARCHAR(100) NOT NULL,
    public_repos INT          NOT NULL,
    created_at   VARCHAR(100) NOT NULL
);
