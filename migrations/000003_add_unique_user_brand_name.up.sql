CREATE UNIQUE INDEX ux_sellers_user_brand_name_active
ON sellers (
        user_id,
        lower(btrim(brand_name))
    )
    WHERE deleted_at IS NULL;