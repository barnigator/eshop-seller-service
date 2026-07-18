CREATE TABLE social_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID NOT NULL,
    type TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_social_links_seller
        FOREIGN KEY (seller_id)
        REFERENCES sellers(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_social_link_type
        CHECK (
            type IN (
                     'telegram',
                     'vk',
                     'instagram',
                     'website',
                     'whatsapp',
                     'youtube',
                     'facebook',
                     'ozon',
                     'wildberries',
                     'yandexmarket',
                     'avito'
            )
        )

);

CREATE INDEX idx_social_links_seller_id
    ON social_links(seller_id);