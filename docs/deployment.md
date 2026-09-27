# Containers and deployment

The image builds the frontend and three Go binaries, then runs the server as a non-root user in distroless. Node, the compiler, Storybook, and skills are excluded from the runtime image.

```sh
docker build -t my-app .
```

Set `DATABASE_URL` and `PUBLIC_URL=https://app.example.com`. The image defaults to `APP_ENV=production` and `HTTP_ADDR=0.0.0.0:3000`. Supply secrets through your platform. `compose.yml` is for local development; its password and disabled TLS are unsuitable for a public database.

These commands assume the environment variables have already been supplied securely. Run the server and worker as separate services:

```sh
docker run --rm --env DATABASE_URL --env PUBLIC_URL --entrypoint /app/manage my-app migrate
docker run --rm -i --env DATABASE_URL --env PUBLIC_URL --entrypoint /app/manage my-app create-user --email you@example.com --admin
docker run --rm -p 127.0.0.1:3000:3000 --env DATABASE_URL --env PUBLIC_URL my-app
docker run --rm --env DATABASE_URL --env PUBLIC_URL --entrypoint /app/worker my-app
```

The user-creation command reads the password from standard input. Run migrations as a separate step before starting a release. Design migrations for compatibility between concurrently running versions; a transaction does not make every DDL change safe for a rolling deployment.

The reverse proxy terminates TLS, preserves the public Host, restricts allowed hostnames, and limits request sizes and rates. Secure cookies require HTTPS in the browser. Keep the database, profiler, and internal interfaces off the public internet. The built-in sign-in limiter is local to each process and sees the connection address; configure an ingress limit that accounts for the proxy and instance count.

`/health/live` checks the process; `/health/ready` checks PostgreSQL connectivity. SIGTERM starts graceful shutdown with the `SHUTDOWN_TIMEOUT` limit, defaulting to ten seconds. The worker has an additional five seconds to cancel work after its drain timeout; allow a longer termination grace period on your platform. Model and SMTP timeouts are bounded separately. Back up PostgreSQL and test restoration.

The template does not provision infrastructure or publish the application to a host. Before deployment, verify DNS/TLS, secrets, database permissions, proxy configuration, backups, and release behavior. Production CSP blocks inline scripts; external integrations need an explicit policy change and browser validation.

## Server and worker configuration

The server and worker need consistent MAIL/AI settings as described in [email](background.md) and [AI](ai.md). Add the corresponding `--env` or secret bindings to the commands above; the base commands leave both features disabled.

The runtime pool budget applies per process. The worker also opens one LISTEN/NOTIFY connection with the same search path. With PgBouncer, LISTEN requires session pooling or a direct connection. The starter currently uses one DSN, so transaction pooling is not supported for the worker.

`/health/ready` checks the web process and database, not worker health. Supervise the worker as a separate service and monitor its exit status and queue growth. Compose binds River UI and Mailpit to loopback for development. A production River UI needs its own authentication and network-access configuration.

## Session revocation upgrade

Migration `003_login_sessions.sql` separates authentication grants from SCS data. Apply migrations before starting the updated server. Existing cookies have no grant and require a fresh sign-in; no old `userID` session value is promoted to authorization. Replace all old server instances before relying on the new revocation contract, since old binaries still trust SCS data. Keep the new table when rolling back application code, and treat rollback to the old authorization code as restoring the original concurrency defect.
