# Go-C++ Interface

## Overview

The Go-C++ Interface is a critical component of the trading platform that enables seamless communication between the Go backend and the high-performance C++ order execution engine. This document details the design, implementation, and usage of this interface, focusing on performance, reliability, and maintainability.

## Architecture

The Go-C++ Interface follows a layered architecture to provide clean separation of concerns and efficient cross-language communication:

```
┌─────────────────────────┐
│      Go Backend         │
├─────────────────────────┤
│   Go Interface Layer    │
├─────────────────────────┤
│ Communication Protocol  │
├─────────────────────────┤
│   C++ Interface Layer   │
├─────────────────────────┤
│  C++ Execution Engine   │
└─────────────────────────┘
```

### Key Components

1. **Go Interface Layer**
   - Go code that interacts with the C++ components
   - Handles marshaling/unmarshaling of data
   - Manages resource lifecycle
   - Provides idiomatic Go API

2. **Communication Protocol**
   - Defines message formats and serialization
   - Handles cross-language data representation
   - Manages error propagation
   - Optimizes for performance

3. **C++ Interface Layer**
   - C++ code exposed to Go through CGO
   - Manages memory and resources
   - Translates between C++ and C-compatible types
   - Handles threading and synchronization

## Communication Mechanism

The Go-C++ Interface uses a combination of techniques to achieve optimal performance and reliability:

### CGO Integration

The primary method of communication is through CGO (cgo), which allows Go code to call C/C++ functions:

```go
// go_interface.go
package execution

/*
#cgo CFLAGS: -I${SRCDIR}/../../cpp/include
#cgo LDFLAGS: -L${SRCDIR}/../../cpp/build -lgo_interface -Wl,-rpath,${SRCDIR}/../../cpp/build
#include "interface/go_interface.h"
*/
import "C"
import (
    "unsafe"
    "errors"
    "sync"
)

// OrderBook represents a wrapper around the C++ order book
type OrderBook struct {
    handle C.OrderBookHandle
    mutex  sync.Mutex
}

// NewOrderBook creates a new order book instance
func NewOrderBook() (*OrderBook, error) {
    handle := C.order_book_create()
    if handle == nil {
        return nil, errors.New("failed to create order book")
    }
    return &OrderBook{handle: handle}, nil
}

// AddOrder adds an order to the order book
func (ob *OrderBook) AddOrder(order Order) (uint64, error) {
    ob.mutex.Lock()
    defer ob.mutex.Unlock()
    
    cOrder := orderToC(order)
    var errorMsg *C.char
    orderId := C.order_book_add_order(ob.handle, cOrder, &errorMsg)
    
    if errorMsg != nil {
        err := errors.New(C.GoString(errorMsg))
        C.free_string(errorMsg)
        return 0, err
    }
    
    return uint64(orderId), nil
}

// Close releases resources associated with the order book
func (ob *OrderBook) Close() {
    ob.mutex.Lock()
    defer ob.mutex.Unlock()
    
    if ob.handle != nil {
        C.order_book_destroy(ob.handle)
        ob.handle = nil
    }
}
```

### C++ Header Interface

The C++ side exposes a C-compatible API that can be called from Go:

```cpp
// go_interface.h
#ifndef GO_INTERFACE_H
#define GO_INTERFACE_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct OrderBookStruct* OrderBookHandle;
typedef struct OrderStruct {
    const char* symbol;
    double price;
    int quantity;
    int side;  // 0 = buy, 1 = sell
    int type;  // 0 = market, 1 = limit
} Order;

// Order book functions
OrderBookHandle order_book_create();
void order_book_destroy(OrderBookHandle handle);
unsigned long long order_book_add_order(OrderBookHandle handle, Order order, char** error_msg);
int order_book_cancel_order(OrderBookHandle handle, unsigned long long order_id, char** error_msg);
int order_book_get_depth(OrderBookHandle handle, const char* symbol, int side, int levels, double* prices, int* quantities, char** error_msg);

// Memory management
void free_string(char* str);

#ifdef __cplusplus
}
#endif

#endif // GO_INTERFACE_H
```

### C++ Implementation

The C++ implementation wraps the internal C++ classes with C-compatible functions:

```cpp
// go_interface.cpp
#include "interface/go_interface.h"
#include "order_book/order_book.h"
#include "common/error_handling.h"
#include <string>
#include <cstring>
#include <memory>
#include <unordered_map>
#include <mutex>

struct OrderBookStruct {
    std::unique_ptr<OrderBook> impl;
    std::mutex mutex;
};

static std::string lastErrorMessage;

// Helper function to allocate and return error message
static char* createErrorMessage(const std::string& message) {
    lastErrorMessage = message;
    char* result = static_cast<char*>(malloc(lastErrorMessage.size() + 1));
    if (result) {
        strcpy(result, lastErrorMessage.c_str());
    }
    return result;
}

extern "C" {

OrderBookHandle order_book_create() {
    try {
        auto* handle = new OrderBookStruct();
        handle->impl = std::make_unique<OrderBook>();
        return handle;
    } catch (const std::exception& e) {
        return nullptr;
    }
}

void order_book_destroy(OrderBookHandle handle) {
    if (handle) {
        delete handle;
    }
}

unsigned long long order_book_add_order(OrderBookHandle handle, Order order, char** error_msg) {
    if (!handle || !handle->impl) {
        *error_msg = createErrorMessage("Invalid order book handle");
        return 0;
    }
    
    try {
        std::lock_guard<std::mutex> lock(handle->mutex);
        
        OrderSide side = order.side == 0 ? OrderSide::Buy : OrderSide::Sell;
        OrderType type = order.type == 0 ? OrderType::Market : OrderType::Limit;
        
        OrderParams params{
            order.symbol,
            order.price,
            order.quantity,
            side,
            type
        };
        
        return handle->impl->addOrder(params);
    } catch (const std::exception& e) {
        *error_msg = createErrorMessage(e.what());
        return 0;
    }
}

int order_book_cancel_order(OrderBookHandle handle, unsigned long long order_id, char** error_msg) {
    if (!handle || !handle->impl) {
        *error_msg = createErrorMessage("Invalid order book handle");
        return 0;
    }
    
    try {
        std::lock_guard<std::mutex> lock(handle->mutex);
        return handle->impl->cancelOrder(order_id) ? 1 : 0;
    } catch (const std::exception& e) {
        *error_msg = createErrorMessage(e.what());
        return 0;
    }
}

int order_book_get_depth(OrderBookHandle handle, const char* symbol, int side, int levels, 
                         double* prices, int* quantities, char** error_msg) {
    if (!handle || !handle->impl) {
        *error_msg = createErrorMessage("Invalid order book handle");
        return 0;
    }
    
    try {
        std::lock_guard<std::mutex> lock(handle->mutex);
        
        OrderSide orderSide = side == 0 ? OrderSide::Buy : OrderSide::Sell;
        auto depth = handle->impl->getMarketDepth(symbol, orderSide, levels);
        
        int count = std::min(static_cast<int>(depth.size()), levels);
        for (int i = 0; i < count; ++i) {
            prices[i] = depth[i].price;
            quantities[i] = depth[i].quantity;
        }
        
        return count;
    } catch (const std::exception& e) {
        *error_msg = createErrorMessage(e.what());
        return 0;
    }
}

void free_string(char* str) {
    if (str) {
        free(str);
    }
}

} // extern "C"
```

### Shared Memory Communication

For high-performance data exchange, the interface uses shared memory:

```cpp
// shared_memory.h
#ifndef SHARED_MEMORY_H
#define SHARED_MEMORY_H

#include <string>
#include <cstdint>

namespace trading {

class SharedMemoryRegion {
public:
    SharedMemoryRegion(const std::string& name, size_t size, bool create);
    ~SharedMemoryRegion();
    
    void* getData() const;
    size_t getSize() const;
    bool isValid() const;
    
private:
    struct Impl;
    std::unique_ptr<Impl> pImpl;
};

} // namespace trading

#endif // SHARED_MEMORY_H
```

```go
// shared_memory.go
package execution

/*
#cgo CFLAGS: -I${SRCDIR}/../../cpp/include
#cgo LDFLAGS: -L${SRCDIR}/../../cpp/build -lshared_memory -Wl,-rpath,${SRCDIR}/../../cpp/build
#include "interface/shared_memory.h"
*/
import "C"
import (
    "unsafe"
    "errors"
)

// SharedMemory represents a shared memory region
type SharedMemory struct {
    handle C.SharedMemoryHandle
    data   unsafe.Pointer
    size   int
}

// NewSharedMemory creates or opens a shared memory region
func NewSharedMemory(name string, size int, create bool) (*SharedMemory, error) {
    cName := C.CString(name)
    defer C.free(unsafe.Pointer(cName))
    
    handle := C.shared_memory_create(cName, C.size_t(size), C.bool(create))
    if handle == nil {
        return nil, errors.New("failed to create shared memory")
    }
    
    data := C.shared_memory_get_data(handle)
    actualSize := int(C.shared_memory_get_size(handle))
    
    return &SharedMemory{
        handle: handle,
        data:   data,
        size:   actualSize,
    }, nil
}

// Close releases the shared memory region
func (sm *SharedMemory) Close() {
    if sm.handle != nil {
        C.shared_memory_destroy(sm.handle)
        sm.handle = nil
        sm.data = nil
        sm.size = 0
    }
}

// Data returns a byte slice representing the shared memory region
func (sm *SharedMemory) Data() []byte {
    if sm.data == nil {
        return nil
    }
    return unsafe.Slice((*byte)(sm.data), sm.size)
}
```

## Data Serialization

The interface uses efficient serialization to minimize overhead:

### Binary Protocol

For high-performance data exchange, a custom binary protocol is used:

```cpp
// serialization.h
#ifndef SERIALIZATION_H
#define SERIALIZATION_H

#include <vector>
#include <string>
#include <cstdint>

namespace trading {

class BinarySerializer {
public:
    BinarySerializer();
    
    // Write methods
    void writeInt8(int8_t value);
    void writeInt16(int16_t value);
    void writeInt32(int32_t value);
    void writeInt64(int64_t value);
    void writeUInt8(uint8_t value);
    void writeUInt16(uint16_t value);
    void writeUInt32(uint32_t value);
    void writeUInt64(uint64_t value);
    void writeFloat(float value);
    void writeDouble(double value);
    void writeString(const std::string& value);
    void writeBytes(const void* data, size_t size);
    
    // Get serialized data
    const std::vector<uint8_t>& getData() const;
    
private:
    std::vector<uint8_t> buffer;
};

class BinaryDeserializer {
public:
    BinaryDeserializer(const void* data, size_t size);
    
    // Read methods
    int8_t readInt8();
    int16_t readInt16();
    int32_t readInt32();
    int64_t readInt64();
    uint8_t readUInt8();
    uint16_t readUInt16();
    uint32_t readUInt32();
    uint64_t readUInt64();
    float readFloat();
    double readDouble();
    std::string readString();
    void readBytes(void* output, size_t size);
    
    // Check if more data is available
    bool hasMore() const;
    
private:
    const uint8_t* data;
    size_t size;
    size_t position;
};

} // namespace trading

#endif // SERIALIZATION_H
```

```go
// serialization.go
package execution

import (
    "encoding/binary"
    "math"
    "io"
    "bytes"
)

// BinarySerializer provides methods to serialize data in binary format
type BinarySerializer struct {
    buffer bytes.Buffer
}

// NewBinarySerializer creates a new binary serializer
func NewBinarySerializer() *BinarySerializer {
    return &BinarySerializer{}
}

// WriteInt8 writes an int8 value
func (s *BinarySerializer) WriteInt8(value int8) {
    s.buffer.WriteByte(byte(value))
}

// WriteInt16 writes an int16 value
func (s *BinarySerializer) WriteInt16(value int16) {
    b := make([]byte, 2)
    binary.LittleEndian.PutUint16(b, uint16(value))
    s.buffer.Write(b)
}

// WriteInt32 writes an int32 value
func (s *BinarySerializer) WriteInt32(value int32) {
    b := make([]byte, 4)
    binary.LittleEndian.PutUint32(b, uint32(value))
    s.buffer.Write(b)
}

// WriteInt64 writes an int64 value
func (s *BinarySerializer) WriteInt64(value int64) {
    b := make([]byte, 8)
    binary.LittleEndian.PutUint64(b, uint64(value))
    s.buffer.Write(b)
}

// WriteFloat64 writes a float64 value
func (s *BinarySerializer) WriteFloat64(value float64) {
    b := make([]byte, 8)
    binary.LittleEndian.PutUint64(b, math.Float64bits(value))
    s.buffer.Write(b)
}

// WriteString writes a string value
func (s *BinarySerializer) WriteString(value string) {
    s.WriteInt32(int32(len(value)))
    s.buffer.WriteString(value)
}

// Bytes returns the serialized data
func (s *BinarySerializer) Bytes() []byte {
    return s.buffer.Bytes()
}
```

### Protocol Buffers

For more complex data structures, Protocol Buffers are used:

```protobuf
// order.proto
syntax = "proto3";
package trading;

message Order {
    string symbol = 1;
    double price = 2;
    int32 quantity = 3;
    enum Side {
        BUY = 0;
        SELL = 1;
    }
    Side side = 4;
    enum Type {
        MARKET = 0;
        LIMIT = 1;
        STOP = 2;
        STOP_LIMIT = 3;
    }
    Type type = 5;
    uint64 client_order_id = 6;
    string user_id = 7;
    map<string, string> properties = 8;
}

message OrderBook {
    string symbol = 1;
    repeated PriceLevel bids = 2;
    repeated PriceLevel asks = 3;
    uint64 timestamp = 4;
}

message PriceLevel {
    double price = 1;
    int32 quantity = 2;
    int32 order_count = 3;
}

message ExecutionReport {
    uint64 order_id = 1;
    string symbol = 2;
    double price = 3;
    int32 quantity = 4;
    int32 filled_quantity = 5;
    Order.Side side = 6;
    enum Status {
        NEW = 0;
        PARTIALLY_FILLED = 1;
        FILLED = 2;
        CANCELED = 3;
        REJECTED = 4;
    }
    Status status = 7;
    string reason = 8;
    uint64 timestamp = 9;
}
```

## Error Handling

The interface implements robust error handling across language boundaries:

### Error Propagation

Errors from C++ are propagated to Go with detailed information:

```cpp
// error_handling.h
#ifndef ERROR_HANDLING_H
#define ERROR_HANDLING_H

#include <string>
#include <exception>

namespace trading {

class TradingException : public std::exception {
public:
    enum class ErrorCode {
        UNKNOWN_ERROR = 0,
        INVALID_ARGUMENT = 1,
        RESOURCE_NOT_FOUND = 2,
        PERMISSION_DENIED = 3,
        RESOURCE_EXHAUSTED = 4,
        INTERNAL_ERROR = 5
    };
    
    TradingException(ErrorCode code, const std::string& message);
    
    const char* what() const noexcept override;
    ErrorCode code() const noexcept;
    
private:
    ErrorCode errorCode;
    std::string errorMessage;
};

} // namespace trading

#endif // ERROR_HANDLING_H
```

```go
// errors.go
package execution

import (
    "fmt"
)

// ErrorCode represents error codes from the C++ execution engine
type ErrorCode int

const (
    ErrorUnknown ErrorCode = iota
    ErrorInvalidArgument
    ErrorResourceNotFound
    ErrorPermissionDenied
    ErrorResourceExhausted
    ErrorInternal
)

// ExecutionError represents an error from the C++ execution engine
type ExecutionError struct {
    Code    ErrorCode
    Message string
}

// Error implements the error interface
func (e *ExecutionError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewExecutionError creates a new execution error
func NewExecutionError(code ErrorCode, message string) *ExecutionError {
    return &ExecutionError{
        Code:    code,
        Message: message,
    }
}
```

### Recovery Mechanisms

The interface includes mechanisms to recover from errors and prevent crashes:

```cpp
// recovery.cpp
#include "interface/recovery.h"
#include <exception>
#include <iostream>

extern "C" {

int execute_with_recovery(void (*func)(void*), void* data, char** error_msg) {
    try {
        func(data);
        return 1;
    } catch (const std::exception& e) {
        if (error_msg) {
            *error_msg = strdup(e.what());
        }
        return 0;
    } catch (...) {
        if (error_msg) {
            *error_msg = strdup("Unknown error occurred");
        }
        return 0;
    }
}

} // extern "C"
```

```go
// recovery.go
package execution

/*
#cgo CFLAGS: -I${SRCDIR}/../../cpp/include
#cgo LDFLAGS: -L${SRCDIR}/../../cpp/build -lrecovery -Wl,-rpath,${SRCDIR}/../../cpp/build
#include "interface/recovery.h"
#include <stdlib.h>

extern void goCallback(void* data);

static int executeWithRecovery(void* data) {
    char* error_msg = NULL;
    int result = execute_with_recovery(goCallback, data, &error_msg);
    if (error_msg != NULL) {
        free(error_msg);
    }
    return result;
}
*/
import "C"
import (
    "runtime/cgo"
    "sync"
    "unsafe"
)

var (
    callbackRegistry = struct {
        sync.Mutex
        callbacks map[unsafe.Pointer]func()
    }{
        callbacks: make(map[unsafe.Pointer]func()),
    }
)

//export goCallback
func goCallback(data unsafe.Pointer) {
    callbackRegistry.Lock()
    callback, ok := callbackRegistry.callbacks[data]
    callbackRegistry.Unlock()
    
    if ok && callback != nil {
        callback()
    }
}

// ExecuteWithRecovery executes a function with C++ exception recovery
func ExecuteWithRecovery(callback func()) bool {
    handle := cgo.NewHandle(callback)
    defer handle.Delete()
    
    ptr := unsafe.Pointer(&handle)
    
    callbackRegistry.Lock()
    callbackRegistry.callbacks[ptr] = callback
    callbackRegistry.Unlock()
    
    defer func() {
        callbackRegistry.Lock()
        delete(callbackRegistry.callbacks, ptr)
        callbackRegistry.Unlock()
    }()
    
    return C.executeWithRecovery(ptr) != 0
}
```

## Resource Management

The interface implements careful resource management to prevent leaks:

### Memory Management

Memory is managed carefully across the language boundary:

```cpp
// memory_management.cpp
#include "interface/memory_management.h"
#include <unordered_map>
#include <mutex>
#include <memory>

namespace {
    std::mutex resourceMutex;
    std::unordered_map<uint64_t, std::shared_ptr<void>> resources;
    uint64_t nextResourceId = 1;
}

extern "C" {

uint64_t register_resource(void* resource, void (*deleter)(void*)) {
    if (!resource || !deleter) {
        return 0;
    }
    
    std::lock_guard<std::mutex> lock(resourceMutex);
    uint64_t id = nextResourceId++;
    
    resources[id] = std::shared_ptr<void>(resource, deleter);
    return id;
}

void* get_resource(uint64_t id) {
    std::lock_guard<std::mutex> lock(resourceMutex);
    auto it = resources.find(id);
    if (it == resources.end()) {
        return nullptr;
    }
    return it->second.get();
}

int release_resource(uint64_t id) {
    std::lock_guard<std::mutex> lock(resourceMutex);
    return resources.erase(id) > 0 ? 1 : 0;
}

} // extern "C"
```

```go
// resource_manager.go
package execution

/*
#cgo CFLAGS: -I${SRCDIR}/../../cpp/include
#cgo LDFLAGS: -L${SRCDIR}/../../cpp/build -lmemory_management -Wl,-rpath,${SRCDIR}/../../cpp/build
#include "interface/memory_management.h"
*/
import "C"
import (
    "sync"
    "runtime"
)

// ResourceManager manages C++ resources from Go
type ResourceManager struct {
    resources map[uint64]interface{}
    mutex     sync.Mutex
}

// NewResourceManager creates a new resource manager
func NewResourceManager() *ResourceManager {
    return &ResourceManager{
        resources: make(map[uint64]interface{}),
    }
}

// Register registers a resource with the manager
func (rm *ResourceManager) Register(resource interface{}, id uint64) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()
    
    rm.resources[id] = resource
    
    // Ensure resource is released when garbage collected
    runtime.SetFinalizer(resource, func(r interface{}) {
        rm.Release(id)
    })
}

// Get retrieves a resource by ID
func (rm *ResourceManager) Get(id uint64) interface{} {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()
    
    return rm.resources[id]
}

// Release releases a resource by ID
func (rm *ResourceManager) Release(id uint64) bool {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()
    
    if _, exists := rm.resources[id]; exists {
        delete(rm.resources, id)
        C.release_resource(C.uint64_t(id))
        return true
    }
    
    return false
}
```

### Lifecycle Management

The interface manages the lifecycle of C++ objects:

```go
// lifecycle.go
package execution

import (
    "sync"
    "runtime"
)

// Lifecycle manages the lifecycle of C++ objects
type Lifecycle struct {
    initialized bool
    closed      bool
    mutex       sync.Mutex
    closeFunc   func()
}

// Initialize initializes the lifecycle
func (l *Lifecycle) Initialize(closeFunc func()) bool {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    if l.initialized {
        return false
    }
    
    l.initialized = true
    l.closed = false
    l.closeFunc = closeFunc
    
    // Set finalizer to ensure Close is called during garbage collection
    runtime.SetFinalizer(l, func(lifecycle *Lifecycle) {
        lifecycle.Close()
    })
    
    return true
}

// Close closes the lifecycle
func (l *Lifecycle) Close() bool {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    if !l.initialized || l.closed {
        return false
    }
    
    if l.closeFunc != nil {
        l.closeFunc()
    }
    
    l.closed = true
    runtime.SetFinalizer(l, nil)
    
    return true
}

// IsInitialized returns whether the lifecycle is initialized
func (l *Lifecycle) IsInitialized() bool {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    return l.initialized
}

// IsClosed returns whether the lifecycle is closed
func (l *Lifecycle) IsClosed() bool {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    return l.closed
}
```

## Performance Monitoring

The interface includes performance monitoring capabilities:

```cpp
// performance_metrics.h
#ifndef PERFORMANCE_METRICS_H
#define PERFORMANCE_METRICS_H

#include <string>
#include <chrono>
#include <atomic>
#include <vector>
#include <unordered_map>
#include <mutex>

namespace trading {

class PerformanceMetrics {
public:
    static PerformanceMetrics& getInstance();
    
    void recordLatency(const std::string& operation, std::chrono::nanoseconds latency);
    void incrementCounter(const std::string& counter, uint64_t value = 1);
    
    struct LatencyStats {
        double min;
        double max;
        double avg;
        double p50;
        double p95;
        double p99;
        uint64_t count;
    };
    
    LatencyStats getLatencyStats(const std::string& operation);
    uint64_t getCounter(const std::string& counter);
    
    std::unordered_map<std::string, LatencyStats> getAllLatencyStats();
    std::unordered_map<std::string, uint64_t> getAllCounters();
    
private:
    PerformanceMetrics();
    ~PerformanceMetrics();
    
    struct LatencyData {
        std::mutex mutex;
        std::vector<double> samples;
        double min = std::numeric_limits<double>::max();
        double max = 0;
        double sum = 0;
        uint64_t count = 0;
    };
    
    std::unordered_map<std::string, LatencyData> latencyData;
    std::unordered_map<std::string, std::atomic<uint64_t>> counters;
    std::mutex metricsMutex;
};

class ScopedLatencyRecorder {
public:
    ScopedLatencyRecorder(const std::string& operation);
    ~ScopedLatencyRecorder();
    
private:
    std::string operation;
    std::chrono::high_resolution_clock::time_point startTime;
};

} // namespace trading

#endif // PERFORMANCE_METRICS_H
```

```go
// performance_metrics.go
package execution

/*
#cgo CFLAGS: -I${SRCDIR}/../../cpp/include
#cgo LDFLAGS: -L${SRCDIR}/../../cpp/build -lperformance_metrics -Wl,-rpath,${SRCDIR}/../../cpp/build
#include "interface/performance_metrics.h"
*/
import "C"
import (
    "time"
    "unsafe"
)

// LatencyStats represents latency statistics
type LatencyStats struct {
    Min   float64
    Max   float64
    Avg   float64
    P50   float64
    P95   float64
    P99   float64
    Count uint64
}

// GetLatencyStats retrieves latency statistics for an operation
func GetLatencyStats(operation string) LatencyStats {
    cOperation := C.CString(operation)
    defer C.free(unsafe.Pointer(cOperation))
    
    var stats C.LatencyStats
    C.get_latency_stats(cOperation, &stats)
    
    return LatencyStats{
        Min:   float64(stats.min),
        Max:   float64(stats.max),
        Avg:   float64(stats.avg),
        P50:   float64(stats.p50),
        P95:   float64(stats.p95),
        P99:   float64(stats.p99),
        Count: uint64(stats.count),
    }
}

// GetCounter retrieves a counter value
func GetCounter(counter string) uint64 {
    cCounter := C.CString(counter)
    defer C.free(unsafe.Pointer(cCounter))
    
    return uint64(C.get_counter(cCounter))
}

// MeasureLatency measures the latency of a function
func MeasureLatency(operation string, f func()) {
    cOperation := C.CString(operation)
    defer C.free(unsafe.Pointer(cOperation))
    
    start := time.Now()
    f()
    latency := time.Since(start)
    
    C.record_latency(cOperation, C.long(latency.Nanoseconds()))
}
```

## Thread Safety

The interface ensures thread safety for concurrent access:

```cpp
// thread_safety.cpp
#include "interface/thread_safety.h"
#include <mutex>
#include <shared_mutex>
#include <atomic>
#include <thread>
#include <condition_variable>

namespace trading {

class ThreadSafeCounter {
public:
    ThreadSafeCounter() : value(0) {}
    
    void increment() {
        ++value;
    }
    
    uint64_t get() const {
        return value.load();
    }
    
private:
    std::atomic<uint64_t> value;
};

class ThreadSafeQueue {
public:
    void push(void* item) {
        std::lock_guard<std::mutex> lock(mutex);
        queue.push_back(item);
        cv.notify_one();
    }
    
    void* pop() {
        std::unique_lock<std::mutex> lock(mutex);
        cv.wait(lock, [this] { return !queue.empty(); });
        
        void* item = queue.front();
        queue.pop_front();
        return item;
    }
    
    bool tryPop(void** item) {
        std::lock_guard<std::mutex> lock(mutex);
        if (queue.empty()) {
            return false;
        }
        
        *item = queue.front();
        queue.pop_front();
        return true;
    }
    
    size_t size() const {
        std::lock_guard<std::mutex> lock(mutex);
        return queue.size();
    }
    
private:
    std::deque<void*> queue;
    mutable std::mutex mutex;
    std::condition_variable cv;
};

} // namespace trading
```

```go
// thread_safety.go
package execution

/*
#cgo CFLAGS: -I${SRCDIR}/../../cpp/include
#cgo LDFLAGS: -L${SRCDIR}/../../cpp/build -lthread_safety -Wl,-rpath,${SRCDIR}/../../cpp/build
#include "interface/thread_safety.h"
*/
import "C"
import (
    "unsafe"
    "sync"
)

// ThreadSafeCounter represents a thread-safe counter
type ThreadSafeCounter struct {
    handle C.CounterHandle
    mutex  sync.Mutex
}

// NewThreadSafeCounter creates a new thread-safe counter
func NewThreadSafeCounter() *ThreadSafeCounter {
    handle := C.counter_create()
    return &ThreadSafeCounter{handle: handle}
}

// Increment increments the counter
func (c *ThreadSafeCounter) Increment() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if c.handle != nil {
        C.counter_increment(c.handle)
    }
}

// Get gets the counter value
func (c *ThreadSafeCounter) Get() uint64 {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if c.handle != nil {
        return uint64(C.counter_get(c.handle))
    }
    return 0
}

// Close releases resources associated with the counter
func (c *ThreadSafeCounter) Close() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if c.handle != nil {
        C.counter_destroy(c.handle)
        c.handle = nil
    }
}
```

## Testing Strategy

The Go-C++ Interface is thoroughly tested to ensure reliability:

### Unit Testing

Unit tests verify individual components:

```go
// go_interface_test.go
package execution

import (
    "testing"
    "sync"
)

func TestOrderBookCreateDestroy(t *testing.T) {
    ob, err := NewOrderBook()
    if err != nil {
        t.Fatalf("Failed to create order book: %v", err)
    }
    
    ob.Close()
}

func TestOrderBookAddOrder(t *testing.T) {
    ob, err := NewOrderBook()
    if err != nil {
        t.Fatalf("Failed to create order book: %v", err)
    }
    defer ob.Close()
    
    order := Order{
        Symbol:   "AAPL",
        Price:    150.0,
        Quantity: 100,
        Side:     SideBuy,
        Type:     TypeLimit,
    }
    
    orderId, err := ob.AddOrder(order)
    if err != nil {
        t.Fatalf("Failed to add order: %v", err)
    }
    
    if orderId == 0 {
        t.Fatal("Expected non-zero order ID")
    }
}

func TestOrderBookConcurrency(t *testing.T) {
    ob, err := NewOrderBook()
    if err != nil {
        t.Fatalf("Failed to create order book: %v", err)
    }
    defer ob.Close()
    
    const numGoroutines = 10
    const numOrdersPerGoroutine = 100
    
    var wg sync.WaitGroup
    wg.Add(numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < numOrdersPerGoroutine; j++ {
                order := Order{
                    Symbol:   "AAPL",
                    Price:    150.0 + float64(j) * 0.1,
                    Quantity: 100 + j,
                    Side:     SideBuy,
                    Type:     TypeLimit,
                }
                
                _, err := ob.AddOrder(order)
                if err != nil {
                    t.Errorf("Failed to add order: %v", err)
                }
            }
        }(i)
    }
    
    wg.Wait()
}
```

### Integration Testing

Integration tests verify cross-language communication:

```go
// integration_test.go
package execution

import (
    "testing"
    "time"
    "math/rand"
)

func TestEndToEndOrderExecution(t *testing.T) {
    // Create order book
    ob, err := NewOrderBook()
    if err != nil {
        t.Fatalf("Failed to create order book: %v", err)
    }
    defer ob.Close()
    
    // Create matching engine
    me, err := NewMatchingEngine()
    if err != nil {
        t.Fatalf("Failed to create matching engine: %v", err)
    }
    defer me.Close()
    
    // Connect order book to matching engine
    if err := me.RegisterOrderBook(ob); err != nil {
        t.Fatalf("Failed to register order book: %v", err)
    }
    
    // Create execution listener
    listener := NewExecutionListener()
    defer listener.Close()
    
    // Register listener with matching engine
    if err := me.RegisterListener(listener); err != nil {
        t.Fatalf("Failed to register listener: %v", err)
    }
    
    // Add buy order
    buyOrder := Order{
        Symbol:   "AAPL",
        Price:    150.0,
        Quantity: 100,
        Side:     SideBuy,
        Type:     TypeLimit,
    }
    
    buyOrderId, err := ob.AddOrder(buyOrder)
    if err != nil {
        t.Fatalf("Failed to add buy order: %v", err)
    }
    
    // Add sell order that should match
    sellOrder := Order{
        Symbol:   "AAPL",
        Price:    150.0,
        Quantity: 100,
        Side:     SideSell,
        Type:     TypeLimit,
    }
    
    sellOrderId, err := ob.AddOrder(sellOrder)
    if err != nil {
        t.Fatalf("Failed to add sell order: %v", err)
    }
    
    // Wait for execution reports
    reports := listener.WaitForReports(2, 1*time.Second)
    if len(reports) != 2 {
        t.Fatalf("Expected 2 execution reports, got %d", len(reports))
    }
    
    // Verify execution reports
    for _, report := range reports {
        if report.OrderId != buyOrderId && report.OrderId != sellOrderId {
            t.Errorf("Unexpected order ID in execution report: %d", report.OrderId)
        }
        
        if report.Status != StatusFilled {
            t.Errorf("Expected order status FILLED, got %v", report.Status)
        }
        
        if report.FilledQuantity != 100 {
            t.Errorf("Expected filled quantity 100, got %d", report.FilledQuantity)
        }
    }
}
```

### Performance Testing

Performance tests verify the efficiency of the interface:

```go
// performance_test.go
package execution

import (
    "testing"
    "time"
    "sync"
    "math/rand"
)

func BenchmarkOrderAddition(b *testing.B) {
    ob, err := NewOrderBook()
    if err != nil {
        b.Fatalf("Failed to create order book: %v", err)
    }
    defer ob.Close()
    
    symbols := []string{"AAPL", "MSFT", "GOOG", "AMZN", "FB"}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        symbol := symbols[i%len(symbols)]
        price := 100.0 + float64(i%100)
        quantity := 100 + i%900
        side := SideBuy
        if i%2 == 0 {
            side = SideSell
        }
        
        order := Order{
            Symbol:   symbol,
            Price:    price,
            Quantity: quantity,
            Side:     side,
            Type:     TypeLimit,
        }
        
        _, err := ob.AddOrder(order)
        if err != nil {
            b.Fatalf("Failed to add order: %v", err)
        }
    }
}

func BenchmarkConcurrentOrderAddition(b *testing.B) {
    ob, err := NewOrderBook()
    if err != nil {
        b.Fatalf("Failed to create order book: %v", err)
    }
    defer ob.Close()
    
    symbols := []string{"AAPL", "MSFT", "GOOG", "AMZN", "FB"}
    
    b.ResetTimer()
    
    b.RunParallel(func(pb *testing.PB) {
        localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
        
        for pb.Next() {
            symbol := symbols[localRand.Intn(len(symbols))]
            price := 100.0 + float64(localRand.Intn(100))
            quantity := 100 + localRand.Intn(900)
            side := SideBuy
            if localRand.Intn(2) == 0 {
                side = SideSell
            }
            
            order := Order{
                Symbol:   symbol,
                Price:    price,
                Quantity: quantity,
                Side:     side,
                Type:     TypeLimit,
            }
            
            _, err := ob.AddOrder(order)
            if err != nil {
                b.Fatalf("Failed to add order: %v", err)
            }
        }
    })
}
```

## Deployment Considerations

The Go-C++ Interface requires careful deployment configuration:

### Library Packaging

The C++ libraries are packaged for deployment:

```cmake
# packaging.cmake
include(CPack)

set(CPACK_PACKAGE_NAME "trading-platform-cpp")
set(CPACK_PACKAGE_VERSION "1.0.0")
set(CPACK_PACKAGE_DESCRIPTION_SUMMARY "Trading Platform C++ Components")
set(CPACK_PACKAGE_VENDOR "Trading Platform")

set(CPACK_GENERATOR "TGZ")
set(CPACK_SOURCE_GENERATOR "TGZ")

set(CPACK_DEBIAN_PACKAGE_DEPENDS "libboost-all-dev, libtbb-dev")
```

### Runtime Configuration

Runtime configuration for the interface:

```go
// config.go
package execution

import (
    "os"
    "path/filepath"
    "runtime"
)

func init() {
    // Set library search path
    if runtime.GOOS == "linux" {
        // Get the directory of the current file
        _, file, _, _ := runtime.Caller(0)
        dir := filepath.Dir(file)
        
        // Set LD_LIBRARY_PATH to include the cpp/build directory
        cppBuildDir := filepath.Join(dir, "..", "..", "cpp", "build")
        
        // Append to existing LD_LIBRARY_PATH
        ldLibraryPath := os.Getenv("LD_LIBRARY_PATH")
        if ldLibraryPath == "" {
            os.Setenv("LD_LIBRARY_PATH", cppBuildDir)
        } else {
            os.Setenv("LD_LIBRARY_PATH", ldLibraryPath+":"+cppBuildDir)
        }
    }
}
```

## Troubleshooting Guide

Common issues and their solutions:

### Library Loading Issues

```go
// troubleshooting.go
package execution

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

// CheckLibraryPath checks if the C++ libraries are in the library path
func CheckLibraryPath() (bool, string) {
    if runtime.GOOS == "linux" {
        ldLibraryPath := os.Getenv("LD_LIBRARY_PATH")
        if ldLibraryPath == "" {
            return false, "LD_LIBRARY_PATH is not set"
        }
        
        // Get the directory of the current file
        _, file, _, _ := runtime.Caller(0)
        dir := filepath.Dir(file)
        
        // Check if cpp/build directory is in LD_LIBRARY_PATH
        cppBuildDir := filepath.Join(dir, "..", "..", "cpp", "build")
        cppBuildDirAbs, _ := filepath.Abs(cppBuildDir)
        
        paths := strings.Split(ldLibraryPath, ":")
        for _, path := range paths {
            pathAbs, _ := filepath.Abs(path)
            if pathAbs == cppBuildDirAbs {
                return true, ""
            }
        }
        
        return false, fmt.Sprintf("cpp/build directory (%s) is not in LD_LIBRARY_PATH", cppBuildDirAbs)
    }
    
    return true, ""
}

// CheckLibraryExists checks if a specific library exists
func CheckLibraryExists(libraryName string) (bool, string) {
    if runtime.GOOS == "linux" {
        // Get the directory of the current file
        _, file, _, _ := runtime.Caller(0)
        dir := filepath.Dir(file)
        
        // Check if library exists in cpp/build directory
        cppBuildDir := filepath.Join(dir, "..", "..", "cpp", "build")
        libraryPath := filepath.Join(cppBuildDir, fmt.Sprintf("lib%s.so", libraryName))
        
        if _, err := os.Stat(libraryPath); os.IsNotExist(err) {
            return false, fmt.Sprintf("Library %s not found at %s", libraryName, libraryPath)
        }
        
        return true, ""
    }
    
    return true, ""
}
```

### Common Errors

```go
// errors_guide.go
package execution

import (
    "fmt"
    "strings"
)

// CommonErrorGuide provides guidance for common errors
func CommonErrorGuide(err error) string {
    if err == nil {
        return ""
    }
    
    errStr := err.Error()
    
    if strings.Contains(errStr, "cannot open shared object file") {
        return `
Library loading error: The C++ shared library could not be found.

Possible solutions:
1. Ensure the C++ components are built: cd cpp && mkdir -p build && cd build && cmake .. && make
2. Check library path: go run cmd/check_libraries/main.go
3. Set LD_LIBRARY_PATH manually: export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/path/to/trading-platform/cpp/build
`
    }
    
    if strings.Contains(errStr, "Invalid order book handle") {
        return `
Invalid handle error: The C++ object handle is invalid.

Possible solutions:
1. Ensure the object is created successfully before use
2. Check if the object has been closed
3. Verify that the object is not being used concurrently without proper synchronization
`
    }
    
    if strings.Contains(errStr, "symbol not found") {
        return `
Symbol not found error: The requested symbol does not exist in the order book.

Possible solutions:
1. Verify the symbol name is correct
2. Ensure the symbol has been added to the order book
3. Check if the symbol has been removed or expired
`
    }
    
    return ""
}
```

## Best Practices

Guidelines for using the Go-C++ Interface effectively:

### Performance Optimization

```go
// best_practices.go
package execution

import (
    "runtime"
    "sync"
)

// BestPractices provides guidance for using the Go-C++ interface effectively
type BestPractices struct{}

// AvoidExcessiveCrossing recommends minimizing cross-language calls
func (BestPractices) AvoidExcessiveCrossing() string {
    return `
Minimize Cross-Language Calls:
- Batch operations when possible
- Process data in bulk rather than individual items
- Keep computation-intensive operations on the C++ side
- Use shared memory for large data transfers
`
}

// OptimizeMemoryUsage recommends memory optimization techniques
func (BestPractices) OptimizeMemoryUsage() string {
    return `
Optimize Memory Usage:
- Reuse objects instead of creating new ones
- Use object pools for frequently created/destroyed objects
- Release C++ resources explicitly when done
- Be aware of garbage collection timing
`
}

// EnsureThreadSafety recommends thread safety practices
func (BestPractices) EnsureThreadSafety() string {
    return `
Ensure Thread Safety:
- Use proper synchronization when accessing shared resources
- Be aware of Go's goroutine scheduling
- Avoid sharing C++ objects across goroutines without synchronization
- Use thread-safe data structures when appropriate
`
}

// MonitorPerformance recommends performance monitoring practices
func (BestPractices) MonitorPerformance() string {
    return `
Monitor Performance:
- Use the performance monitoring tools provided
- Track latency of cross-language calls
- Profile memory usage regularly
- Benchmark critical paths
`
}
```

## Conclusion

The Go-C++ Interface is a critical component of the trading platform that enables seamless communication between the Go backend and the high-performance C++ order execution engine. By carefully designing for performance, reliability, and maintainability, this interface provides the best of both worlds: the development speed and ecosystem of Go combined with the raw performance of C++.

The interface uses a combination of techniques including CGO, shared memory, and efficient serialization to minimize overhead while maintaining safety and reliability. Comprehensive error handling, resource management, and performance monitoring ensure robust operation in production environments.

By following the guidelines and best practices outlined in this document, developers can effectively utilize the Go-C++ Interface to build high-performance trading applications that meet the demanding requirements of modern financial markets.
