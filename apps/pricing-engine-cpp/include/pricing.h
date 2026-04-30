#pragma once

#include <cstdint>
#include <string>
#include <vector>

namespace pricing {

struct LineItem {
    std::string sku;
    int64_t unit_price_cents;
    int qty;
};

struct CartTotals {
    int64_t subtotal_cents;
    int64_t discount_cents;
    int64_t tax_cents;
    int64_t shipping_cents;
    int64_t total_cents;
};

CartTotals compute(const std::vector<LineItem>& items,
                   const std::string& coupon,
                   double tax_rate,
                   int64_t shipping_cents);

int64_t apply_coupon(int64_t subtotal_cents, const std::string& coupon);

int64_t compute_tax(int64_t taxable_cents, double rate);

bool valid_coupon(const std::string& coupon);

}
