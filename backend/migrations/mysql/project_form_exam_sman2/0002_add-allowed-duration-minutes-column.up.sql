ALTER TABLE exams ADD allowed_duration_minutes INT DEFAULT 120;
ALTER TABLE participants ADD allowed_duration_minutes INT DEFAULT 120;