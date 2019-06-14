FROM mysql:5.6.41
ENV MYSQL_ROOT_PASSWORD password
ENV MYSQL_DATABASE poppulodb
ENV MYSQL_USER testuser
ENV MYSQL_PASSWORD testuser
COPY poppulodb.sql /docker-entrypoint-initdb.d/poppulodb.sql
EXPOSE 3306