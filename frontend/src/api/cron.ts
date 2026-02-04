import request from '@/utils/request'

export interface CronTask {
  id: number
  name: string
  cron_expr: string
  task_type: string
  account_ids: string
  status: number
  last_run: string | null
  next_run: string | null
}

export function getCronTasks(): Promise<CronTask[]> {
  return request.get<CronTask[], CronTask[]>('/cron')
}

export function createCronTask(data: Partial<CronTask>): Promise<CronTask> {
  return request.post<CronTask, CronTask>('/cron', data)
}

export function updateCronTask(id: number, data: Partial<CronTask>): Promise<CronTask> {
  return request.put<CronTask, CronTask>(`/cron/${id}`, data)
}

export function deleteCronTask(id: number): Promise<null> {
  return request.delete<null, null>(`/cron/${id}`)
}

export function triggerCronTask(id: number): Promise<null> {
  return request.post<null, null>(`/cron/${id}/trigger`)
}
