-- 学术资料已由 call_report_materials 结构化关联，不再使用旧自由文本约束。
ALTER TABLE call_reports
    DROP CONSTRAINT IF EXISTS chk_call_reports_material_note;

-- 防止绕过应用层写入状态与关键时间字段不一致的数据。
ALTER TABLE visits
    ADD CONSTRAINT chk_visits_checked_out_after_checked_in
        CHECK (checked_out_at IS NULL OR (checked_in_at IS NOT NULL AND checked_out_at >= checked_in_at)),
    ADD CONSTRAINT chk_visits_status_required_timestamps
        CHECK (
            (status = 'PLANNED' AND checked_in_at IS NULL AND checked_out_at IS NULL)
            OR (status = 'CHECKED_IN' AND checked_in_at IS NOT NULL AND checked_out_at IS NULL)
            OR (status IN ('CHECKED_OUT', 'COMPLETED') AND checked_in_at IS NOT NULL AND checked_out_at IS NOT NULL)
            OR status = 'CANCELLED'
        );

-- 月度 Dashboard 以签退时间做范围扫描，并只统计已产生最终合规结论的记录。
CREATE INDEX IF NOT EXISTS idx_visits_checked_out_final_compliance
    ON visits(checked_out_at, product_id)
    WHERE status IN ('CHECKED_OUT', 'COMPLETED')
      AND compliance_status IN ('COMPLIANT', 'NON_COMPLIANT');
