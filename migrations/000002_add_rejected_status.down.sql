ALTER TABLE sellers
DROP CONSTRAINT sellers_status_valid;

ALTER TABLE sellers
ADD CONSTRAINT sellers_status_valid
CHECK (
    status IN (
        'pending',
        'active',
        'blocked,',
        'archived'
        )
    );