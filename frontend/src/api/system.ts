import request from '@/utils/request'

export function getConfigs(category: string): Promise<Record<string, string>> {
  return request.get<Record<string, string>, Record<string, string>>(`/config/${category}`)
}

export function updateConfigs(category: string, data: Record<string, string>): Promise<null> {
  return request.put<null, null>(`/config/${category}`, data)
}

export function testTelegram(): Promise<{ message: string }> {
  return request.post<{ message: string }, { message: string }>('/config/telegram/test')
}

export interface CheckinLog {
  id: number
  account_id: number
  success: boolean
  message: string
  created_at: string
}

export interface CheckinLogSummary {
  logs: CheckinLog[]
  today_checkin_account_count: number
}

export function getLogs(): Promise<CheckinLogSummary> {
  return request.get('/logs')
}
