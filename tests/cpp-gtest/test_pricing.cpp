#include "pricing.h"

#include <gtest/gtest.h>

using pricing::LineItem;

TEST(Pricing, EmptyCart) {
    auto t = pricing::compute({}, "", 0.0875, 0);
    EXPECT_EQ(t.subtotal_cents, 0);
    EXPECT_EQ(t.total_cents, 0);
}

TEST(Pricing, BasicSubtotal) {
    std::vector<LineItem> items{{"SKU-001", 1999, 2}, {"SKU-002", 1000, 1}};
    auto t = pricing::compute(items, "", 0.0, 0);
    EXPECT_EQ(t.subtotal_cents, 4998);
    EXPECT_EQ(t.discount_cents, 0);
    EXPECT_EQ(t.tax_cents, 0);
}

TEST(Pricing, TaxApplied) {
    std::vector<LineItem> items{{"SKU-001", 10000, 1}};
    auto t = pricing::compute(items, "", 0.10, 0);
    EXPECT_EQ(t.tax_cents, 1000);
    EXPECT_EQ(t.total_cents, 11000);
}

TEST(Pricing, ShippingAdded) {
    std::vector<LineItem> items{{"SKU-001", 5000, 1}};
    auto t = pricing::compute(items, "", 0.0, 599);
    EXPECT_EQ(t.shipping_cents, 599);
    EXPECT_EQ(t.total_cents, 5599);
}

TEST(Pricing, NegativePriceThrows) {
    std::vector<LineItem> items{{"SKU-X", -100, 1}};
    EXPECT_THROW(pricing::compute(items, "", 0.0, 0), std::invalid_argument);
}

TEST(Pricing, ZeroQtyThrows) {
    std::vector<LineItem> items{{"SKU-X", 100, 0}};
    EXPECT_THROW(pricing::compute(items, "", 0.0, 0), std::invalid_argument);
}
