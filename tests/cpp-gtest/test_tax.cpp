#include "pricing.h"

#include <gtest/gtest.h>

TEST(Tax, ZeroRate) {
    EXPECT_EQ(pricing::compute_tax(10000, 0.0), 0);
}

TEST(Tax, BasicRate) {
    EXPECT_EQ(pricing::compute_tax(10000, 0.0875), 875);
}

TEST(Tax, RoundsHalfUp) {
    EXPECT_EQ(pricing::compute_tax(123, 0.05), 6);
}

TEST(Tax, NegativeRateTreatedAsZero) {
    EXPECT_EQ(pricing::compute_tax(10000, -0.5), 0);
}
