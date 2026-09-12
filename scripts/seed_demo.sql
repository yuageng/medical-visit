BEGIN;

INSERT INTO medical_representatives (id, employee_no, name)
VALUES ('00000000-0000-0000-0000-000000000001', 'MR-021', '刘慧娟')
ON CONFLICT (id) DO UPDATE
SET employee_no = EXCLUDED.employee_no,
    name = EXCLUDED.name,
    updated_at = NOW();

INSERT INTO hospitals (id, name, address, latitude, longitude)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '北京市第一人民医院',
    '北京市海淀区',
    31.230400,
    121.473700
)
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    address = EXCLUDED.address,
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude,
    updated_at = NOW();

INSERT INTO departments (id, hospital_id, name)
VALUES (
    '00000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000002',
    '心内科'
)
ON CONFLICT (id) DO UPDATE
SET hospital_id = EXCLUDED.hospital_id,
    name = EXCLUDED.name,
    updated_at = NOW();

INSERT INTO hcps (id, name, title, hospital_id, department_id)
VALUES (
    '00000000-0000-0000-0000-000000000004',
    '耿元贞',
    '主任医师',
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000003'
)
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    title = EXCLUDED.title,
    hospital_id = EXCLUDED.hospital_id,
    department_id = EXCLUDED.department_id,
    updated_at = NOW();

INSERT INTO products (id, code, name)
VALUES
    ('00000000-0000-0000-0000-000000000005', 'PROD-XNP', '心宁平'),
    ('00000000-0000-0000-0000-000000000006', 'PROD-SMA', '舒脉安')
ON CONFLICT (id) DO UPDATE
SET code = EXCLUDED.code,
    name = EXCLUDED.name,
    updated_at = NOW();

INSERT INTO academic_materials (id, code, name, product_id, active)
VALUES
    (
        '00000000-0000-0000-0000-000000000007',
        'MAT-XNP-001',
        '心宁平临床研究摘要',
        '00000000-0000-0000-0000-000000000005',
        TRUE
    ),
    (
        '00000000-0000-0000-0000-000000000008',
        'MAT-SMA-001',
        '舒脉安学术资料',
        '00000000-0000-0000-0000-000000000006',
        TRUE
    )
ON CONFLICT (id) DO UPDATE
SET code = EXCLUDED.code,
    name = EXCLUDED.name,
    product_id = EXCLUDED.product_id,
    active = EXCLUDED.active,
    updated_at = NOW();

COMMIT;
