import { queryOptions } from '@tanstack/react-query'
import type { ReportData } from './types'

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

async function fetchJson<T>(url: string): Promise<T> {
  const res = await fetch(url)
  if (!res.ok) {
    throw new Error(`Request to ${url} failed with status ${res.status}`)
  }
  return res.json() as Promise<T>
}

function prepareReport(data: ReportData): ReportData {
  data.graphData?.projectsByFilters?.unshift(['Filter', 'Version'])
  data.graphData?.componentsByFilters?.unshift(['Filter', 'Version'])
  return data
}

export const cachedReportQuery = queryOptions({
  queryKey: ['report', 'cached'] as const,
  queryFn: () => fetchJson<ReportData>('/data.json').then(prepareReport),
  staleTime: Infinity,
})

export const liveReportQuery = queryOptions({
  queryKey: ['report', 'live'] as const,
  queryFn: () =>
    fetchJson<ReportData>(`${API_URL}/generate-report`).then(prepareReport),
  enabled: false,
  retry: false,
})

export const serverStatusQuery = queryOptions({
  queryKey: ['server-status'] as const,
  queryFn: () => fetchJson<{ ok: boolean }>(`${API_URL}/ping`).then((r) => r.ok),
  retry: false,
  refetchInterval: 30_000,
  refetchOnWindowFocus: false,
})
