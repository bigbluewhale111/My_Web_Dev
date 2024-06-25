# Use an official node image as a base
FROM node:16-alpine as build

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY package.json ./
COPY package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application to the container
COPY . .

# Build the Vite React application
RUN npm run build

# Serve the built application using a simple HTTP server
FROM node:16-alpine
WORKDIR /app
RUN npm install -g serve
COPY --from=build /app/dist /app

# Expose port 8000
EXPOSE 8000

# Start the HTTP server
CMD ["npx", "serve", "-l", "8000"]