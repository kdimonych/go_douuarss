# This docker file is the base for each GO service
FROM service-base:latest

# Set working directory
WORKDIR /app

# Install dependencies
RUN mkdir -p /app/bin/dockers/news_checker

CMD ["bash", "/app/scripts/build_and_run.sh", "-b", "/app/bin/dockers/news_checker", "/app/services/news_checker/"]
