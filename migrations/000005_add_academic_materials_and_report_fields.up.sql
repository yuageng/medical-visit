ALTER TABLE call_reports
    ADD COLUMN IF NOT EXISTS additional_notes TEXT;

-- 学术资料改由结构化关联表维护，不再要求旧的自由文本说明。
ALTER TABLE call_reports
    DROP CONSTRAINT IF EXISTS chk_call_reports_material_note;

CREATE TABLE IF NOT EXISTS academic_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    product_id UUID,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_academic_materials_product FOREIGN KEY (product_id)
        REFERENCES products(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_academic_materials_product_id ON academic_materials(product_id);
CREATE INDEX IF NOT EXISTS idx_academic_materials_active ON academic_materials(active);

CREATE TABLE IF NOT EXISTS call_report_materials (
    call_report_id UUID NOT NULL,
    material_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (call_report_id, material_id),
    CONSTRAINT fk_call_report_materials_report FOREIGN KEY (call_report_id)
        REFERENCES call_reports(id) ON DELETE CASCADE,
    CONSTRAINT fk_call_report_materials_material FOREIGN KEY (material_id)
        REFERENCES academic_materials(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_call_report_materials_material_id
    ON call_report_materials(material_id);
