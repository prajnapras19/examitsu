CREATE TABLE participant_sessions(
    id BIGINT NOT NULL AUTO_INCREMENT,

    serial VARCHAR(255) NOT NULL,
    participant_id BIGINT NOT NULL,
    is_authorized TINYINT(1) NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    not_archived BOOLEAN GENERATED ALWAYS AS (IF(deleted_at IS NULL, 1, NULL)) VIRTUAL,

    CONSTRAINT PK_id PRIMARY KEY (id),
    CONSTRAINT FOREIGN KEY (participant_id) REFERENCES participants(id),
    CONSTRAINT UNIQUE (serial, not_archived),
    INDEX (participant_id, is_authorized, not_archived)
);