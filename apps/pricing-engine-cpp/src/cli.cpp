#include "pricing.h"

#include <iostream>
#include <vector>

int main(int argc, char** argv) {
    std::vector<pricing::LineItem> items{
        {"SKU-001", 1999, 2},
        {"SKU-002", 4500, 1},
    };
    std::string coupon = argc > 1 ? argv[1] : "";
    auto t = pricing::compute(items, coupon, 0.0875, 599);
    std::cout << "subtotal=" << t.subtotal_cents
              << " discount=" << t.discount_cents
              << " tax=" << t.tax_cents
              << " shipping=" << t.shipping_cents
              << " total=" << t.total_cents << "\n";
    return 0;
}
