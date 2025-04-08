#!/usr/bin/env bash

if [ -z "$NEED_RECOVERY" ]; then
  echo "No recovery needed"
  exit 0
fi

cat << EOF > /tmp/minio.json
{
  "version": "10",
  "aliases": {
    "minio": {
      "url": "http://minio:9000",
      "accessKey": "$MINIO_USER",
      "secretKey": "$MINIO_PASSWORD",
      "auth": "s3v4",
      "path": "auto"
    }
  }
}
EOF

REDIS_ARGS=${REDIS_PASSWORD:+-a "$REDIS_PASSWORD"}

echo "Importing config to Redis..."
redis-cli -h redis $REDIS_ARGS -n 7 -x SET $ID < /recovery/config.json

echo "Importing database to MariaDB..."
mysql -h mariadb -u "$DB_USER" --password="$DB_PASS" gdps_$ID < /recovery/gdps_$ID.sql

echo "Importing savedata to Minio..."
mc cp -r /recovery/gdps_savedata minio/gdps_savedata

echo "Done!"