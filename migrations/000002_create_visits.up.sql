-- 创建拜访表
CREATE TABLE IF NOT EXISTS visits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- 业务归属
    mr_id UUID NOT NULL,
    hcp_id UUID NOT NULL,
    hospital_id UUID NOT NULL,
    department_id UUID NOT NULL,
    product_id UUID NOT NULL,
    
    -- 计划信息
    planned_start_at TIMESTAMPTZ NOT NULL,
    plan_note TEXT,
    
    -- 状态
    status VARCHAR(30) NOT NULL DEFAULT 'PLANNED',
    
    -- 签到信息
    checked_in_at TIMESTAMPTZ,
    check_in_latitude NUMERIC(9,6),
    check_in_longitude NUMERIC(9,6),
    check_in_distance_m NUMERIC(10,2),
    
    -- 签退信息
    checked_out_at TIMESTAMPTZ,
    check_out_latitude NUMERIC(9,6),
    check_out_longitude NUMERIC(9,6),
    check_out_distance_m NUMERIC(10,2),
    
    -- 时长
    duration_seconds INTEGER,
    
    -- 合规结果
    compliance_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    anomaly_reasons JSONB NOT NULL DEFAULT '[]'::jsonb,
    
    -- 审计与乐观锁
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- 外键约束
    CONSTRAINT fk_visits_mr FOREIGN KEY (mr_id) 
        REFERENCES medical_representatives(id) ON DELETE RESTRICT,
    CONSTRAINT fk_visits_hcp FOREIGN KEY (hcp_id) 
        REFERENCES hcps(id) ON DELETE RESTRICT,
    CONSTRAINT fk_visits_hospital FOREIGN KEY (hospital_id) 
        REFERENCES hospitals(id) ON DELETE RESTRICT,
    CONSTRAINT fk_visits_department FOREIGN KEY (department_id) 
        REFERENCES departments(id) ON DELETE RESTRICT,
    CONSTRAINT fk_visits_product FOREIGN KEY (product_id) 
        REFERENCES products(id) ON DELETE RESTRICT,
    
    -- 状态约束
    CONSTRAINT chk_visits_status CHECK (
        status IN ('PLANNED', 'CHECKED_IN', 'CHECKED_OUT', 'COMPLETED', 'CANCELLED')
    ),
    CONSTRAINT chk_visits_compliance_status CHECK (
        compliance_status IN ('PENDING', 'COMPLIANT', 'NON_COMPLIANT')
    ),
    
    -- 经纬度范围约束
    CONSTRAINT chk_visits_check_in_latitude CHECK (
        check_in_latitude IS NULL OR (check_in_latitude >= -90 AND check_in_latitude <= 90)
    ),
    CONSTRAINT chk_visits_check_in_longitude CHECK (
        check_in_longitude IS NULL OR (check_in_longitude >= -180 AND check_in_longitude <= 180)
    ),
    CONSTRAINT chk_visits_check_out_latitude CHECK (
        check_out_latitude IS NULL OR (check_out_latitude >= -90 AND check_out_latitude <= 90)
    ),
    CONSTRAINT chk_visits_check_out_longitude CHECK (
        check_out_longitude IS NULL OR (check_out_longitude >= -180 AND check_out_longitude <= 180)
    ),
    
    -- 业务逻辑约束
    CONSTRAINT chk_visits_check_in_complete CHECK (
        (checked_in_at IS NULL AND check_in_latitude IS NULL AND check_in_longitude IS NULL AND check_in_distance_m IS NULL)
        OR (checked_in_at IS NOT NULL AND check_in_latitude IS NOT NULL AND check_in_longitude IS NOT NULL AND check_in_distance_m IS NOT NULL)
    ),
    CONSTRAINT chk_visits_check_out_complete CHECK (
        (checked_out_at IS NULL AND check_out_latitude IS NULL AND check_out_longitude IS NULL AND check_out_distance_m IS NULL)
        OR (checked_out_at IS NOT NULL AND check_out_latitude IS NOT NULL AND check_out_longitude IS NOT NULL AND check_out_distance_m IS NOT NULL)
    ),
    CONSTRAINT chk_visits_duration_non_negative CHECK (
        duration_seconds IS NULL OR duration_seconds >= 0
    )
);

-- 核心索引
CREATE INDEX idx_visits_mr_id ON visits(mr_id, planned_start_at DESC);
CREATE INDEX idx_visits_product_checked_out ON visits(product_id, checked_out_at DESC) 
    WHERE checked_out_at IS NOT NULL;
CREATE INDEX idx_visits_status ON visits(status, planned_start_at DESC);
CREATE INDEX idx_visits_hospital_id ON visits(hospital_id);
CREATE INDEX idx_visits_hcp_id ON visits(hcp_id);
CREATE INDEX idx_visits_compliance_status ON visits(compliance_status, checked_out_at DESC);
CREATE INDEX idx_visits_created_at ON visits(created_at DESC);

COMMENT ON TABLE visits IS '拜访核心业务表，保存计划、签到、签退、时长和合规结果';
COMMENT ON COLUMN visits.status IS '拜访流程状态：PLANNED/CHECKED_IN/CHECKED_OUT/COMPLETED/CANCELLED';
COMMENT ON COLUMN visits.compliance_status IS '合规状态：PENDING/COMPLIANT/NON_COMPLIANT';
COMMENT ON COLUMN visits.anomaly_reasons IS '异常原因机器码数组，如 ["DURATION_TOO_SHORT", "CHECK_IN_TOO_FAR"]';
