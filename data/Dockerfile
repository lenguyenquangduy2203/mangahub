FROM alpine:3.19

RUN apk add --no-cache sqlite

WORKDIR /data

# Copy SQL schema and initialization script
COPY init.sql /data/init.sql
COPY init-db.sh /usr/local/bin/init-db.sh

# Make sure the script is executable
RUN chmod +x /usr/local/bin/init-db.sh

# Use exec form for CMD
CMD ["sh", "/usr/local/bin/init-db.sh"]