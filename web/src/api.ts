import type { CheckInRequest, CheckOutRequest, CreateVisitRequest, DashboardResponse, Report, SaveReportRequest, Visit } from './types'

const API_BASE = import.meta.env.VITE_API_BASE ?? ''

type ApiEnvelope<T> = { code: string; message: string; data?: T; details?: unknown; request_id?: string }

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
  })
  const body = await response.json().catch(() => ({} as ApiEnvelope<T>)) as ApiEnvelope<T>
  if (!response.ok) throw new Error(body.message || `请求失败（${response.status}）`)
  return body.data as T
}

export const api = {
  health: () => request<{ status: string }>('/healthz'),
  listVisits: (mrId: string) => request<Visit[]>(`/api/v1/visits?mr_id=${encodeURIComponent(mrId)}`),
  createVisit: (payload: CreateVisitRequest) => request<Visit>('/api/v1/visits', { method: 'POST', body: JSON.stringify(payload) }),
  getVisit: (id: string) => request<Visit>(`/api/v1/visits/${id}`),
  checkIn: (id: string, payload: CheckInRequest) => request<Visit>(`/api/v1/visits/${id}/check-in`, { method: 'POST', body: JSON.stringify(payload) }),
  checkOut: (id: string, payload: CheckOutRequest) => request<Visit>(`/api/v1/visits/${id}/check-out`, { method: 'POST', body: JSON.stringify(payload) }),
  saveReport: (id: string, payload: SaveReportRequest) => request<Report>(`/api/v1/visits/${id}/report`, { method: 'POST', body: JSON.stringify(payload) }),
  dashboard: (month: string) => request<DashboardResponse>(`/api/v1/dashboard/visits/monthly?month=${encodeURIComponent(month)}`),
}
