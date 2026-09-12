export type VisitStatus = 'PLANNED' | 'CHECKED_IN' | 'CHECKED_OUT' | 'COMPLETED' | 'CANCELLED'
export type ComplianceStatus = 'PENDING' | 'COMPLIANT' | 'NON_COMPLIANT'

export interface AnomalyReason {
  code: string
  message: string
}

export interface Visit {
  id: string
  mr_id: string
  mr_name?: string
  hcp_id: string
  hcp_name?: string
  hcp_title?: string
  hospital_id: string
  hospital_name?: string
  department_id: string
  department_name?: string
  product_id: string
  product_name?: string
  planned_start_at: string
  plan_note?: string
  status: VisitStatus
  checked_in_at?: string
  check_in_latitude?: number
  check_in_longitude?: number
  check_in_distance_m?: number
  checked_out_at?: string
  check_out_latitude?: number
  check_out_longitude?: number
  check_out_distance_m?: number
  duration_seconds?: number
  compliance_status: ComplianceStatus
  anomaly_reasons?: AnomalyReason[]
  created_at: string
  updated_at: string
}

export interface CreateVisitRequest {
  mr_id: string
  hcp_id: string
  hospital_id: string
  department_id: string
  product_id: string
  planned_start_at: string
  plan_note: string
}

export interface CheckInRequest {
  latitude: number
  longitude: number
  check_in_time: string
}

export interface CheckOutRequest {
  latitude: number
  longitude: number
  check_out_time: string
}

export interface SaveReportRequest {
  conversation_summary: string
  doctor_feedback: string
  materials_distributed: boolean
  material_ids: string[]
  additional_notes: string
}

export interface Report {
  id: string
  visit_id: string
  conversation_summary: string
  doctor_feedback?: string
  materials_distributed: boolean
  material_ids?: string[]
  additional_notes?: string
  created_at: string
  updated_at: string
}

export interface MonthlyProductStat {
  product_id: string
  product_name: string
  visit_count: number
  normal_count: number
  abnormal_count: number
}

export interface DashboardResponse {
  month: string
  items: MonthlyProductStat[]
}

export interface DemoOption {
  id: string
  name: string
  detail?: string
}

export const demoOptions = {
  mrs: [{ id: '00000000-0000-0000-0000-000000000001', name: '刘慧娟', detail: 'MR-021' }],
  hospitals: [{ id: '00000000-0000-0000-0000-000000000002', name: '北京市第一人民医院', detail: '北京市海淀区' }],
  departments: [{ id: '00000000-0000-0000-0000-000000000003', name: '心内科', detail: '目标科室' }],
  hcps: [{ id: '00000000-0000-0000-0000-000000000004', name: '耿元贞', detail: '主任医师 · 心内科' }],
  products: [
    { id: '00000000-0000-0000-0000-000000000005', name: '心宁平', detail: '抗凝治疗产品' },
    { id: '00000000-0000-0000-0000-000000000006', name: '舒脉安', detail: '心血管产品' },
  ],
}
