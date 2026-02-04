import dayjs from 'dayjs'

export function formatTime(time: string | null | undefined, format = 'YYYY-MM-DD HH:mm:ss'): string {
  if (!time) return '-'
  return dayjs(time).format(format)
}

export function formatDate(time: string | null | undefined): string {
  return formatTime(time, 'YYYY-MM-DD')
}

export function fromNow(time: string | null | undefined): string {
  if (!time) return '-'
  const d = dayjs(time)
  const now = dayjs()
  const diff = now.diff(d, 'minute')

  if (diff < 1) return '刚刚'
  if (diff < 60) return `${diff}分钟前`
  if (diff < 1440) return `${Math.floor(diff / 60)}小时前`
  if (diff < 43200) return `${Math.floor(diff / 1440)}天前`
  return d.format('YYYY-MM-DD')
}
