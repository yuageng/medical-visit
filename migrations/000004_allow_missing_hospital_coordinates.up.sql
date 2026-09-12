ALTER TABLE hospitals
    ALTER COLUMN latitude DROP NOT NULL,
    ALTER COLUMN longitude DROP NOT NULL;

ALTER TABLE hospitals
    DROP CONSTRAINT IF EXISTS chk_hospitals_latitude,
    DROP CONSTRAINT IF EXISTS chk_hospitals_longitude;

ALTER TABLE hospitals
    ADD CONSTRAINT chk_hospitals_latitude CHECK (latitude IS NULL OR (latitude >= -90 AND latitude <= 90)),
    ADD CONSTRAINT chk_hospitals_longitude CHECK (longitude IS NULL OR (longitude >= -180 AND longitude <= 180)),
    ADD CONSTRAINT chk_hospitals_coordinates_complete CHECK ((latitude IS NULL) = (longitude IS NULL));
