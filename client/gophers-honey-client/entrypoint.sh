#!/bin/sh
JSON_STRING='window.configs = { \
  "VUE_APP_API_ROOT":"'"${VUE_APP_API_ROOT}"'" \
}'
sed -i "s@// CONFIGURATIONS_PLACEHOLDER@${JSON_STRING}@" /app/index.html
exec "$@"