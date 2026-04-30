cmake_minimum_required(VERSION 3.16)
project(pricing_engine CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_library(pricing_engine STATIC
    src/pricing.cpp
    src/coupons.cpp
    src/tax.cpp
)
target_include_directories(pricing_engine PUBLIC include)

add_executable(pricing_cli src/cli.cpp)
target_link_libraries(pricing_cli PRIVATE pricing_engine)

option(BUILD_TESTING "Build pricing tests" OFF)
if(BUILD_TESTING)
    include(FetchContent)
    FetchContent_Declare(
        googletest
        URL https://github.com/google/googletest/archive/refs/tags/v1.14.0.zip
    )
    FetchContent_MakeAvailable(googletest)
    enable_testing()
    add_executable(pricing_tests
        ../../tests/cpp-gtest/test_pricing.cpp
        ../../tests/cpp-gtest/test_coupons.cpp
        ../../tests/cpp-gtest/test_tax.cpp
    )
    target_link_libraries(pricing_tests PRIVATE pricing_engine GTest::gtest_main)
    include(GoogleTest)
    gtest_discover_tests(pricing_tests)
endif()
