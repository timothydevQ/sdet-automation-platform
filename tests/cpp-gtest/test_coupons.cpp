#include "pricing.h"

#include <gtest/gtest.h>

TEST(Coupons, NoCoupon) {
    EXPECT_EQ(pricing::apply_coupon(10000, ""), 10000);
}

TEST(Coupons, Welcome10) {
    EXPECT_EQ(pricing::apply_coupon(10000, "WELCOME10"), 9000);
}

TEST(Coupons, Bulk20OnlyAboveThreshold) {
    EXPECT_EQ(pricing::apply_coupon(4999, "BULK20"), 4999);
    EXPECT_EQ(pricing::apply_coupon(10000, "BULK20"), 8000);
}

TEST(Coupons, Flat5DoesNotGoNegative) {
    EXPECT_EQ(pricing::apply_coupon(300, "FLAT5"), 0);
    EXPECT_EQ(pricing::apply_coupon(1000, "FLAT5"), 500);
}

TEST(Coupons, VIP25) {
    EXPECT_EQ(pricing::apply_coupon(10000, "VIP25"), 7500);
}

TEST(Coupons, UnknownCouponIgnored) {
    EXPECT_EQ(pricing::apply_coupon(10000, "BOGUS"), 10000);
}

TEST(Coupons, CouponIsCaseInsensitive) {
    EXPECT_EQ(pricing::apply_coupon(10000, "welcome10"), 9000);
    EXPECT_EQ(pricing::apply_coupon(10000, "Welcome10"), 9000);
}
