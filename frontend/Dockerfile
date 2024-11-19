# Stage 1: Build
FROM node:18-alpine AS build

# Set working directory
WORKDIR /app

# Accept build args
ARG REACT_APP_BACKEND_URL

# Add the build arg to the environment
ENV REACT_APP_BACKEND_URL=$REACT_APP_BACKEND_URL

# Copy package.json and package-lock.json to install dependencies
COPY package.json package-lock.json ./
RUN npm install --frozen-lockfile

# Copy only necessary files to avoid .env and node_modules
COPY . .

# Remove .env if they were accidentally copied
RUN rm -f .env || true

# Build the React app
RUN npm run build

# Stage 2: Production environment to serve the app with Node.js
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy build output from the previous stage
COPY --from=build /app/build ./build

# Install a lightweight HTTP server to serve the build
RUN npm install -g serve

# Expose the port the app runs on
EXPOSE 3000

# Start the app with serve on port 3000
CMD ["serve", "-s", "build", "-l", "3000"]