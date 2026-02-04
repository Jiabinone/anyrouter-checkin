import request from '@/utils/request'

export interface Account {
  id: number
  name: string
  user_id: number
  username: string
  role: number
  status: number
  last_checkin: string | null
  last_result: string | null
}

export interface SessionInfo {
  user_id: number
  username: string
  role: number
  status: number
  group: string
}

export interface CheckinResult {
  success: boolean
  result: string
}

export function getAccounts(): Promise<Account[]> {
  return request.get<Account[], Account[]>('/accounts')
}

export function createAccount(data: { name: string; session: string }): Promise<Account> {
  return request.post<Account, Account>('/accounts', data)
}

export function updateAccount(id: number, data: { name: string; session: string }): Promise<Account> {
  return request.put<Account, Account>(`/accounts/${id}`, data)
}

export function deleteAccount(id: number): Promise<null> {
  return request.delete<null, null>(`/accounts/${id}`)
}

export function checkinAccount(id: number): Promise<CheckinResult> {
  return request.post<CheckinResult, CheckinResult>(`/accounts/${id}/checkin`)
}

export function verifySession(session: string): Promise<SessionInfo> {
  return request.post<SessionInfo, SessionInfo>('/accounts/verify', { session })
}
