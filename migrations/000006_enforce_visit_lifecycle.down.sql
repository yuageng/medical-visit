DROP INDEX IF EXISTS idx_visits_checked_out_final_compliance;

ALTER TABLE visits
    DROP CONSTRAINT IF EXISTS chk_visits_status_required_timestamps,
    DROP CONSTRAINT IF EXISTS chk_visits_checked_out_after_checked_in;

ALTER TABLE call_reports
    ADD CONSTRAINT chk_call_reports_material_note CHECK (
        (materials_distributed = FALSE) OR
        (materials_distributed = TRUE AND material_note IS NOT NULL AND LENGTH(TRIM(material_note)) > 0)
    );
