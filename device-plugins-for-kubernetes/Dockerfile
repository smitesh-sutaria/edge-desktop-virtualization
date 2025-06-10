# Use a minimal base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary into the container
COPY device-plugin /app/device-plugin

# Ensure the binary is executable
RUN chmod +x /app/device-plugin

# Run the device plugin
ENTRYPOINT ["/app/device-plugin"]

