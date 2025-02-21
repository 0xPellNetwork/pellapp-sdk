# PellApp-sdk

## What is the PellApp SDK

The PellApp SDK is a task processing framework similar to the Cosmos SDK, designed specifically for developing and extending decentralized verification service (DVS) applications based on pelldvs. With the PellApp SDK, DVS developers can easily integrate their applications and related business logic without needing to deeply understand the underlying implementation details of pelldvs.

### Key Features

- **Comprehensive RPC Function Integration**: The PellApp SDK integrates all RPC functions of DVS, allowing developers to focus solely on implementing business logic.
- **Unified DVS Message Protocol**: Developers only need to register specific DVS messages, and the PellApp SDK will automatically handle all DVS requests.
- **Simplified Request and Response Handling**: By registering the corresponding message handler in `RegisterMsgHandler`, developers can easily handle all DVS requests and responses.

### Running Unit Tests

To ensure the integrity and functionality of your application, it's important to run unit tests. Follow these steps to execute the unit tests for the PellApp SDK:

1. **Install Go**: Ensure that Go is installed on your system. You can download it from [the official Go website](https://golang.org/dl/).

2. **Run Tests**: Execute the unit tests using the Go testing tool. Navigate to the root directory of your project and run:

   ```bash
   go test ./...
   ```

3. **Review Results**: Check the output for any failed tests and address any issues as needed.

### Summary

The PellApp SDK provides an efficient, modular way to handle requests and responses in DVS applications. With its powerful features and flexible architecture, developers can quickly build and extend decentralized verification service applications, driving innovation and development in blockchain technology.
