-- 创建医药代表表
CREATE TABLE IF NOT EXISTS medical_representatives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_no VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_medical_representatives_employee_no ON medical_representatives(employee_no);

COMMENT ON TABLE medical_representatives IS '医药代表（MR）基础信息表';
COMMENT ON COLUMN medical_representatives.employee_no IS '员工编号，全局唯一';

-- 创建医院表
CREATE TABLE IF NOT EXISTS hospitals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    address VARCHAR(500),
    latitude NUMERIC(9,6),
    longitude NUMERIC(9,6),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_hospitals_latitude CHECK (latitude IS NULL OR (latitude >= -90 AND latitude <= 90)),
    CONSTRAINT chk_hospitals_longitude CHECK (longitude IS NULL OR (longitude >= -180 AND longitude <= 180)),
    CONSTRAINT chk_hospitals_coordinates_complete CHECK ((latitude IS NULL) = (longitude IS NULL))
);

CREATE INDEX idx_hospitals_name ON hospitals(name);

COMMENT ON TABLE hospitals IS '医院基础信息及合规定位基准表';
COMMENT ON COLUMN hospitals.latitude IS '医院基准纬度，用于签到签退距离合规计算';
COMMENT ON COLUMN hospitals.longitude IS '医院基准经度，用于签到签退距离合规计算';

-- 创建科室表
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_departments_hospital FOREIGN KEY (hospital_id) 
        REFERENCES hospitals(id) ON DELETE RESTRICT,
    CONSTRAINT uq_departments_hospital_name UNIQUE (hospital_id, name)
);

CREATE INDEX idx_departments_hospital_id ON departments(hospital_id);

COMMENT ON TABLE departments IS '医院科室表';

-- 创建医生表
CREATE TABLE IF NOT EXISTS hcps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    title VARCHAR(100),
    hospital_id UUID NOT NULL,
    department_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_hcps_hospital FOREIGN KEY (hospital_id) 
        REFERENCES hospitals(id) ON DELETE RESTRICT,
    CONSTRAINT fk_hcps_department FOREIGN KEY (department_id) 
        REFERENCES departments(id) ON DELETE RESTRICT
);

CREATE INDEX idx_hcps_hospital_id ON hcps(hospital_id);
CREATE INDEX idx_hcps_department_id ON hcps(department_id);
CREATE INDEX idx_hcps_name ON hcps(name);

COMMENT ON TABLE hcps IS '医生（HCP - Healthcare Professional）基础信息表';

-- 创建产品表
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_code ON products(code);
CREATE INDEX idx_products_name ON products(name);

COMMENT ON TABLE products IS '药品产品基础信息表';
