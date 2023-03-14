FROM golang:alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy the application files to the container's working directory
COPY . .

# Build the Go application inside the container
RUN go build -o main .

# Expose port 9090 to the outside world
EXPOSE 9070
# Command to run the executable
CMD ["./main"]