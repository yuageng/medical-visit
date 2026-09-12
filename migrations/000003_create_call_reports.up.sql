-- 创建拜访记录表
CREATE TABLE IF NOT EXISTS call_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_id UUID NOT NULL UNIQUE,
    discussion_summary TEXT NOT NULL,
    hcp_feedback TEXT,
    materials_distributed BOOLEAN NOT NULL DEFAULT FALSE,
    material_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_call_reports_visit FOREIGN KEY (visit_id) 
        REFERENCES visits(id) ON DELETE RESTRICT,
    CONSTRAINT chk_call_reports_material_note CHECK (
        (materials_distributed = FALSE) OR 
        (materials_distributed = TRUE AND material_note IS NOT NULL AND LENGTH(TRIM(material_note)) > 0)
    )
);

CREATE UNIQUE INDEX uq_call_reports_visit_id ON call_reports(visit_id);

COMMENT ON TABLE call_reports IS '拜访后学术沟通记录表';
COMMENT ON COLUMN call_reports.materials_distributed IS '是否派发合规学术资料';
COMMENT ON COLUMN call_reports.material_note IS '资料说明，派发时必填';
