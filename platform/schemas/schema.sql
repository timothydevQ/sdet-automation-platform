-- Schema for ingested test results.
-- Lives in its own database/namespace; not the application postgres.

CREATE TABLE IF NOT EXISTS test_runs (
    id BIGSERIAL PRIMARY KEY,
    suite TEXT NOT NULL,
    branch TEXT,
    commit_sha TEXT,
    workflow_run TEXT,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ,
    total INT NOT NULL DEFAULT 0,
    passed INT NOT NULL DEFAULT 0,
    failed INT NOT NULL DEFAULT 0,
    skipped INT NOT NULL DEFAULT 0,
    duration_ms BIGINT,
    metadata JSONB
);

CREATE INDEX IF NOT EXISTS test_runs_suite_idx ON test_runs (suite);
CREATE INDEX IF NOT EXISTS test_runs_started_idx ON test_runs (started_at DESC);

CREATE TABLE IF NOT EXISTS test_cases (
    id BIGSERIAL PRIMARY KEY,
    run_id BIGINT NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    suite TEXT NOT NULL,
    classname TEXT,
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('passed', 'failed', 'skipped', 'error')),
    duration_ms BIGINT,
    retries INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS test_cases_name_idx ON test_cases (suite, classname, name);
CREATE INDEX IF NOT EXISTS test_cases_run_idx ON test_cases (run_id);

CREATE TABLE IF NOT EXISTS test_failures (
    id BIGSERIAL PRIMARY KEY,
    case_id BIGINT NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    message TEXT,
    stack TEXT,
    failure_type TEXT
);

CREATE TABLE IF NOT EXISTS test_artifacts (
    id BIGSERIAL PRIMARY KEY,
    case_id BIGINT NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    kind TEXT NOT NULL,
    path TEXT NOT NULL,
    size_bytes BIGINT
);

CREATE TABLE IF NOT EXISTS flake_scores (
    id BIGSERIAL PRIMARY KEY,
    suite TEXT NOT NULL,
    classname TEXT,
    name TEXT NOT NULL,
    score DOUBLE PRECISION NOT NULL,
    classification TEXT NOT NULL,
    failure_rate DOUBLE PRECISION,
    pass_after_retry_rate DOUBLE PRECISION,
    duration_variance DOUBLE PRECISION,
    recent_failure_streak INT,
    samples INT NOT NULL,
    computed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (suite, classname, name)
);

CREATE TABLE IF NOT EXISTS performance_results (
    id BIGSERIAL PRIMARY KEY,
    run_id BIGINT REFERENCES test_runs(id) ON DELETE CASCADE,
    scenario TEXT NOT NULL,
    metric TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    threshold DOUBLE PRECISION,
    passed BOOLEAN
);
