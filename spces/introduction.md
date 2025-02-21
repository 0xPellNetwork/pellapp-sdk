# PellApp-sdk Introduction

## What is the PellApp SDK

The PellApp SDK is a task processing framework similar to the Cosmos SDK, designed specifically for developing and extending decentralized verification service (DVS) applications based on pelldvs. With the PellApp SDK, DVS developers can easily integrate their applications and related business logic without needing to deeply understand the underlying implementation details of pelldvs.

### Key Features

- **Comprehensive RPC Function Integration**: The PellApp SDK integrates all RPC functions of DVS, allowing developers to focus solely on implementing business logic.
- **Unified DVS Message Protocol**: Developers only need to register specific DVS messages, and the PellApp SDK will automatically handle all DVS requests.
- **Simplified Request and Response Handling**: By registering the corresponding message handler in `RegisterMsgHandler`, developers can easily handle all DVS requests and responses.

### Handling DVS Requests

1. **Modular Processing Steps**:
   - **Message Processor**: Each DVS request message corresponds to a processor registered in `MsgRouterMgr`. The processor maps protocol messages (proto msg) to the corresponding processing function.
   - **Message Execution Result Processor**: After each message processor executes, it returns an execution result containing event information from the execution process. The `ResultHandler` is responsible for handling these execution results.

### Handling DVS Responses

- **Message Processor**: Each DVS response message corresponds to a processor registered in `MsgRouterMgr`. The processor maps protocol messages (proto msg) to the corresponding processing function.
- **Response Message Execution Result Processor**: After each message processor executes, it returns an execution result containing event information from the execution process. The `ResultHandler` is responsible for handling these execution results.

### Summary

The PellApp SDK provides an efficient, modular way to handle requests and responses in DVS applications. With its powerful features and flexible architecture, developers can quickly build and extend decentralized verification service applications, driving innovation and development in blockchain technology.
