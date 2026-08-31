-- Migrasi awal QuestDay (MVP: user, quest, scoring).
-- PostgreSQL 16. Format golang-migrate (migrate membungkus tiap file dalam transaksi).
-- id uuid TANPA default: aplikasi men-generate UUIDv7. Tanpa pgcrypto.
-- Tabel v2 (badges, unlocks) sengaja tidak dibuat.

-- users
CREATE TABLE users (
    id            uuid PRIMARY KEY,
    email         text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    display_name  text NOT NULL,
    timezone      text NOT NULL DEFAULT 'Asia/Jakarta',
    created_at    timestamptz NOT NULL DEFAULT now()
);

-- quests (definisi)
CREATE TABLE quests (
    id         uuid PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      text NOT NULL,
    note       text NOT NULL DEFAULT '',
    category   text NOT NULL,
    difficulty text NOT NULL,
    recurrence text NOT NULL DEFAULT 'daily',
    active     boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON quests (user_id);

-- quest_logs (instance harian; date = tanggal lokal user, ADR-006)
CREATE TABLE quest_logs (
    id             uuid PRIMARY KEY,
    quest_id       uuid NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date           date NOT NULL,
    status         text NOT NULL,
    points_awarded integer NOT NULL DEFAULT 0,
    completed_at   timestamptz NOT NULL DEFAULT now(),
    UNIQUE (quest_id, date)
);

CREATE INDEX ON quest_logs (user_id, date);

-- wallets (scoring)
CREATE TABLE wallets (
    user_id      uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_points integer NOT NULL DEFAULT 0,
    xp           integer NOT NULL DEFAULT 0,
    level        integer NOT NULL DEFAULT 1
);

-- streaks (scoring)
CREATE TABLE streaks (
    user_id     uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current     integer NOT NULL DEFAULT 0,
    longest     integer NOT NULL DEFAULT 0,
    last_active date
);

-- point_transactions (ledger scoring; points bisa negatif)
CREATE TABLE point_transactions (
    id         uuid PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_id   uuid REFERENCES quests(id) ON DELETE SET NULL,
    points     integer NOT NULL,
    date       date NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON point_transactions (user_id, created_at);
