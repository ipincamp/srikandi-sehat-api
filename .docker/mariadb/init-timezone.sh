#!/bin/sh
set -e
mysql_tzinfo_to_sql /usr/share/zoneinfo | mysql -u root -p"$MYSQL_ROOT_PASSWORD" mysql
echo ">>>> Timezone data successfully loaded into MariaDB."
