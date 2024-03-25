# Derived from official mysql image (our base image)
FROM mysql:8.0.33

LABEL version="1.0.0"
LABEL maintainer="Topzone"

# Add a database
ENV MYSQL_DATABASE devops
ENV MYSQL_PORT 3306

# Add the content of the sql-scripts/ directory to your image
# All scripts in docker-entrypoint-initdb.d/ are automatically
# executed during container startup
COPY sql /docker-entrypoint-initdb.d/

