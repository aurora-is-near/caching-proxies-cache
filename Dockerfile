FROM alpine:latest

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY . .

# Ensure the binary is executable
RUN chmod +x /app/app