export interface Documentation {
  frontend: string
  design: string
}

export interface Update {
  name: string
  current: string
  update: string
}

export interface GithubNotification {
  reason: string
  title: string
  url: string
  updatedAt: string
}

export interface Repository {
  name: string
  version: string
  filter: string
  url: string
  documentation: Documentation
  notifications: Array<GithubNotification>
  updates: Array<Update>
}

export interface SummaryData {
  updated: Array<string>
  inconsistent: Array<string>
  vulnerable: Array<string>
}

export interface GraphData {
  componentsByFilters: Array<[string, unknown]>
  projectsByFilters: Array<[string, unknown]>
}

export interface DependencyVersionInfo {
  quantity: number
  projects: Array<string>
}

export interface DependencyVersions {
  type: 'good' | 'warning' | 'bad' | 'terrible'
  versions: Record<string, DependencyVersionInfo>
}

export interface ReportData {
  generatedAt: string
  summary: SummaryData
  projects: Array<Repository>
  components: Array<Repository>
  dependenciesByVersions: Record<string, DependencyVersions>
  graphData?: GraphData
}
