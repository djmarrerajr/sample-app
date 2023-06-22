FROM cockroachdb/cockroach:v23.1.1

COPY ./init_crdb.sh /cockroach
RUN chmod +x /cockroach/init_crdb.sh

COPY ./bootstrap_crdb.sql /cockroach