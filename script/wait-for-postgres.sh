# ejecutar para este archivo:
# chmod +x script/wait-for-postgres.sh
#!/bin/sh

# Esperar hasta que postgres esté disponible
echo "Waiting for postgres..."

max_attempts=30
attempt=1
# Imprime las variables de entorno para verificar
echo "DB_POSTGRES_USER: $DB_POSTGRES_USER"
echo "DB_POSTGRES_HOST: $DB_POSTGRES_HOST"
echo "DB_POSTGRES_NAME: $DB_POSTGRES_NAME"
until PGPASSWORD=$DB_POSTGRES_PASS psql -h "$DB_POSTGRES_HOST" -U "$DB_POSTGRES_USER" -d "$DB_POSTGRES_NAME" -c '\q' >/dev/null 2>&1; do
    if [ $attempt -gt $max_attempts ]; then
        echo "Max attempts reached. Postgres is still unavailable."
        exit 1
    fi

    echo "Postgres is not yet available. Waiting..."
    attempt=$(( attempt + 1 ))
    sleep 60
done

echo "Postgres is available. Starting application..."

# Ejecutar tu aplicación
exec "$@"
