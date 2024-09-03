# Project Setup and Execution Guide

## Overview

This project involves interacting with a Jedi server and running a smart contract. Follow the steps below to build and run the necessary components and execute tests.

## Prerequisites

Ensure you have the following installed on your machine:
- Docker
- Node.js and npm (Node Package Manager)
- Hardhat (for running Ethereum smart contract tests)

## Step-by-Step Instructions

### 1. Start the Jedi Server

The Jedi server needs to be built and run using Docker. Follow these steps to start the server:

1. Navigate to the `go-jedi` directory:
   ```bash
   cd go-jedi
2. Build the Docker image for Jedi:
   ```bash
   docker build . -t jedi

3. Run the Docker container with the necessary volume mounts:
      ```bash
      docker run -v ./go:/go -v ./:/app -w /app jede
      ```
This command mounts the go directory and the current directory into the container, setting /app as the working directory.

### 2. Run the Smart Contract Tests
After starting the Jedi server, you need to interact with it using a smart contract. Follow these steps to set up and run the tests:

1. Navigate to the kyc-contract directory:
    ```bash
    cd kyc-contract
    ```

2. Install the required npm packages:
    ```bash
    npm install
    ```

3. Run the tests using Hardhat:
    ```bash
    npx hardhat test
    ```
    
### Customizing Tests
If you wish to modify the test data or the test cases, you can edit the kyc-contract/test/MedicalRecord-test.js file. Be sure to adjust the test data and logic according to your requirements.

