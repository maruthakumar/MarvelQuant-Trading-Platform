# C++ Environment Setup

## Overview

This document provides comprehensive instructions for setting up the C++ development environment for the trading platform's high-performance order execution engine. The environment is designed to support modern C++ development with a focus on performance, reliability, and maintainability.

## System Requirements

### Hardware Recommendations

- **CPU**: Modern multi-core processor (8+ cores recommended)
- **Memory**: 16GB+ RAM
- **Storage**: SSD with at least 50GB free space
- **Network**: Low-latency network interface for production environments

### Software Requirements

- **Operating System**: Linux (Ubuntu 20.04 LTS or newer recommended)
- **Compiler**: GCC 10+ or Clang 12+
- **Build System**: CMake 3.15+
- **Version Control**: Git 2.25+
- **Package Manager**: Conan 1.40+ (optional but recommended)
- **Containerization**: Docker 20.10+ (for development and testing)

## Installation Guide

### Basic Development Environment

1. **Install Essential Tools**

   ```bash
   sudo apt update
   sudo apt install -y build-essential git cmake ninja-build ccache
   ```

2. **Install Modern Compiler**

   ```bash
   # For GCC
   sudo apt install -y gcc-10 g++-10
   sudo update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-10 100
   sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-10 100

   # For Clang
   sudo apt install -y clang-12 lld-12
   sudo update-alternatives --install /usr/bin/clang clang /usr/bin/clang-12 100
   sudo update-alternatives --install /usr/bin/clang++ clang++ /usr/bin/clang++-12 100
   ```

3. **Install Required Libraries**

   ```bash
   sudo apt install -y libboost-all-dev libtbb-dev libfmt-dev
   sudo apt install -y libgtest-dev libgmock-dev libunwind-dev
   sudo apt install -y libbenchmark-dev
   ```

4. **Install Conan Package Manager (Optional)**

   ```bash
   pip3 install conan
   conan profile new default --detect
   conan profile update settings.compiler.libcxx=libstdc++11 default
   ```

### Docker-Based Development Environment

For a consistent development environment across all team members, we provide a Docker-based setup:

1. **Install Docker**

   ```bash
   sudo apt update
   sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
   sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
   sudo apt update
   sudo apt install -y docker-ce docker-ce-cli containerd.io
   sudo usermod -aG docker $USER
   ```

2. **Build Development Container**

   Create a `Dockerfile` in the project root:

   ```dockerfile
   FROM ubuntu:20.04

   ENV DEBIAN_FRONTEND=noninteractive

   RUN apt-get update && apt-get install -y \
       build-essential \
       cmake \
       ninja-build \
       ccache \
       git \
       gcc-10 \
       g++-10 \
       clang-12 \
       lld-12 \
       libboost-all-dev \
       libtbb-dev \
       libfmt-dev \
       libgtest-dev \
       libgmock-dev \
       libunwind-dev \
       libbenchmark-dev \
       python3-pip \
       && rm -rf /var/lib/apt/lists/*

   RUN update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-10 100 \
       && update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-10 100

   RUN pip3 install conan
   RUN conan profile new default --detect \
       && conan profile update settings.compiler.libcxx=libstdc++11 default

   WORKDIR /app
   ```

3. **Build and Run Development Container**

   ```bash
   docker build -t trading-platform-cpp-dev .
   docker run -it --rm -v $(pwd):/app trading-platform-cpp-dev
   ```

## Project Setup

### Directory Structure

Create the following directory structure for the C++ components:

```
cpp/
├── include/                 # Public header files
│   ├── order_book/          # Order book data structures
│   ├── matching_engine/     # Matching engine components
│   ├── market_data/         # Market data handlers
│   └── common/              # Common utilities and interfaces
├── src/                     # Implementation files
│   ├── order_book/          # Order book implementation
│   ├── matching_engine/     # Matching engine implementation
│   ├── market_data/         # Market data handlers implementation
│   ├── memory/              # Memory management implementation
│   └── interface/           # Go-C++ interface implementation
├── tests/                   # Unit and integration tests
│   ├── unit/                # Unit tests for components
│   ├── integration/         # Integration tests
│   └── performance/         # Performance benchmarks
├── benchmarks/              # Performance benchmarking code
├── examples/                # Example usage of components
├── tools/                   # Development tools and scripts
└── cmake/                   # CMake modules and configuration
```

### CMake Configuration

1. **Create Root CMakeLists.txt**

   ```cmake
   cmake_minimum_required(VERSION 3.15)
   project(TradingPlatformCpp VERSION 0.1.0 LANGUAGES CXX)

   # Set C++ standard
   set(CMAKE_CXX_STANDARD 17)
   set(CMAKE_CXX_STANDARD_REQUIRED ON)
   set(CMAKE_CXX_EXTENSIONS OFF)

   # Set compiler flags
   if(CMAKE_CXX_COMPILER_ID MATCHES "GNU|Clang")
       set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra -Wpedantic -Werror")
       set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} -O3 -march=native")
       set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -g -O0 -fno-omit-frame-pointer")
   endif()

   # Enable testing
   enable_testing()

   # Find required packages
   find_package(Boost REQUIRED COMPONENTS system thread)
   find_package(TBB REQUIRED)
   find_package(fmt REQUIRED)
   find_package(GTest REQUIRED)
   find_package(benchmark REQUIRED)

   # Add subdirectories
   add_subdirectory(src)
   add_subdirectory(tests)
   add_subdirectory(benchmarks)
   add_subdirectory(examples)
   ```

2. **Create src/CMakeLists.txt**

   ```cmake
   # Add library targets
   add_library(order_book SHARED
       order_book/order_book.cpp
       order_book/limit_order.cpp
       order_book/market_order.cpp
   )
   target_include_directories(order_book PUBLIC ${CMAKE_SOURCE_DIR}/include)
   target_link_libraries(order_book PRIVATE Boost::boost TBB::tbb)

   add_library(matching_engine SHARED
       matching_engine/matching_engine.cpp
       matching_engine/price_time_priority.cpp
   )
   target_include_directories(matching_engine PUBLIC ${CMAKE_SOURCE_DIR}/include)
   target_link_libraries(matching_engine PRIVATE order_book Boost::boost TBB::tbb)

   add_library(market_data SHARED
       market_data/market_data_handler.cpp
       market_data/feed_normalizer.cpp
   )
   target_include_directories(market_data PUBLIC ${CMAKE_SOURCE_DIR}/include)
   target_link_libraries(market_data PRIVATE Boost::boost TBB::tbb fmt::fmt)

   add_library(memory_manager SHARED
       memory/memory_pool.cpp
       memory/object_recycler.cpp
   )
   target_include_directories(memory_manager PUBLIC ${CMAKE_SOURCE_DIR}/include)
   target_link_libraries(memory_manager PRIVATE Boost::boost)

   add_library(go_interface SHARED
       interface/go_interface.cpp
       interface/serialization.cpp
   )
   target_include_directories(go_interface PUBLIC ${CMAKE_SOURCE_DIR}/include)
   target_link_libraries(go_interface PRIVATE 
       order_book 
       matching_engine 
       market_data 
       memory_manager 
       Boost::boost 
       fmt::fmt
   )
   ```

3. **Create tests/CMakeLists.txt**

   ```cmake
   # Add test targets
   add_executable(order_book_tests
       unit/order_book/order_book_test.cpp
       unit/order_book/limit_order_test.cpp
       unit/order_book/market_order_test.cpp
   )
   target_link_libraries(order_book_tests PRIVATE order_book GTest::gtest GTest::gtest_main)
   add_test(NAME OrderBookTests COMMAND order_book_tests)

   add_executable(matching_engine_tests
       unit/matching_engine/matching_engine_test.cpp
       unit/matching_engine/price_time_priority_test.cpp
   )
   target_link_libraries(matching_engine_tests PRIVATE matching_engine GTest::gtest GTest::gtest_main)
   add_test(NAME MatchingEngineTests COMMAND matching_engine_tests)

   add_executable(market_data_tests
       unit/market_data/market_data_handler_test.cpp
       unit/market_data/feed_normalizer_test.cpp
   )
   target_link_libraries(market_data_tests PRIVATE market_data GTest::gtest GTest::gtest_main)
   add_test(NAME MarketDataTests COMMAND market_data_tests)

   add_executable(memory_manager_tests
       unit/memory/memory_pool_test.cpp
       unit/memory/object_recycler_test.cpp
   )
   target_link_libraries(memory_manager_tests PRIVATE memory_manager GTest::gtest GTest::gtest_main)
   add_test(NAME MemoryManagerTests COMMAND memory_manager_tests)

   add_executable(go_interface_tests
       unit/interface/go_interface_test.cpp
       unit/interface/serialization_test.cpp
   )
   target_link_libraries(go_interface_tests PRIVATE go_interface GTest::gtest GTest::gtest_main)
   add_test(NAME GoInterfaceTests COMMAND go_interface_tests)

   add_executable(integration_tests
       integration/end_to_end_test.cpp
   )
   target_link_libraries(integration_tests PRIVATE 
       order_book 
       matching_engine 
       market_data 
       memory_manager 
       go_interface 
       GTest::gtest 
       GTest::gtest_main
   )
   add_test(NAME IntegrationTests COMMAND integration_tests)
   ```

## Build Instructions

### Basic Build

1. **Create Build Directory**

   ```bash
   mkdir -p build && cd build
   ```

2. **Configure Project**

   ```bash
   cmake ..
   ```

3. **Build Project**

   ```bash
   cmake --build . -j$(nproc)
   ```

4. **Run Tests**

   ```bash
   ctest -V
   ```

### Build Configurations

1. **Debug Build**

   ```bash
   cmake -DCMAKE_BUILD_TYPE=Debug ..
   cmake --build . -j$(nproc)
   ```

2. **Release Build**

   ```bash
   cmake -DCMAKE_BUILD_TYPE=Release ..
   cmake --build . -j$(nproc)
   ```

3. **Build with Clang**

   ```bash
   CC=clang CXX=clang++ cmake ..
   cmake --build . -j$(nproc)
   ```

4. **Build with Address Sanitizer**

   ```bash
   cmake -DCMAKE_BUILD_TYPE=Debug -DENABLE_ASAN=ON ..
   cmake --build . -j$(nproc)
   ```

5. **Build with Thread Sanitizer**

   ```bash
   cmake -DCMAKE_BUILD_TYPE=Debug -DENABLE_TSAN=ON ..
   cmake --build . -j$(nproc)
   ```

## Development Tools

### Code Formatting

1. **Install clang-format**

   ```bash
   sudo apt install -y clang-format
   ```

2. **Create .clang-format Configuration**

   ```yaml
   BasedOnStyle: Google
   IndentWidth: 4
   ColumnLimit: 100
   AlignConsecutiveAssignments: true
   AlignConsecutiveDeclarations: true
   AllowShortFunctionsOnASingleLine: Empty
   AllowShortIfStatementsOnASingleLine: false
   AllowShortLoopsOnASingleLine: false
   BinPackArguments: false
   BinPackParameters: false
   BreakBeforeBraces: Allman
   BreakConstructorInitializersBeforeComma: true
   ConstructorInitializerAllOnOneLineOrOnePerLine: true
   PointerAlignment: Left
   SortIncludes: true
   ```

3. **Format Code**

   ```bash
   find include src tests -name "*.cpp" -o -name "*.h" | xargs clang-format -i
   ```

### Static Analysis

1. **Install clang-tidy**

   ```bash
   sudo apt install -y clang-tidy
   ```

2. **Create .clang-tidy Configuration**

   ```yaml
   Checks: 'clang-diagnostic-*,clang-analyzer-*,cppcoreguidelines-*,performance-*,readability-*,modernize-*,-modernize-use-trailing-return-type,-readability-magic-numbers'
   WarningsAsErrors: ''
   HeaderFilterRegex: '.*'
   FormatStyle: file
   ```

3. **Run Static Analysis**

   ```bash
   find include src -name "*.cpp" -o -name "*.h" | xargs clang-tidy -p build
   ```

### Memory Analysis

1. **Install Valgrind**

   ```bash
   sudo apt install -y valgrind
   ```

2. **Run Memory Check**

   ```bash
   valgrind --leak-check=full --show-leak-kinds=all ./build/tests/unit/order_book_tests
   ```

## Continuous Integration

### GitHub Actions Configuration

Create `.github/workflows/cpp.yml`:

```yaml
name: C++ CI

on:
  push:
    branches: [ main ]
    paths:
      - 'cpp/**'
      - '.github/workflows/cpp.yml'
  pull_request:
    branches: [ main ]
    paths:
      - 'cpp/**'
      - '.github/workflows/cpp.yml'

jobs:
  build:
    runs-on: ubuntu-20.04
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Install dependencies
      run: |
        sudo apt-get update
        sudo apt-get install -y build-essential cmake ninja-build ccache
        sudo apt-get install -y gcc-10 g++-10
        sudo apt-get install -y libboost-all-dev libtbb-dev libfmt-dev
        sudo apt-get install -y libgtest-dev libgmock-dev libunwind-dev libbenchmark-dev
        sudo update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-10 100
        sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-10 100
    
    - name: Configure CMake
      run: |
        cd cpp
        mkdir -p build
        cd build
        cmake ..
    
    - name: Build
      run: |
        cd cpp/build
        cmake --build . -j$(nproc)
    
    - name: Test
      run: |
        cd cpp/build
        ctest -V
    
    - name: Static Analysis
      run: |
        sudo apt-get install -y clang-tidy
        cd cpp
        find include src -name "*.cpp" -o -name "*.h" | xargs clang-tidy -p build
```

### Jenkins Pipeline Configuration

Create `Jenkinsfile`:

```groovy
pipeline {
    agent {
        docker {
            image 'trading-platform-cpp-dev:latest'
            args '-v ${WORKSPACE}:/app'
        }
    }
    
    stages {
        stage('Configure') {
            steps {
                sh '''
                cd cpp
                mkdir -p build
                cd build
                cmake ..
                '''
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                cd cpp/build
                cmake --build . -j$(nproc)
                '''
            }
        }
        
        stage('Test') {
            steps {
                sh '''
                cd cpp/build
                ctest -V
                '''
            }
        }
        
        stage('Benchmark') {
            steps {
                sh '''
                cd cpp/build
                ./benchmarks/order_book_benchmark
                ./benchmarks/matching_engine_benchmark
                '''
            }
        }
        
        stage('Static Analysis') {
            steps {
                sh '''
                cd cpp
                find include src -name "*.cpp" -o -name "*.h" | xargs clang-tidy -p build
                '''
            }
        }
        
        stage('Package') {
            steps {
                sh '''
                cd cpp/build
                cpack -G TGZ
                '''
            }
        }
    }
    
    post {
        always {
            archiveArtifacts artifacts: 'cpp/build/*.tar.gz', fingerprint: true
            junit 'cpp/build/test-results/*.xml'
        }
    }
}
```

## Performance Profiling

### CPU Profiling with perf

1. **Install perf**

   ```bash
   sudo apt install -y linux-tools-common linux-tools-generic
   ```

2. **Run Profiling**

   ```bash
   perf record -g ./build/benchmarks/matching_engine_benchmark
   perf report
   ```

### Memory Profiling with Heaptrack

1. **Install Heaptrack**

   ```bash
   sudo apt install -y heaptrack heaptrack-gui
   ```

2. **Run Profiling**

   ```bash
   heaptrack ./build/benchmarks/order_book_benchmark
   heaptrack_gui heaptrack.order_book_benchmark.*.gz
   ```

## Debugging

### GDB Configuration

1. **Install GDB**

   ```bash
   sudo apt install -y gdb
   ```

2. **Create .gdbinit**

   ```
   set print pretty on
   set print object on
   set print static-members on
   set print demangle on
   set demangle-style gnu-v3
   set print sevenbit-strings off
   ```

3. **Debug Application**

   ```bash
   gdb ./build/tests/unit/matching_engine_tests
   ```

### Core Dump Analysis

1. **Enable Core Dumps**

   ```bash
   ulimit -c unlimited
   ```

2. **Set Core Pattern**

   ```bash
   sudo sh -c 'echo "core.%e.%p.%t" > /proc/sys/kernel/core_pattern'
   ```

3. **Analyze Core Dump**

   ```bash
   gdb ./build/tests/unit/matching_engine_tests core.matching_engine_tests.12345.1612345678
   ```

## Conclusion

This document provides a comprehensive guide for setting up the C++ development environment for the trading platform's high-performance order execution engine. By following these instructions, developers can create a consistent, efficient, and productive environment for working on the C++ components of the system.

The environment is designed to support modern C++ development practices with a focus on performance, reliability, and maintainability. The tools and configurations provided enable thorough testing, profiling, and analysis to ensure the highest quality code.

For any issues or questions regarding the C++ environment setup, please contact the development team lead.
