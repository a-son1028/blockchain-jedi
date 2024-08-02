# Use the official Node.js image from the Docker Hub with a version greater than 16
FROM node:16-alpine

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy package.json and package-lock.json into the container
COPY package*.json ./

# Install the project dependencies
RUN npm install

# Copy the rest of the application files into the container
COPY . .

# Set the default command to run your tests using Hardhat
CMD ["npx", "hardhat", "test"]
