FROM postgres:9.6.6-alpine
ENV POSTGRES_USER postgres
ENV POSTGRES_DB todo
ENV POSTGRES_PASSWORD password
ENV POSTGRES_PORT 5432
EXPOSE 5432

ADD database.sql /docker-entrypoint-initdb.d/