-- 005_user_lockout.sql — Add failed login attempt tracking and lockout columns to usuarios

ALTER TABLE usuarios ADD COLUMN intentos_fallidos INTEGER NOT NULL DEFAULT 0;
ALTER TABLE usuarios ADD COLUMN bloqueado_hasta TEXT;
