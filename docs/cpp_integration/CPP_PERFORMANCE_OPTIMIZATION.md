# C++ Performance Optimization

## Overview

This document provides comprehensive guidance on performance optimization techniques for the C++ order execution engine in the trading platform. Performance is critical for trading systems, especially in high-frequency trading scenarios where microseconds can make a significant difference in execution quality and profitability.

## Performance Goals

The C++ order execution engine aims to achieve the following performance targets:

- **Order Processing Latency**: < 10 microseconds (99th percentile)
- **Matching Engine Throughput**: > 1 million orders per second
- **Market Data Processing**: > 5 million updates per second
- **Memory Footprint**: < 2GB for full order book
- **CPU Utilization**: Efficient scaling across available cores

## Memory Optimization

### Memory-First Design

Memory access patterns are often the primary bottleneck in high-performance systems. The C++ order execution engine employs a memory-first design approach:

#### Custom Memory Allocators

```cpp
// memory_pool.h
#ifndef MEMORY_POOL_H
#define MEMORY_POOL_H

#include <cstddef>
#include <vector>
#include <mutex>
#include <memory>

namespace trading {

template<typename T, size_t BlockSize = 4096>
class MemoryPool {
public:
    MemoryPool() : currentBlock_(nullptr), currentSlot_(nullptr), lastSlot_(nullptr), freeSlots_(nullptr) {
        // Allocate first block
        allocateBlock();
    }
    
    ~MemoryPool() {
        // Free all allocated memory blocks
        for (auto block : blocks_) {
            std::free(block);
        }
    }
    
    // Allocate a new object
    T* allocate() {
        std::lock_guard<std::mutex> lock(mutex_);
        
        if (freeSlots_ != nullptr) {
            // Use a previously freed slot
            T* result = reinterpret_cast<T*>(freeSlots_);
            freeSlots_ = freeSlots_->next;
            return result;
        }
        
        if (currentSlot_ >= lastSlot_) {
            // Current block is full, allocate a new one
            allocateBlock();
        }
        
        return reinterpret_cast<T*>(currentSlot_++);
    }
    
    // Deallocate an object
    void deallocate(T* p) {
        std::lock_guard<std::mutex> lock(mutex_);
        
        if (p == nullptr) {
            return;
        }
        
        // Add the slot to the free list
        reinterpret_cast<Slot*>(p)->next = freeSlots_;
        freeSlots_ = reinterpret_cast<Slot*>(p);
    }
    
private:
    // Slot in the free list
    union Slot {
        T value;
        Slot* next;
    };
    
    // Allocate a new block of memory
    void allocateBlock() {
        // Calculate block size
        size_t blockSize = BlockSize * sizeof(Slot);
        
        // Allocate memory
        Slot* newBlock = reinterpret_cast<Slot*>(std::malloc(blockSize));
        if (newBlock == nullptr) {
            throw std::bad_alloc();
        }
        
        // Add block to the list
        blocks_.push_back(newBlock);
        
        // Update pointers
        currentBlock_ = newBlock;
        currentSlot_ = newBlock;
        lastSlot_ = newBlock + BlockSize;
    }
    
    Slot* currentBlock_;
    Slot* currentSlot_;
    Slot* lastSlot_;
    Slot* freeSlots_;
    std::vector<Slot*> blocks_;
    std::mutex mutex_;
};

} // namespace trading

#endif // MEMORY_POOL_H
```

#### Object Recycling

```cpp
// object_recycler.h
#ifndef OBJECT_RECYCLER_H
#define OBJECT_RECYCLER_H

#include <vector>
#include <mutex>
#include <memory>
#include <functional>

namespace trading {

template<typename T>
class ObjectRecycler {
public:
    ObjectRecycler(size_t initialCapacity = 1000, 
                  std::function<T*()> factory = []() { return new T(); },
                  std::function<void(T*)> reset = [](T* obj) { /* Default no-op */ }) 
        : factory_(factory), reset_(reset) {
        // Pre-allocate objects
        for (size_t i = 0; i < initialCapacity; ++i) {
            freeObjects_.push_back(factory_());
        }
    }
    
    ~ObjectRecycler() {
        // Free all objects
        for (auto obj : freeObjects_) {
            delete obj;
        }
    }
    
    // Get an object from the pool
    T* acquire() {
        std::lock_guard<std::mutex> lock(mutex_);
        
        if (freeObjects_.empty()) {
            // Create a new object if none available
            return factory_();
        }
        
        // Get object from the pool
        T* obj = freeObjects_.back();
        freeObjects_.pop_back();
        return obj;
    }
    
    // Return an object to the pool
    void release(T* obj) {
        if (obj == nullptr) {
            return;
        }
        
        // Reset the object state
        reset_(obj);
        
        std::lock_guard<std::mutex> lock(mutex_);
        freeObjects_.push_back(obj);
    }
    
    // Get number of available objects
    size_t available() const {
        std::lock_guard<std::mutex> lock(mutex_);
        return freeObjects_.size();
    }
    
private:
    std::vector<T*> freeObjects_;
    std::function<T*()> factory_;
    std::function<void(T*)> reset_;
    mutable std::mutex mutex_;
};

} // namespace trading

#endif // OBJECT_RECYCLER_H
```

#### Memory-Mapped Files

For large data sets or persistence, memory-mapped files provide efficient access:

```cpp
// memory_mapped_file.h
#ifndef MEMORY_MAPPED_FILE_H
#define MEMORY_MAPPED_FILE_H

#include <string>
#include <cstddef>
#include <stdexcept>
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

namespace trading {

class MemoryMappedFile {
public:
    MemoryMappedFile(const std::string& filename, bool readOnly = true) 
        : data_(nullptr), size_(0), readOnly_(readOnly), fd_(-1) {
        // Open file
        fd_ = open(filename.c_str(), readOnly ? O_RDONLY : O_RDWR | O_CREAT, 0644);
        if (fd_ == -1) {
            throw std::runtime_error("Failed to open file: " + filename);
        }
        
        // Get file size
        struct stat sb;
        if (fstat(fd_, &sb) == -1) {
            close(fd_);
            throw std::runtime_error("Failed to get file size");
        }
        
        size_ = sb.st_size;
        
        // Map file into memory
        data_ = mmap(nullptr, size_, readOnly ? PROT_READ : PROT_READ | PROT_WRITE, 
                    readOnly ? MAP_PRIVATE : MAP_SHARED, fd_, 0);
        
        if (data_ == MAP_FAILED) {
            close(fd_);
            throw std::runtime_error("Failed to map file into memory");
        }
    }
    
    ~MemoryMappedFile() {
        if (data_ != nullptr && data_ != MAP_FAILED) {
            munmap(data_, size_);
        }
        
        if (fd_ != -1) {
            close(fd_);
        }
    }
    
    // Get pointer to mapped memory
    void* getData() const {
        return data_;
    }
    
    // Get size of mapped memory
    size_t getSize() const {
        return size_;
    }
    
    // Flush changes to disk
    void flush() {
        if (!readOnly_ && data_ != nullptr && data_ != MAP_FAILED) {
            msync(data_, size_, MS_SYNC);
        }
    }
    
private:
    void* data_;
    size_t size_;
    bool readOnly_;
    int fd_;
};

} // namespace trading

#endif // MEMORY_MAPPED_FILE_H
```

### Cache-Friendly Data Structures

Optimizing data structures for CPU cache efficiency is crucial for performance:

#### Cache Line Alignment

```cpp
// cache_aligned.h
#ifndef CACHE_ALIGNED_H
#define CACHE_ALIGNED_H

#include <cstddef>

namespace trading {

// Cache line size (64 bytes on most x86 processors)
constexpr size_t CACHE_LINE_SIZE = 64;

// Macro for cache line alignment
#define CACHE_ALIGNED alignas(CACHE_LINE_SIZE)

// Template for cache-aligned allocator
template<typename T>
class CacheAlignedAllocator {
public:
    using value_type = T;
    
    CacheAlignedAllocator() noexcept {}
    
    template<typename U>
    CacheAlignedAllocator(const CacheAlignedAllocator<U>&) noexcept {}
    
    T* allocate(size_t n) {
        if (n == 0) {
            return nullptr;
        }
        
        if (n > std::numeric_limits<size_t>::max() / sizeof(T)) {
            throw std::bad_alloc();
        }
        
        void* p = aligned_alloc(CACHE_LINE_SIZE, n * sizeof(T));
        if (p == nullptr) {
            throw std::bad_alloc();
        }
        
        return static_cast<T*>(p);
    }
    
    void deallocate(T* p, size_t) noexcept {
        free(p);
    }
};

template<typename T, typename U>
bool operator==(const CacheAlignedAllocator<T>&, const CacheAlignedAllocator<U>&) {
    return true;
}

template<typename T, typename U>
bool operator!=(const CacheAlignedAllocator<T>&, const CacheAlignedAllocator<U>&) {
    return false;
}

// Cache-aligned vector
template<typename T>
using CacheAlignedVector = std::vector<T, CacheAlignedAllocator<T>>;

} // namespace trading

#endif // CACHE_ALIGNED_H
```

#### Structure of Arrays vs. Array of Structures

```cpp
// soa_order_book.h
#ifndef SOA_ORDER_BOOK_H
#define SOA_ORDER_BOOK_H

#include <vector>
#include <string>
#include <unordered_map>
#include "cache_aligned.h"

namespace trading {

// Structure of Arrays (SoA) design for order book
class SoAOrderBook {
public:
    SoAOrderBook(const std::string& symbol) : symbol_(symbol) {}
    
    // Add order to the book
    uint64_t addOrder(double price, int quantity, bool isBuy) {
        uint64_t orderId = nextOrderId_++;
        
        // Add order to appropriate side
        if (isBuy) {
            buyOrderIds_.push_back(orderId);
            buyPrices_.push_back(price);
            buyQuantities_.push_back(quantity);
        } else {
            sellOrderIds_.push_back(orderId);
            sellPrices_.push_back(price);
            sellQuantities_.push_back(quantity);
        }
        
        // Store order index for quick lookup
        orderIndexMap_[orderId] = isBuy ? 
            buyOrderIds_.size() - 1 : sellOrderIds_.size() - 1;
        orderSideMap_[orderId] = isBuy;
        
        return orderId;
    }
    
    // Cancel order
    bool cancelOrder(uint64_t orderId) {
        auto indexIt = orderIndexMap_.find(orderId);
        if (indexIt == orderIndexMap_.end()) {
            return false;
        }
        
        size_t index = indexIt->second;
        bool isBuy = orderSideMap_[orderId];
        
        if (isBuy) {
            // Move last order to this position
            if (index < buyOrderIds_.size() - 1) {
                buyOrderIds_[index] = buyOrderIds_.back();
                buyPrices_[index] = buyPrices_.back();
                buyQuantities_[index] = buyQuantities_.back();
                
                // Update index map
                orderIndexMap_[buyOrderIds_[index]] = index;
            }
            
            // Remove last element
            buyOrderIds_.pop_back();
            buyPrices_.pop_back();
            buyQuantities_.pop_back();
        } else {
            // Move last order to this position
            if (index < sellOrderIds_.size() - 1) {
                sellOrderIds_[index] = sellOrderIds_.back();
                sellPrices_[index] = sellPrices_.back();
                sellQuantities_[index] = sellQuantities_.back();
                
                // Update index map
                orderIndexMap_[sellOrderIds_[index]] = index;
            }
            
            // Remove last element
            sellOrderIds_.pop_back();
            sellPrices_.pop_back();
            sellQuantities_.pop_back();
        }
        
        // Remove from maps
        orderIndexMap_.erase(orderId);
        orderSideMap_.erase(orderId);
        
        return true;
    }
    
private:
    std::string symbol_;
    uint64_t nextOrderId_ = 1;
    
    // Buy side (Structure of Arrays)
    CACHE_ALIGNED std::vector<uint64_t> buyOrderIds_;
    CACHE_ALIGNED std::vector<double> buyPrices_;
    CACHE_ALIGNED std::vector<int> buyQuantities_;
    
    // Sell side (Structure of Arrays)
    CACHE_ALIGNED std::vector<uint64_t> sellOrderIds_;
    CACHE_ALIGNED std::vector<double> sellPrices_;
    CACHE_ALIGNED std::vector<int> sellQuantities_;
    
    // Maps for quick lookup
    std::unordered_map<uint64_t, size_t> orderIndexMap_;
    std::unordered_map<uint64_t, bool> orderSideMap_;
};

} // namespace trading

#endif // SOA_ORDER_BOOK_H
```

#### Flat Data Structures

```cpp
// flat_map.h
#ifndef FLAT_MAP_H
#define FLAT_MAP_H

#include <vector>
#include <algorithm>
#include <utility>
#include <functional>

namespace trading {

template<typename Key, typename Value, typename Compare = std::less<Key>>
class FlatMap {
public:
    using value_type = std::pair<Key, Value>;
    using iterator = typename std::vector<value_type>::iterator;
    using const_iterator = typename std::vector<value_type>::const_iterator;
    
    FlatMap() = default;
    
    // Find element by key
    iterator find(const Key& key) {
        auto it = std::lower_bound(data_.begin(), data_.end(), key,
            [](const value_type& p, const Key& k) {
                return Compare()(p.first, k);
            });
        
        if (it != data_.end() && !Compare()(key, it->first) && !Compare()(it->first, key)) {
            return it;
        }
        
        return data_.end();
    }
    
    const_iterator find(const Key& key) const {
        auto it = std::lower_bound(data_.begin(), data_.end(), key,
            [](const value_type& p, const Key& k) {
                return Compare()(p.first, k);
            });
        
        if (it != data_.end() && !Compare()(key, it->first) && !Compare()(it->first, key)) {
            return it;
        }
        
        return data_.end();
    }
    
    // Insert or update element
    std::pair<iterator, bool> insert(const value_type& value) {
        auto it = std::lower_bound(data_.begin(), data_.end(), value.first,
            [](const value_type& p, const Key& k) {
                return Compare()(p.first, k);
            });
        
        if (it != data_.end() && !Compare()(value.first, it->first) && !Compare()(it->first, value.first)) {
            // Key already exists, update value
            it->second = value.second;
            return {it, false};
        }
        
        // Insert new element
        it = data_.insert(it, value);
        return {it, true};
    }
    
    // Erase element by key
    size_t erase(const Key& key) {
        auto it = find(key);
        if (it == data_.end()) {
            return 0;
        }
        
        data_.erase(it);
        return 1;
    }
    
    // Access element by key (creates if not exists)
    Value& operator[](const Key& key) {
        auto it = find(key);
        if (it != data_.end()) {
            return it->second;
        }
        
        // Insert new element
        auto result = insert({key, Value()});
        return result.first->second;
    }
    
    // Clear all elements
    void clear() {
        data_.clear();
    }
    
    // Check if empty
    bool empty() const {
        return data_.empty();
    }
    
    // Get size
    size_t size() const {
        return data_.size();
    }
    
    // Iterators
    iterator begin() { return data_.begin(); }
    iterator end() { return data_.end(); }
    const_iterator begin() const { return data_.begin(); }
    const_iterator end() const { return data_.end(); }
    
private:
    std::vector<value_type> data_;
};

} // namespace trading

#endif // FLAT_MAP_H
```

## CPU Optimization

### SIMD Instructions

Single Instruction Multiple Data (SIMD) instructions enable parallel processing of data:

```cpp
// simd_utils.h
#ifndef SIMD_UTILS_H
#define SIMD_UTILS_H

#include <immintrin.h>
#include <vector>

namespace trading {

// SIMD-optimized price comparison
inline bool compareOrderBookLevels(const double* prices1, const double* prices2, size_t count) {
    // Process 4 doubles at a time using AVX
    size_t i = 0;
    for (; i + 3 < count; i += 4) {
        __m256d a = _mm256_loadu_pd(prices1 + i);
        __m256d b = _mm256_loadu_pd(prices2 + i);
        __m256d cmp = _mm256_cmp_pd(a, b, _CMP_NEQ_OQ);
        
        int mask = _mm256_movemask_pd(cmp);
        if (mask != 0) {
            return false;
        }
    }
    
    // Process remaining elements
    for (; i < count; ++i) {
        if (prices1[i] != prices2[i]) {
            return false;
        }
    }
    
    return true;
}

// SIMD-optimized volume calculation
inline double calculateTotalVolume(const double* prices, const int* quantities, size_t count) {
    __m256d sum = _mm256_setzero_pd();
    
    // Process 4 elements at a time
    size_t i = 0;
    for (; i + 3 < count; i += 4) {
        __m256d price = _mm256_loadu_pd(prices + i);
        
        // Load quantities and convert to double
        __m128i q = _mm_loadu_si128(reinterpret_cast<const __m128i*>(quantities + i));
        __m256d quantity = _mm256_cvtepi32_pd(q);
        
        // Multiply price by quantity and add to sum
        __m256d volume = _mm256_mul_pd(price, quantity);
        sum = _mm256_add_pd(sum, volume);
    }
    
    // Horizontal sum of the 4 doubles in sum
    __m128d sum128 = _mm_add_pd(_mm256_extractf128_pd(sum, 0), _mm256_extractf128_pd(sum, 1));
    __m128d sum64 = _mm_add_sd(sum128, _mm_unpackhi_pd(sum128, sum128));
    double result = _mm_cvtsd_f64(sum64);
    
    // Process remaining elements
    for (; i < count; ++i) {
        result += prices[i] * quantities[i];
    }
    
    return result;
}

} // namespace trading

#endif // SIMD_UTILS_H
```

### Branch Prediction Optimization

Optimizing code for branch prediction improves performance:

```cpp
// branch_prediction.h
#ifndef BRANCH_PREDICTION_H
#define BRANCH_PREDICTION_H

namespace trading {

// Macros for branch prediction hints
#define LIKELY(x) __builtin_expect(!!(x), 1)
#define UNLIKELY(x) __builtin_expect(!!(x), 0)

// Example usage in order matching
inline bool matchOrders(double buyPrice, double sellPrice, bool isMarketOrder) {
    // Market orders always match
    if (UNLIKELY(isMarketOrder)) {
        return true;
    }
    
    // For limit orders, buy price must be >= sell price
    return buyPrice >= sellPrice;
}

// Optimized order validation
inline bool validateOrder(const Order& order) {
    // Most orders have positive quantity
    if (LIKELY(order.quantity > 0)) {
        // Most orders have valid price
        if (LIKELY(order.price > 0.0 || order.type == OrderType::Market)) {
            // Most orders have valid symbol
            if (LIKELY(!order.symbol.empty())) {
                return true;
            }
        }
    }
    
    return false;
}

} // namespace trading

#endif // BRANCH_PREDICTION_H
```

### Loop Optimization

Optimizing loops for performance:

```cpp
// loop_optimization.h
#ifndef LOOP_OPTIMIZATION_H
#define LOOP_OPTIMIZATION_H

#include <vector>

namespace trading {

// Loop unrolling example
template<typename T>
inline void sumArray(const T* array, size_t size, T& result) {
    result = 0;
    
    // Process 4 elements at a time
    size_t i = 0;
    T sum1 = 0, sum2 = 0, sum3 = 0, sum4 = 0;
    
    for (; i + 3 < size; i += 4) {
        sum1 += array[i];
        sum2 += array[i + 1];
        sum3 += array[i + 2];
        sum4 += array[i + 3];
    }
    
    // Process remaining elements
    for (; i < size; ++i) {
        result += array[i];
    }
    
    result += sum1 + sum2 + sum3 + sum4;
}

// Loop tiling example
template<typename T>
inline void matrixMultiply(const T* A, const T* B, T* C, int n) {
    constexpr int BLOCK_SIZE = 32;
    
    // Zero the result matrix
    for (int i = 0; i < n * n; ++i) {
        C[i] = 0;
    }
    
    // Loop tiling
    for (int i0 = 0; i0 < n; i0 += BLOCK_SIZE) {
        for (int j0 = 0; j0 < n; j0 += BLOCK_SIZE) {
            for (int k0 = 0; k0 < n; k0 += BLOCK_SIZE) {
                // Process block
                for (int i = i0; i < std::min(i0 + BLOCK_SIZE, n); ++i) {
                    for (int j = j0; j < std::min(j0 + BLOCK_SIZE, n); ++j) {
                        T sum = 0;
                        for (int k = k0; k < std::min(k0 + BLOCK_SIZE, n); ++k) {
                            sum += A[i * n + k] * B[k * n + j];
                        }
                        C[i * n + j] += sum;
                    }
                }
            }
        }
    }
}

} // namespace trading

#endif // LOOP_OPTIMIZATION_H
```

## Concurrency Optimization

### Lock-Free Data Structures

Lock-free data structures improve performance in multi-threaded environments:

```cpp
// lock_free_queue.h
#ifndef LOCK_FREE_QUEUE_H
#define LOCK_FREE_QUEUE_H

#include <atomic>
#include <memory>

namespace trading {

template<typename T>
class LockFreeQueue {
private:
    struct Node {
        std::shared_ptr<T> data;
        std::atomic<Node*> next;
        
        Node() : next(nullptr) {}
    };
    
    std::atomic<Node*> head_;
    std::atomic<Node*> tail_;
    
public:
    LockFreeQueue() {
        Node* dummy = new Node();
        head_.store(dummy);
        tail_.store(dummy);
    }
    
    ~LockFreeQueue() {
        while (pop() != nullptr) {}
        
        Node* dummy = head_.load();
        delete dummy;
    }
    
    void push(T value) {
        std::shared_ptr<T> newData = std::make_shared<T>(std::move(value));
        Node* newNode = new Node();
        
        Node* oldTail = tail_.load();
        while (true) {
            Node* next = nullptr;
            if (oldTail->next.compare_exchange_weak(next, newNode)) {
                oldTail->data = newData;
                tail_.store(newNode);
                return;
            }
            
            oldTail = tail_.load();
        }
    }
    
    std::shared_ptr<T> pop() {
        Node* oldHead = head_.load();
        
        while (oldHead != tail_.load()) {
            Node* next = oldHead->next.load();
            if (head_.compare_exchange_weak(oldHead, next)) {
                std::shared_ptr<T> result = oldHead->data;
                delete oldHead;
                return result;
            }
            
            oldHead = head_.load();
        }
        
        return nullptr;
    }
};

} // namespace trading

#endif // LOCK_FREE_QUEUE_H
```

### Thread Pool

Efficient thread management with a thread pool:

```cpp
// thread_pool.h
#ifndef THREAD_POOL_H
#define THREAD_POOL_H

#include <vector>
#include <queue>
#include <memory>
#include <thread>
#include <mutex>
#include <condition_variable>
#include <future>
#include <functional>
#include <stdexcept>

namespace trading {

class ThreadPool {
public:
    ThreadPool(size_t numThreads) : stop(false) {
        for (size_t i = 0; i < numThreads; ++i) {
            workers.emplace_back([this] {
                while (true) {
                    std::function<void()> task;
                    
                    {
                        std::unique_lock<std::mutex> lock(queueMutex);
                        condition.wait(lock, [this] { 
                            return stop || !tasks.empty(); 
                        });
                        
                        if (stop && tasks.empty()) {
                            return;
                        }
                        
                        task = std::move(tasks.front());
                        tasks.pop();
                    }
                    
                    task();
                }
            });
        }
    }
    
    template<class F, class... Args>
    auto enqueue(F&& f, Args&&... args) 
        -> std::future<typename std::result_of<F(Args...)>::type> {
        using return_type = typename std::result_of<F(Args...)>::type;
        
        auto task = std::make_shared<std::packaged_task<return_type()>>(
            std::bind(std::forward<F>(f), std::forward<Args>(args)...)
        );
        
        std::future<return_type> result = task->get_future();
        
        {
            std::unique_lock<std::mutex> lock(queueMutex);
            
            if (stop) {
                throw std::runtime_error("enqueue on stopped ThreadPool");
            }
            
            tasks.emplace([task]() { (*task)(); });
        }
        
        condition.notify_one();
        return result;
    }
    
    ~ThreadPool() {
        {
            std::unique_lock<std::mutex> lock(queueMutex);
            stop = true;
        }
        
        condition.notify_all();
        
        for (std::thread& worker : workers) {
            worker.join();
        }
    }
    
private:
    std::vector<std::thread> workers;
    std::queue<std::function<void()>> tasks;
    
    std::mutex queueMutex;
    std::condition_variable condition;
    bool stop;
};

} // namespace trading

#endif // THREAD_POOL_H
```

### Work Stealing

Efficient task distribution with work stealing:

```cpp
// work_stealing_queue.h
#ifndef WORK_STEALING_QUEUE_H
#define WORK_STEALING_QUEUE_H

#include <atomic>
#include <vector>
#include <memory>
#include <deque>
#include <mutex>
#include <thread>
#include <functional>
#include <condition_variable>

namespace trading {

class WorkStealingQueue {
public:
    using Task = std::function<void()>;
    
    WorkStealingQueue() {}
    
    void push(Task task) {
        std::lock_guard<std::mutex> lock(mutex_);
        queue_.push_front(std::move(task));
    }
    
    bool pop(Task& task) {
        std::lock_guard<std::mutex> lock(mutex_);
        if (queue_.empty()) {
            return false;
        }
        
        task = std::move(queue_.front());
        queue_.pop_front();
        return true;
    }
    
    bool steal(Task& task) {
        std::lock_guard<std::mutex> lock(mutex_);
        if (queue_.empty()) {
            return false;
        }
        
        task = std::move(queue_.back());
        queue_.pop_back();
        return true;
    }
    
    bool empty() const {
        std::lock_guard<std::mutex> lock(mutex_);
        return queue_.empty();
    }
    
private:
    std::deque<Task> queue_;
    mutable std::mutex mutex_;
};

class WorkStealingThreadPool {
public:
    using Task = std::function<void()>;
    
    WorkStealingThreadPool(size_t numThreads) : stop_(false) {
        queues_.resize(numThreads);
        for (size_t i = 0; i < numThreads; ++i) {
            queues_[i] = std::make_unique<WorkStealingQueue>();
        }
        
        for (size_t i = 0; i < numThreads; ++i) {
            threads_.emplace_back([this, i] { workerThread(i); });
        }
    }
    
    ~WorkStealingThreadPool() {
        {
            std::unique_lock<std::mutex> lock(mutex_);
            stop_ = true;
        }
        
        condition_.notify_all();
        
        for (std::thread& thread : threads_) {
            thread.join();
        }
    }
    
    template<typename F, typename... Args>
    void enqueue(F&& f, Args&&... args) {
        auto task = std::bind(std::forward<F>(f), std::forward<Args>(args)...);
        
        {
            std::unique_lock<std::mutex> lock(mutex_);
            
            // Round-robin task distribution
            size_t idx = nextQueueIndex_++;
            if (nextQueueIndex_ >= queues_.size()) {
                nextQueueIndex_ = 0;
            }
            
            queues_[idx]->push(task);
        }
        
        condition_.notify_one();
    }
    
private:
    void workerThread(size_t index) {
        while (true) {
            Task task;
            bool hasTask = false;
            
            // Try to get task from own queue
            hasTask = queues_[index]->pop(task);
            
            // If no task, try to steal from other queues
            if (!hasTask) {
                for (size_t i = 0; i < queues_.size(); ++i) {
                    if (i == index) {
                        continue;
                    }
                    
                    if (queues_[i]->steal(task)) {
                        hasTask = true;
                        break;
                    }
                }
            }
            
            // If still no task, wait for notification
            if (!hasTask) {
                std::unique_lock<std::mutex> lock(mutex_);
                
                // Check if should stop
                if (stop_ && allQueuesEmpty()) {
                    return;
                }
                
                condition_.wait(lock, [this, index] {
                    return stop_ || !queues_[index]->empty() || !allQueuesEmpty();
                });
                
                // Check again if should stop
                if (stop_ && allQueuesEmpty()) {
                    return;
                }
                
                continue;
            }
            
            // Execute task
            task();
        }
    }
    
    bool allQueuesEmpty() const {
        for (const auto& queue : queues_) {
            if (!queue->empty()) {
                return false;
            }
        }
        return true;
    }
    
    std::vector<std::unique_ptr<WorkStealingQueue>> queues_;
    std::vector<std::thread> threads_;
    std::mutex mutex_;
    std::condition_variable condition_;
    bool stop_;
    size_t nextQueueIndex_ = 0;
};

} // namespace trading

#endif // WORK_STEALING_QUEUE_H
```

## I/O Optimization

### Zero-Copy Data Transfer

Minimizing data copying improves performance:

```cpp
// zero_copy.h
#ifndef ZERO_COPY_H
#define ZERO_COPY_H

#include <cstddef>
#include <cstring>
#include <memory>
#include <vector>

namespace trading {

// Zero-copy buffer
class ZeroCopyBuffer {
public:
    ZeroCopyBuffer(size_t capacity) : capacity_(capacity), size_(0) {
        data_ = std::make_unique<char[]>(capacity);
    }
    
    // Get pointer to data
    char* data() {
        return data_.get();
    }
    
    const char* data() const {
        return data_.get();
    }
    
    // Get size
    size_t size() const {
        return size_;
    }
    
    // Get capacity
    size_t capacity() const {
        return capacity_;
    }
    
    // Set size
    void setSize(size_t size) {
        if (size > capacity_) {
            throw std::runtime_error("Size exceeds capacity");
        }
        size_ = size;
    }
    
    // Resize buffer
    void resize(size_t newCapacity) {
        if (newCapacity < size_) {
            size_ = newCapacity;
        }
        
        auto newData = std::make_unique<char[]>(newCapacity);
        std::memcpy(newData.get(), data_.get(), size_);
        
        data_ = std::move(newData);
        capacity_ = newCapacity;
    }
    
private:
    std::unique_ptr<char[]> data_;
    size_t capacity_;
    size_t size_;
};

// Zero-copy message
class ZeroCopyMessage {
public:
    ZeroCopyMessage() : buffer_(nullptr), size_(0), owned_(false) {}
    
    // Create message with owned buffer
    ZeroCopyMessage(size_t size) : size_(size), owned_(true) {
        buffer_ = new char[size];
    }
    
    // Create message with external buffer
    ZeroCopyMessage(char* buffer, size_t size) : buffer_(buffer), size_(size), owned_(false) {}
    
    // Move constructor
    ZeroCopyMessage(ZeroCopyMessage&& other) noexcept 
        : buffer_(other.buffer_), size_(other.size_), owned_(other.owned_) {
        other.buffer_ = nullptr;
        other.size_ = 0;
        other.owned_ = false;
    }
    
    // Move assignment
    ZeroCopyMessage& operator=(ZeroCopyMessage&& other) noexcept {
        if (this != &other) {
            if (owned_ && buffer_ != nullptr) {
                delete[] buffer_;
            }
            
            buffer_ = other.buffer_;
            size_ = other.size_;
            owned_ = other.owned_;
            
            other.buffer_ = nullptr;
            other.size_ = 0;
            other.owned_ = false;
        }
        return *this;
    }
    
    // Destructor
    ~ZeroCopyMessage() {
        if (owned_ && buffer_ != nullptr) {
            delete[] buffer_;
        }
    }
    
    // Get buffer
    char* buffer() {
        return buffer_;
    }
    
    const char* buffer() const {
        return buffer_;
    }
    
    // Get size
    size_t size() const {
        return size_;
    }
    
    // Take ownership of buffer
    void takeOwnership() {
        owned_ = true;
    }
    
    // Release ownership of buffer
    char* releaseBuffer() {
        owned_ = false;
        return buffer_;
    }
    
private:
    char* buffer_;
    size_t size_;
    bool owned_;
    
    // Disable copy
    ZeroCopyMessage(const ZeroCopyMessage&) = delete;
    ZeroCopyMessage& operator=(const ZeroCopyMessage&) = delete;
};

} // namespace trading

#endif // ZERO_COPY_H
```

### Memory-Mapped I/O

Efficient I/O with memory mapping:

```cpp
// memory_mapped_io.h
#ifndef MEMORY_MAPPED_IO_H
#define MEMORY_MAPPED_IO_H

#include <string>
#include <vector>
#include <stdexcept>
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

namespace trading {

// Memory-mapped file reader
class MemoryMappedFileReader {
public:
    MemoryMappedFileReader(const std::string& filename) 
        : fd_(-1), data_(nullptr), size_(0) {
        // Open file
        fd_ = open(filename.c_str(), O_RDONLY);
        if (fd_ == -1) {
            throw std::runtime_error("Failed to open file: " + filename);
        }
        
        // Get file size
        struct stat sb;
        if (fstat(fd_, &sb) == -1) {
            close(fd_);
            throw std::runtime_error("Failed to get file size");
        }
        
        size_ = sb.st_size;
        
        // Map file into memory
        data_ = mmap(nullptr, size_, PROT_READ, MAP_PRIVATE, fd_, 0);
        
        if (data_ == MAP_FAILED) {
            close(fd_);
            throw std::runtime_error("Failed to map file into memory");
        }
    }
    
    ~MemoryMappedFileReader() {
        if (data_ != nullptr && data_ != MAP_FAILED) {
            munmap(data_, size_);
        }
        
        if (fd_ != -1) {
            close(fd_);
        }
    }
    
    // Get pointer to mapped memory
    const void* data() const {
        return data_;
    }
    
    // Get size of mapped memory
    size_t size() const {
        return size_;
    }
    
private:
    int fd_;
    void* data_;
    size_t size_;
};

// Memory-mapped file writer
class MemoryMappedFileWriter {
public:
    MemoryMappedFileWriter(const std::string& filename, size_t size) 
        : fd_(-1), data_(nullptr), size_(size) {
        // Open file
        fd_ = open(filename.c_str(), O_RDWR | O_CREAT | O_TRUNC, 0644);
        if (fd_ == -1) {
            throw std::runtime_error("Failed to open file: " + filename);
        }
        
        // Set file size
        if (ftruncate(fd_, size) == -1) {
            close(fd_);
            throw std::runtime_error("Failed to set file size");
        }
        
        // Map file into memory
        data_ = mmap(nullptr, size, PROT_READ | PROT_WRITE, MAP_SHARED, fd_, 0);
        
        if (data_ == MAP_FAILED) {
            close(fd_);
            throw std::runtime_error("Failed to map file into memory");
        }
    }
    
    ~MemoryMappedFileWriter() {
        if (data_ != nullptr && data_ != MAP_FAILED) {
            msync(data_, size_, MS_SYNC);
            munmap(data_, size_);
        }
        
        if (fd_ != -1) {
            close(fd_);
        }
    }
    
    // Get pointer to mapped memory
    void* data() {
        return data_;
    }
    
    // Get size of mapped memory
    size_t size() const {
        return size_;
    }
    
    // Flush changes to disk
    void flush() {
        if (data_ != nullptr && data_ != MAP_FAILED) {
            msync(data_, size_, MS_SYNC);
        }
    }
    
private:
    int fd_;
    void* data_;
    size_t size_;
};

} // namespace trading

#endif // MEMORY_MAPPED_IO_H
```

## Profiling and Benchmarking

### Performance Measurement

Tools for measuring performance:

```cpp
// performance_measurement.h
#ifndef PERFORMANCE_MEASUREMENT_H
#define PERFORMANCE_MEASUREMENT_H

#include <chrono>
#include <string>
#include <vector>
#include <algorithm>
#include <numeric>
#include <iostream>
#include <iomanip>
#include <atomic>

namespace trading {

// High-resolution timer
class Timer {
public:
    Timer() : start_(std::chrono::high_resolution_clock::now()) {}
    
    // Reset timer
    void reset() {
        start_ = std::chrono::high_resolution_clock::now();
    }
    
    // Get elapsed time in nanoseconds
    int64_t elapsedNanos() const {
        auto now = std::chrono::high_resolution_clock::now();
        return std::chrono::duration_cast<std::chrono::nanoseconds>(now - start_).count();
    }
    
    // Get elapsed time in microseconds
    double elapsedMicros() const {
        return elapsedNanos() / 1000.0;
    }
    
    // Get elapsed time in milliseconds
    double elapsedMillis() const {
        return elapsedNanos() / 1000000.0;
    }
    
    // Get elapsed time in seconds
    double elapsedSeconds() const {
        return elapsedNanos() / 1000000000.0;
    }
    
private:
    std::chrono::high_resolution_clock::time_point start_;
};

// Scoped timer for measuring function execution time
class ScopedTimer {
public:
    ScopedTimer(const std::string& name) : name_(name), timer_() {}
    
    ~ScopedTimer() {
        std::cout << name_ << ": " << timer_.elapsedMicros() << " us" << std::endl;
    }
    
private:
    std::string name_;
    Timer timer_;
};

// Performance statistics
class PerformanceStats {
public:
    PerformanceStats() : count_(0), min_(std::numeric_limits<double>::max()), max_(0), sum_(0) {}
    
    // Add sample
    void addSample(double value) {
        samples_.push_back(value);
        sum_ += value;
        count_++;
        
        min_ = std::min(min_, value);
        max_ = std::max(max_, value);
    }
    
    // Get count
    size_t count() const {
        return count_;
    }
    
    // Get minimum
    double min() const {
        return min_;
    }
    
    // Get maximum
    double max() const {
        return max_;
    }
    
    // Get average
    double avg() const {
        return count_ > 0 ? sum_ / count_ : 0;
    }
    
    // Get percentile
    double percentile(double p) const {
        if (samples_.empty()) {
            return 0;
        }
        
        std::vector<double> sorted = samples_;
        std::sort(sorted.begin(), sorted.end());
        
        double index = p * (sorted.size() - 1);
        size_t lowerIndex = static_cast<size_t>(index);
        size_t upperIndex = std::min(lowerIndex + 1, sorted.size() - 1);
        double weight = index - lowerIndex;
        
        return sorted[lowerIndex] * (1 - weight) + sorted[upperIndex] * weight;
    }
    
    // Get standard deviation
    double stddev() const {
        if (count_ <= 1) {
            return 0;
        }
        
        double mean = avg();
        double sumSquaredDiff = 0;
        
        for (double sample : samples_) {
            double diff = sample - mean;
            sumSquaredDiff += diff * diff;
        }
        
        return std::sqrt(sumSquaredDiff / (count_ - 1));
    }
    
    // Print statistics
    void print(const std::string& name) const {
        std::cout << "=== " << name << " ===" << std::endl;
        std::cout << "Count: " << count_ << std::endl;
        std::cout << "Min: " << min_ << " us" << std::endl;
        std::cout << "Max: " << max_ << " us" << std::endl;
        std::cout << "Avg: " << avg() << " us" << std::endl;
        std::cout << "StdDev: " << stddev() << " us" << std::endl;
        std::cout << "P50: " << percentile(0.5) << " us" << std::endl;
        std::cout << "P90: " << percentile(0.9) << " us" << std::endl;
        std::cout << "P99: " << percentile(0.99) << " us" << std::endl;
        std::cout << "P99.9: " << percentile(0.999) << " us" << std::endl;
    }
    
private:
    std::vector<double> samples_;
    size_t count_;
    double min_;
    double max_;
    double sum_;
};

// Thread-safe performance counter
class PerformanceCounter {
public:
    PerformanceCounter() : count_(0), sum_(0) {}
    
    // Increment counter
    void increment() {
        count_.fetch_add(1, std::memory_order_relaxed);
    }
    
    // Add value
    void add(double value) {
        count_.fetch_add(1, std::memory_order_relaxed);
        sum_.fetch_add(value, std::memory_order_relaxed);
    }
    
    // Get count
    uint64_t count() const {
        return count_.load(std::memory_order_relaxed);
    }
    
    // Get sum
    double sum() const {
        return sum_.load(std::memory_order_relaxed);
    }
    
    // Get average
    double avg() const {
        uint64_t c = count();
        return c > 0 ? sum() / c : 0;
    }
    
    // Reset counter
    void reset() {
        count_.store(0, std::memory_order_relaxed);
        sum_.store(0, std::memory_order_relaxed);
    }
    
private:
    std::atomic<uint64_t> count_;
    std::atomic<double> sum_;
};

} // namespace trading

#endif // PERFORMANCE_MEASUREMENT_H
```

### Benchmarking Framework

Framework for benchmarking components:

```cpp
// benchmark_framework.h
#ifndef BENCHMARK_FRAMEWORK_H
#define BENCHMARK_FRAMEWORK_H

#include <functional>
#include <string>
#include <vector>
#include <chrono>
#include <iostream>
#include <iomanip>
#include <algorithm>
#include <numeric>
#include <random>
#include <thread>
#include <atomic>
#include <mutex>
#include <condition_variable>
#include "performance_measurement.h"

namespace trading {

// Benchmark configuration
struct BenchmarkConfig {
    size_t numThreads = 1;
    size_t numIterations = 1000000;
    size_t warmupIterations = 100000;
    bool measureLatency = true;
    bool printProgress = true;
    size_t progressInterval = 100000;
};

// Benchmark result
struct BenchmarkResult {
    std::string name;
    double throughput;  // operations per second
    PerformanceStats latencyStats;
};

// Benchmark framework
class BenchmarkFramework {
public:
    // Run benchmark
    static BenchmarkResult runBenchmark(
        const std::string& name,
        std::function<void()> setup,
        std::function<void()> teardown,
        std::function<void()> benchmark,
        const BenchmarkConfig& config = BenchmarkConfig()) {
        
        BenchmarkResult result;
        result.name = name;
        
        // Setup
        setup();
        
        // Warmup
        for (size_t i = 0; i < config.warmupIterations; ++i) {
            benchmark();
        }
        
        // Prepare threads
        std::vector<std::thread> threads;
        std::atomic<size_t> completedIterations(0);
        std::atomic<bool> startFlag(false);
        std::atomic<bool> stopFlag(false);
        std::mutex mutex;
        std::condition_variable cv;
        
        std::vector<PerformanceStats> threadStats(config.numThreads);
        
        // Create threads
        for (size_t t = 0; t < config.numThreads; ++t) {
            threads.emplace_back([&, t]() {
                // Wait for start signal
                {
                    std::unique_lock<std::mutex> lock(mutex);
                    cv.wait(lock, [&]() { return startFlag.load(); });
                }
                
                // Run benchmark
                size_t localIterations = config.numIterations / config.numThreads;
                if (t == 0) {
                    localIterations += config.numIterations % config.numThreads;
                }
                
                for (size_t i = 0; i < localIterations && !stopFlag.load(); ++i) {
                    Timer timer;
                    
                    benchmark();
                    
                    if (config.measureLatency) {
                        double latency = timer.elapsedMicros();
                        threadStats[t].addSample(latency);
                    }
                    
                    size_t completed = completedIterations.fetch_add(1) + 1;
                    
                    if (config.printProgress && completed % config.progressInterval == 0) {
                        std::cout << "Progress: " << completed << " / " << config.numIterations 
                                  << " (" << (completed * 100.0 / config.numIterations) << "%)" << std::endl;
                    }
                }
            });
        }
        
        // Start benchmark
        Timer timer;
        {
            std::lock_guard<std::mutex> lock(mutex);
            startFlag.store(true);
        }
        cv.notify_all();
        
        // Wait for completion or timeout
        const int maxDurationSeconds = 60;
        for (int i = 0; i < maxDurationSeconds * 10 && completedIterations.load() < config.numIterations; ++i) {
            std::this_thread::sleep_for(std::chrono::milliseconds(100));
        }
        
        // Stop benchmark
        stopFlag.store(true);
        
        // Wait for threads to finish
        for (auto& thread : threads) {
            thread.join();
        }
        
        double elapsed = timer.elapsedSeconds();
        size_t completed = completedIterations.load();
        
        // Calculate throughput
        result.throughput = completed / elapsed;
        
        // Merge latency stats
        if (config.measureLatency) {
            for (const auto& stats : threadStats) {
                for (size_t i = 0; i < stats.count(); ++i) {
                    result.latencyStats.addSample(stats.percentile(static_cast<double>(i) / stats.count()));
                }
            }
        }
        
        // Print results
        std::cout << "=== Benchmark Results: " << name << " ===" << std::endl;
        std::cout << "Threads: " << config.numThreads << std::endl;
        std::cout << "Iterations: " << completed << std::endl;
        std::cout << "Elapsed Time: " << elapsed << " seconds" << std::endl;
        std::cout << "Throughput: " << result.throughput << " ops/sec" << std::endl;
        
        if (config.measureLatency) {
            result.latencyStats.print("Latency (us)");
        }
        
        // Teardown
        teardown();
        
        return result;
    }
};

} // namespace trading

#endif // BENCHMARK_FRAMEWORK_H
```

## Compiler Optimization

### Compiler Flags

Optimizing compiler flags for performance:

```cmake
# compiler_optimization.cmake

# Set C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Set optimization flags for different build types
set(CMAKE_CXX_FLAGS_DEBUG "-g -O0 -fno-omit-frame-pointer")
set(CMAKE_CXX_FLAGS_RELEASE "-O3 -DNDEBUG -march=native -flto -fno-exceptions")
set(CMAKE_CXX_FLAGS_RELWITHDEBINFO "-O2 -g -DNDEBUG -march=native -fno-omit-frame-pointer")
set(CMAKE_CXX_FLAGS_MINSIZEREL "-Os -DNDEBUG -march=native -flto -fno-exceptions")

# Set common flags
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra -Wpedantic -Werror")

# Enable link-time optimization for release builds
set(CMAKE_INTERPROCEDURAL_OPTIMIZATION_RELEASE ON)
set(CMAKE_INTERPROCEDURAL_OPTIMIZATION_MINSIZEREL ON)

# Set platform-specific flags
if(CMAKE_CXX_COMPILER_ID MATCHES "GNU")
    # GCC-specific flags
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fno-strict-aliasing -fno-builtin-memcmp")
    set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} -ffunction-sections -fdata-sections")
    set(CMAKE_EXE_LINKER_FLAGS_RELEASE "${CMAKE_EXE_LINKER_FLAGS_RELEASE} -Wl,--gc-sections")
elseif(CMAKE_CXX_COMPILER_ID MATCHES "Clang")
    # Clang-specific flags
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fno-strict-aliasing")
    set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} -ffunction-sections -fdata-sections")
    set(CMAKE_EXE_LINKER_FLAGS_RELEASE "${CMAKE_EXE_LINKER_FLAGS_RELEASE} -Wl,--gc-sections")
endif()

# Enable sanitizers for debug builds
option(ENABLE_ASAN "Enable Address Sanitizer" OFF)
option(ENABLE_TSAN "Enable Thread Sanitizer" OFF)
option(ENABLE_UBSAN "Enable Undefined Behavior Sanitizer" OFF)

if(ENABLE_ASAN)
    set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -fsanitize=address -fno-omit-frame-pointer")
    set(CMAKE_EXE_LINKER_FLAGS_DEBUG "${CMAKE_EXE_LINKER_FLAGS_DEBUG} -fsanitize=address")
endif()

if(ENABLE_TSAN)
    set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -fsanitize=thread -fno-omit-frame-pointer")
    set(CMAKE_EXE_LINKER_FLAGS_DEBUG "${CMAKE_EXE_LINKER_FLAGS_DEBUG} -fsanitize=thread")
endif()

if(ENABLE_UBSAN)
    set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -fsanitize=undefined -fno-omit-frame-pointer")
    set(CMAKE_EXE_LINKER_FLAGS_DEBUG "${CMAKE_EXE_LINKER_FLAGS_DEBUG} -fsanitize=undefined")
endif()

# Enable profile-guided optimization
option(ENABLE_PGO "Enable Profile-Guided Optimization" OFF)
option(PGO_GENERATE "Generate profile data" OFF)
option(PGO_USE "Use profile data" OFF)

if(ENABLE_PGO)
    if(PGO_GENERATE)
        set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fprofile-generate")
        set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -fprofile-generate")
    elseif(PGO_USE)
        set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fprofile-use -fprofile-correction")
        set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -fprofile-use")
    endif()
endif()
```

### Inline Assembly

Using inline assembly for critical sections:

```cpp
// inline_assembly.h
#ifndef INLINE_ASSEMBLY_H
#define INLINE_ASSEMBLY_H

#include <cstdint>

namespace trading {

// Fast integer division by constant
inline int fastDivideBy10(int value) {
    int result;
    
    #if defined(__x86_64__) || defined(_M_X64)
    __asm__ (
        "movl %1, %%eax\n"
        "movl $1717986919, %%edx\n"
        "imull %%edx\n"
        "sarl $2, %%edx\n"
        "movl %%edx, %0\n"
        : "=r" (result)
        : "r" (value)
        : "eax", "edx"
    );
    #else
    result = value / 10;
    #endif
    
    return result;
}

// Fast population count (count number of set bits)
inline int popCount(uint64_t value) {
    int result;
    
    #if defined(__x86_64__) || defined(_M_X64)
    __asm__ (
        "popcnt %1, %0\n"
        : "=r" (result)
        : "r" (value)
    );
    #else
    result = __builtin_popcountll(value);
    #endif
    
    return result;
}

// Fast bit scan forward (find index of first set bit)
inline int bitScanForward(uint64_t value) {
    int result;
    
    #if defined(__x86_64__) || defined(_M_X64)
    __asm__ (
        "bsfq %1, %0\n"
        : "=r" (result)
        : "r" (value)
    );
    #else
    result = __builtin_ctzll(value);
    #endif
    
    return result;
}

// Fast bit scan reverse (find index of last set bit)
inline int bitScanReverse(uint64_t value) {
    int result;
    
    #if defined(__x86_64__) || defined(_M_X64)
    __asm__ (
        "bsrq %1, %0\n"
        : "=r" (result)
        : "r" (value)
    );
    #else
    result = 63 - __builtin_clzll(value);
    #endif
    
    return result;
}

} // namespace trading

#endif // INLINE_ASSEMBLY_H
```

## Performance Tuning

### Runtime Configuration

Configuring the system for optimal performance:

```cpp
// performance_tuning.h
#ifndef PERFORMANCE_TUNING_H
#define PERFORMANCE_TUNING_H

#include <thread>
#include <vector>
#include <string>
#include <fstream>
#include <sstream>
#include <iostream>
#include <cstdlib>
#include <cstring>
#include <unistd.h>
#include <sys/resource.h>
#include <pthread.h>

namespace trading {

class PerformanceTuner {
public:
    // Set thread affinity
    static bool setThreadAffinity(std::thread& thread, int cpuId) {
        cpu_set_t cpuset;
        CPU_ZERO(&cpuset);
        CPU_SET(cpuId, &cpuset);
        
        int rc = pthread_setaffinity_np(thread.native_handle(), sizeof(cpu_set_t), &cpuset);
        return rc == 0;
    }
    
    // Set thread priority
    static bool setThreadPriority(std::thread& thread, int priority) {
        sched_param param;
        param.sched_priority = priority;
        
        int rc = pthread_setschedparam(thread.native_handle(), SCHED_FIFO, &param);
        return rc == 0;
    }
    
    // Set process priority
    static bool setProcessPriority(int priority) {
        return setpriority(PRIO_PROCESS, 0, priority) == 0;
    }
    
    // Lock memory to prevent swapping
    static bool lockMemory() {
        return mlockall(MCL_CURRENT | MCL_FUTURE) == 0;
    }
    
    // Set CPU governor to performance mode
    static bool setCpuGovernorToPerformance() {
        int numCpus = std::thread::hardware_concurrency();
        bool success = true;
        
        for (int i = 0; i < numCpus; ++i) {
            std::string path = "/sys/devices/system/cpu/cpu" + std::to_string(i) + "/cpufreq/scaling_governor";
            std::ofstream file(path);
            
            if (file.is_open()) {
                file << "performance";
                file.close();
            } else {
                success = false;
            }
        }
        
        return success;
    }
    
    // Disable CPU power saving
    static bool disableCpuPowerSaving() {
        std::system("echo 1 > /sys/devices/system/cpu/intel_pstate/no_turbo");
        std::system("echo 0 > /sys/devices/system/cpu/cpufreq/boost");
        
        return true;
    }
    
    // Set file descriptor limit
    static bool setFileDescriptorLimit(int limit) {
        struct rlimit rlim;
        rlim.rlim_cur = limit;
        rlim.rlim_max = limit;
        
        return setrlimit(RLIMIT_NOFILE, &rlim) == 0;
    }
    
    // Configure huge pages
    static bool configureHugePages(int numPages) {
        std::string command = "echo " + std::to_string(numPages) + " > /proc/sys/vm/nr_hugepages";
        return std::system(command.c_str()) == 0;
    }
    
    // Disable NUMA balancing
    static bool disableNumaBalancing() {
        return std::system("echo 0 > /proc/sys/kernel/numa_balancing") == 0;
    }
    
    // Set network parameters
    static bool setNetworkParameters() {
        std::system("echo 1 > /proc/sys/net/ipv4/tcp_low_latency");
        std::system("echo 1 > /proc/sys/net/ipv4/tcp_nodelay");
        std::system("echo 1 > /proc/sys/net/ipv4/tcp_fastopen");
        
        return true;
    }
    
    // Apply all performance tuning settings
    static bool applyAllSettings() {
        bool success = true;
        
        // Check if running as root
        if (geteuid() != 0) {
            std::cerr << "Warning: Some performance tuning settings require root privileges" << std::endl;
        }
        
        // Apply settings
        success &= lockMemory();
        success &= setCpuGovernorToPerformance();
        success &= disableCpuPowerSaving();
        success &= setFileDescriptorLimit(1000000);
        success &= configureHugePages(1024);
        success &= disableNumaBalancing();
        success &= setNetworkParameters();
        
        return success;
    }
};

} // namespace trading

#endif // PERFORMANCE_TUNING_H
```

### Adaptive Performance

Adapting to system conditions:

```cpp
// adaptive_performance.h
#ifndef ADAPTIVE_PERFORMANCE_H
#define ADAPTIVE_PERFORMANCE_H

#include <thread>
#include <chrono>
#include <atomic>
#include <vector>
#include <algorithm>
#include <functional>
#include <mutex>
#include <condition_variable>

namespace trading {

class AdaptivePerformanceMonitor {
public:
    AdaptivePerformanceMonitor(
        std::function<void(int)> setThreadCountCallback,
        std::function<void(bool)> setHighPerformanceModeCallback)
        : running_(false), threadCount_(std::thread::hardware_concurrency()),
          highPerformanceMode_(false), loadThresholdHigh_(0.8), loadThresholdLow_(0.2),
          setThreadCountCallback_(setThreadCountCallback),
          setHighPerformanceModeCallback_(setHighPerformanceModeCallback) {}
    
    // Start monitoring
    void start() {
        if (running_.exchange(true)) {
            return;
        }
        
        monitorThread_ = std::thread([this]() {
            while (running_.load()) {
                // Get current system load
                double load = getSystemLoad();
                
                // Adjust thread count based on load
                int optimalThreadCount = calculateOptimalThreadCount(load);
                if (optimalThreadCount != threadCount_) {
                    threadCount_ = optimalThreadCount;
                    setThreadCountCallback_(threadCount_);
                }
                
                // Adjust performance mode based on load
                bool newHighPerformanceMode = load > loadThresholdHigh_;
                if (newHighPerformanceMode != highPerformanceMode_) {
                    highPerformanceMode_ = newHighPerformanceMode;
                    setHighPerformanceModeCallback_(highPerformanceMode_);
                }
                
                // Sleep for monitoring interval
                std::this_thread::sleep_for(std::chrono::seconds(1));
            }
        });
    }
    
    // Stop monitoring
    void stop() {
        if (!running_.exchange(false)) {
            return;
        }
        
        if (monitorThread_.joinable()) {
            monitorThread_.join();
        }
    }
    
    // Set load thresholds
    void setLoadThresholds(double low, double high) {
        loadThresholdLow_ = low;
        loadThresholdHigh_ = high;
    }
    
private:
    // Get current system load
    double getSystemLoad() {
        std::ifstream loadFile("/proc/loadavg");
        double load = 0.0;
        
        if (loadFile.is_open()) {
            loadFile >> load;
        }
        
        return load / std::thread::hardware_concurrency();
    }
    
    // Calculate optimal thread count based on load
    int calculateOptimalThreadCount(double load) {
        int maxThreads = std::thread::hardware_concurrency();
        
        if (load > loadThresholdHigh_) {
            // High load, use all available threads
            return maxThreads;
        } else if (load < loadThresholdLow_) {
            // Low load, use minimum threads
            return std::max(1, maxThreads / 4);
        } else {
            // Medium load, scale linearly
            double loadRange = loadThresholdHigh_ - loadThresholdLow_;
            double loadRatio = (load - loadThresholdLow_) / loadRange;
            int threadRange = maxThreads - maxThreads / 4;
            
            return std::max(1, maxThreads / 4 + static_cast<int>(loadRatio * threadRange));
        }
    }
    
    std::atomic<bool> running_;
    int threadCount_;
    bool highPerformanceMode_;
    double loadThresholdLow_;
    double loadThresholdHigh_;
    std::thread monitorThread_;
    std::function<void(int)> setThreadCountCallback_;
    std::function<void(bool)> setHighPerformanceModeCallback_;
};

class AdaptiveExecutionEngine {
public:
    AdaptiveExecutionEngine() 
        : threadCount_(std::thread::hardware_concurrency()),
          highPerformanceMode_(false),
          monitor_([this](int count) { setThreadCount(count); },
                  [this](bool mode) { setHighPerformanceMode(mode); }) {
        // Initialize thread pool
        workers_.resize(threadCount_);
        for (int i = 0; i < threadCount_; ++i) {
            startWorker(i);
        }
        
        // Start performance monitor
        monitor_.start();
    }
    
    ~AdaptiveExecutionEngine() {
        // Stop performance monitor
        monitor_.stop();
        
        // Stop all workers
        {
            std::unique_lock<std::mutex> lock(mutex_);
            stop_ = true;
        }
        condition_.notify_all();
        
        for (auto& worker : workers_) {
            if (worker.joinable()) {
                worker.join();
            }
        }
    }
    
    // Enqueue task
    template<typename F, typename... Args>
    void enqueue(F&& f, Args&&... args) {
        auto task = std::bind(std::forward<F>(f), std::forward<Args>(args)...);
        
        {
            std::unique_lock<std::mutex> lock(mutex_);
            tasks_.push_back(std::move(task));
        }
        
        condition_.notify_one();
    }
    
private:
    // Set thread count
    void setThreadCount(int count) {
        std::unique_lock<std::mutex> lock(mutex_);
        
        if (count == threadCount_) {
            return;
        }
        
        if (count > threadCount_) {
            // Add workers
            int oldCount = threadCount_;
            threadCount_ = count;
            
            workers_.resize(count);
            for (int i = oldCount; i < count; ++i) {
                startWorker(i);
            }
        } else {
            // Remove workers
            int oldCount = threadCount_;
            threadCount_ = count;
            
            // Signal workers to exit
            condition_.notify_all();
            
            // Wait for excess workers to exit
            for (int i = count; i < oldCount; ++i) {
                if (workers_[i].joinable()) {
                    workers_[i].join();
                }
            }
            
            workers_.resize(count);
        }
    }
    
    // Set high performance mode
    void setHighPerformanceMode(bool mode) {
        std::unique_lock<std::mutex> lock(mutex_);
        
        if (mode == highPerformanceMode_) {
            return;
        }
        
        highPerformanceMode_ = mode;
        
        if (highPerformanceMode_) {
            // Enable optimizations for high performance
            // ...
        } else {
            // Disable optimizations to save resources
            // ...
        }
    }
    
    // Start worker thread
    void startWorker(int index) {
        workers_[index] = std::thread([this, index]() {
            while (true) {
                std::function<void()> task;
                
                {
                    std::unique_lock<std::mutex> lock(mutex_);
                    
                    condition_.wait(lock, [this, index]() {
                        return stop_ || !tasks_.empty() || index >= threadCount_;
                    });
                    
                    if (stop_ || index >= threadCount_) {
                        return;
                    }
                    
                    if (!tasks_.empty()) {
                        task = std::move(tasks_.front());
                        tasks_.pop_front();
                    }
                }
                
                if (task) {
                    task();
                }
            }
        });
    }
    
    int threadCount_;
    bool highPerformanceMode_;
    std::vector<std::thread> workers_;
    std::deque<std::function<void()>> tasks_;
    std::mutex mutex_;
    std::condition_variable condition_;
    bool stop_ = false;
    AdaptivePerformanceMonitor monitor_;
};

} // namespace trading

#endif // ADAPTIVE_PERFORMANCE_H
```

## Conclusion

This document provides a comprehensive guide to performance optimization techniques for the C++ order execution engine in the trading platform. By applying these techniques, the engine can achieve the performance targets required for high-frequency trading scenarios.

The key optimization areas covered include:

1. **Memory Optimization**: Custom allocators, object recycling, cache-friendly data structures, and memory-mapped files.
2. **CPU Optimization**: SIMD instructions, branch prediction optimization, and loop optimization.
3. **Concurrency Optimization**: Lock-free data structures, thread pools, and work stealing.
4. **I/O Optimization**: Zero-copy data transfer and memory-mapped I/O.
5. **Profiling and Benchmarking**: Performance measurement and benchmarking framework.
6. **Compiler Optimization**: Compiler flags and inline assembly.
7. **Performance Tuning**: Runtime configuration and adaptive performance.

By implementing these techniques, the C++ order execution engine can provide the ultra-low latency and high throughput required for modern trading systems, while maintaining reliability and correctness.
