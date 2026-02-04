import request from '@/utils/request'

export interface LoginResponse {
  token: string
}

export interface ProfileResponse {
  username: string
}

export function login(username: string, password: string): Promise<LoginResponse> {
  return request.post<LoginResponse, LoginResponse>('/auth/login', { username, password })
}

export function getProfile(): Promise<ProfileResponse> {
  return request.get<ProfileResponse, ProfileResponse>('/auth/profile')
}

export function changePassword(oldPassword: string, newPassword: string): Promise<null> {
  return request.put<null, null>('/auth/password', { old_password: oldPassword, new_password: newPassword })
}
