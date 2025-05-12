-- Create reviews table if it doesn't exist
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for product_id
CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews(product_id);

-- Create index for user_id
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);

-- Create index for created_at (for sorting by newest)
CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at DESC);