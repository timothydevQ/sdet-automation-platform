# pricing-engine-cpp

Small C++ library that computes cart totals: subtotal, discount, tax, shipping, total.

Build:
```
cmake -S . -B build
cmake --build build
./build/pricing_cli BULK20
```

Tests live in `tests/cpp-gtest/` and are built when `-DBUILD_TESTING=ON`.
