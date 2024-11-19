#!/bin/sh

if [ "$ROLE" = "worker" ]; then
  touch .env
  ./main
else
  cat << EOF > .env
MYSQL_MIGRATIONS_FOLDER=/app/migrations/mysql/project_form_exam_sman2
EOF

  make mysql-migrate-up
  if [ $? -ne 0 ]; then
    echo "Make command failed"
    exit 1
  fi

  ./main
fi