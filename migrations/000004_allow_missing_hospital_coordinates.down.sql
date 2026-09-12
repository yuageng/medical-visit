ALTER TABLE hospitals
    DROP CONSTRAINT IF EXISTS chk_hospitals_coordinates_complete;

ALTER TABLE hospitals
    DROP CONSTRAINT IF EXISTS chk_hospitals_latitude,
    DROP CONSTRAINT IF EXISTS chk_hospitals_longitude;

ALTER TABLE hospitals
    ADD CONSTRAINT chk_hospitals_latitude CHECK (latitude >= -90 AND latitude <= 90),
    ADD CONSTRAINT chk_hospitals_longitude CHECK (longitude >= -180 AND longitude <= 180);

ALTER TABLE hospitals
    ALTER COLUMN latitude SET NOT NULL,
    ALTER COLUMN longitude SET NOT NULL;
