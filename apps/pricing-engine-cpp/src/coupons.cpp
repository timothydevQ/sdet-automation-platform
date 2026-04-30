#include "pricing.h"

#include <algorithm>

namespace pricing {

static std::string upper(std::string s) {
    std::transform(s.begin(), s.end(), s.begin(), ::toupper);
    return s;
}

bool valid_coupon(const std::string& coupon) {
    if (coupon.empty()) return true;
    std::string c = upper(coupon);
    return c == "WELCOME10" || c == "BULK20" || c == "FLAT5" || c == "VIP25";
}

int64_t apply_coupon(int64_t subtotal_cents, const std::string& coupon) {
    if (coupon.empty() || !valid_coupon(coupon)) return subtotal_cents;
    std::string c = upper(coupon);
    if (c == "WELCOME10") {
        return subtotal_cents - subtotal_cents / 10;
    }
    if (c == "BULK20") {
        if (subtotal_cents < 5000) return subtotal_cents;
        return subtotal_cents - (subtotal_cents * 20 + 50) / 100;
    }
    if (c == "FLAT5") {
        return std::max<int64_t>(0, subtotal_cents - 500);
    }
    if (c == "VIP25") {
        return subtotal_cents - (subtotal_cents * 25) / 100;
    }
    return subtotal_cents;
}

}
