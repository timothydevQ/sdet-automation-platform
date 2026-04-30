#include "pricing.h"

#include <cmath>

namespace pricing {

int64_t compute_tax(int64_t taxable_cents, double rate) {
    if (rate <= 0.0) return 0;
    double raw = static_cast<double>(taxable_cents) * rate;
    return static_cast<int64_t>(std::llround(raw));
}

}
