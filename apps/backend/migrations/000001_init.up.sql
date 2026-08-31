-- Migrasi awal QuestDay.
-- Isi DDL sesuai entitas di tiap module. Kerangka tabel & alasannya di bawah.
-- Buat file baru untuk perubahan berikutnya: `make migrate-create name=...`.

-- users
--   id            (pk, uuid/text)
--   email         (unique, not null)          -> deteksi ErrEmailTaken
--   password_hash (not null)
--   display_name  (not null)
--   timezone      (not null, default 'Asia/Jakarta')  -> batas hari lokal
--   created_at    (timestamptz, default now())
-- TODO: CREATE TABLE users (...); CREATE UNIQUE INDEX ... ON users(email);

-- quests (definisi)
--   id, user_id (fk users), title, note,
--   category, difficulty, recurrence, active (bool default true),
--   created_at
-- TODO: CREATE TABLE quests (...); INDEX (user_id);

-- quest_logs (instance harian)
--   id, quest_id (fk quests), user_id (fk users),
--   date (DATE — tanggal lokal user), status, points_awarded, completed_at
--   UNIQUE (quest_id, date)  -> cegah dobel selesai di hari sama
-- TODO: CREATE TABLE quest_logs (...); INDEX (user_id, date);

-- wallets (scoring)
--   user_id (pk, fk users), total_points, xp, level
-- TODO: CREATE TABLE wallets (...);

-- streaks (scoring)
--   user_id (pk, fk users), current, longest, last_active (DATE)
-- TODO: CREATE TABLE streaks (...);

-- point_transactions (ledger scoring, untuk audit & rollback)
--   id, user_id, quest_id, points (bisa negatif), date, created_at
-- TODO: CREATE TABLE point_transactions (...); INDEX (user_id, created_at);

-- (v2) badges, unlocks  -> tunda sampai module achievement digarap.
