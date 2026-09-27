# Render Free public demo

This deployment is an educational, publicly accessible demo with known seed PINs
and synthetic business records. It is NOT a production system. Do not enter real
customers, sales, inventory, API keys or other confidential information.

- Select Free (no disk, no paid database). `/data/pos.db` is ephemeral; data may be
  lost when the container is replaced, redeployed or spun down. A fresh database
  recreates the embedded synthetic seed. Same-filesystem restart does not reseed.
- Use `Dockerfile.render` and the verified deployment branch/commit. `.dockerignore`
  allows only required source/assets; no `.env`, local database, Git history,
  presentations or private preview artifacts enter the build context.
- Generate `SESSION_SECRET` in Render. Set `APP_ENV=local`,
  `DATABASE_PATH=/data/pos.db`, `SESSION_COOKIE_SECURE=true`. The explicit HTTPS
  cookie policy works behind Render's TLS proxy without trusting forwarding
  headers. Local HTTP development defaults to false. Lambda stays HTTPS-only.
- Leave `OPENROUTER_API_KEY` empty. Chat shows AI unavailable and makes no external
  request, including no paid fallback. Enabling an AI provider requires separate
  authorization; no local credential is transferred by this deployment.
- Automatic redeploy is disabled. Publish only an explicitly approved source
  boundary. After deploying, verify `/health`, login, English clients and ES/EN
  switching. A sleeping Free service can take time to start.

Ledgerless legacy SQLite databases are deliberately unsupported. Do not upload,
modify or delete the local `data/pos.db` to make this disposable deployment work.
