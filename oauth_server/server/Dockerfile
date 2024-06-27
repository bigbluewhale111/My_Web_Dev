# Use the official Golang image as a base
FROM alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application to the container
COPY ./dist ./

COPY ./main ./

# Expose port 7000 (or whatever port your Go server uses)
EXPOSE 7000

RUN chmod +x ./main

# Command to run the executable
CMD ["./main"]