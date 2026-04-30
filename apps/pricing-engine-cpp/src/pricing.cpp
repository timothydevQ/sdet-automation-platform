#include "pricing.h"

#include <stdexcept>

namespace pricing {

CartTotals compute(const std::vector<LineItem>& items,
                   const std::string& coupon,
                   double tax_rate,
                   int64_t shipping_cents) {
    CartTotals t{};
    for (const auto& it : items) {
        if (it.unit_price_cents < 0) {
            throw std::invalid_argument("negative price");
        }
        if (it.qty <= 0) {
            throw std::invalid_argument("non-positive qty");
        }
        t.subtotal_cents += it.unit_price_cents * it.qty;
    }
    int64_t after_coupon = apply_coupon(t.subtotal_cents, coupon);
    t.discount_cents = t.subtotal_cents - after_coupon;
    t.tax_cents = compute_tax(after_coupon, tax_rate);
    t.shipping_cents = shipping_cents;
    t.total_cents = after_coupon + t.tax_cents + t.shipping_cents;
    return t;
}

}
