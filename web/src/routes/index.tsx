import { useEffect, useMemo, useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { RefreshCwIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import RepositoryCard from '@/components/RepositoryCard'
import DependencyRow from '@/components/DependencyRow'
import PieChart from '@/components/PieChart'
import {
  cachedReportQuery,
  liveReportQuery,
  serverStatusQuery,
} from '@/lib/queries'
import type { DependencyVersions } from '@/lib/types'

export const Route = createFileRoute('/')({ component: App })

type DependencyFilter = '*' | DependencyVersions['type']

function App() {
  const { t } = useTranslation()
  const [filter, setFilter] = useState<DependencyFilter>('*')

  const cachedQuery = useQuery(cachedReportQuery)
  const liveQuery = useQuery(liveReportQuery)
  const statusQuery = useQuery(serverStatusQuery)

  const serverIsDown = statusQuery.isError || statusQuery.data === false

  useEffect(() => {
    if (serverIsDown) {
      toast.error('Server is unreachable', {
        description: 'Live data refresh is disabled until the server is back.',
      })
    }
  }, [serverIsDown])

  useEffect(() => {
    if (liveQuery.isError) {
      console.error(liveQuery.error)
      toast.error('Failed to generate a live report', {
        description: 'Showing the last cached report instead.',
      })
    }
  }, [liveQuery.isError, liveQuery.error])

  const data = liveQuery.data ?? cachedQuery.data

  const dependenciesByVersions = useMemo(() => {
    if (!data) {
      return {}
    }
    if (filter === '*') {
      return data.dependenciesByVersions
    }
    return Object.fromEntries(
      Object.entries(data.dependenciesByVersions).filter(
        ([, dependency]) => dependency.type === filter,
      ),
    )
  }, [data, filter])

  if (cachedQuery.isLoading) {
    return (
      <main className="mx-auto flex max-w-6xl items-center justify-center px-4 py-24">
        <Spinner className="size-8" />
      </main>
    )
  }

  if (!data) {
    return null
  }

  return (
    <main className="mx-auto max-w-6xl px-4 py-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h2 className="text-xl font-semibold">{t('home.projects')}</h2>
        <div className="flex items-center gap-3 text-sm text-muted-foreground">
          <span>
            {t('home.navbar.generated')}:{' '}
            {new Date(data.generatedAt).toLocaleString()}
          </span>
          <Button
            variant="outline"
            disabled={serverIsDown || liveQuery.isFetching}
            onClick={() => liveQuery.refetch()}
          >
            {liveQuery.isFetching ? (
              <Spinner data-icon="inline-start" />
            ) : (
              <RefreshCwIcon data-icon="inline-start" />
            )}
            {t('home.navbar.refresh')}
          </Button>
        </div>
      </div>

      <div className="mt-4 grid gap-4 lg:grid-cols-[minmax(0,240px)_minmax(0,1fr)_minmax(0,320px)]">
        <div className="rounded-xl bg-card p-4 ring-1 ring-foreground/10">
          <h3 className="mb-3 text-sm font-semibold">
            {t('home.summary.title')}
          </h3>
          <div className="flex flex-col gap-4">
            <SummaryStat
              label={t('home.updated')}
              info={t('home.summary.infoUpdated')}
              emoji="👍"
              items={data.summary.updated}
            />
            <SummaryStat
              label={t('home.inconsistent')}
              info={t('home.summary.infoInconsistent')}
              emoji="🚧"
              items={data.summary.inconsistent}
            />
            <SummaryStat
              label={t('home.vulnerable')}
              info={t('home.summary.infoWithVulnerability')}
              emoji={data.summary.vulnerable.length > 0 ? '🔓' : '🔒'}
              items={data.summary.vulnerable}
            />
          </div>
        </div>

        <div className="rounded-xl bg-card p-4 ring-1 ring-foreground/10">
          <h3 className="mb-3 text-sm font-semibold">
            {data.projects.length} {t('home.graph.frontendProjects')}
          </h3>
          <div className="grid max-h-96 grid-cols-1 gap-2 overflow-y-auto pr-1 sm:grid-cols-2 xl:grid-cols-3">
            {data.projects.map((project) => (
              <RepositoryCard key={project.name} repository={project} />
            ))}
          </div>
        </div>

        <div className="rounded-xl bg-card p-4 ring-1 ring-foreground/10">
          <h3 className="mb-3 text-sm font-semibold">
            {t('home.graph.projectXLibrary')}
          </h3>
          <PieChart data={data.graphData?.projectsByFilters} />
        </div>
      </div>

      <h2 className="mt-8 text-xl font-semibold">{t('home.components')}</h2>
      <div className="mt-4 grid gap-4 lg:grid-cols-2">
        <div className="rounded-xl bg-card p-4 ring-1 ring-foreground/10">
          <h3 className="mb-3 text-sm font-semibold">
            {data.components.length} {t('home.components')}
          </h3>
          <div className="grid max-h-96 grid-cols-1 gap-2 overflow-y-auto pr-1 sm:grid-cols-2 xl:grid-cols-3">
            {data.components.map((component) => (
              <RepositoryCard key={component.name} repository={component} />
            ))}
          </div>
        </div>

        <div className="rounded-xl bg-card p-4 ring-1 ring-foreground/10">
          <h3 className="mb-3 text-sm font-semibold">
            {t('home.graph.componentXLibrary')}
          </h3>
          <PieChart data={data.graphData?.componentsByFilters} />
        </div>
      </div>

      <div className="mt-8 rounded-xl bg-card p-4 ring-1 ring-foreground/10">
        <h3 className="text-sm font-semibold">
          {t('home.dependencies.howOftenDependencyDifferentVersions')}
        </h3>
        <h4 className="mt-1 text-xs text-muted-foreground">
          {t('home.dependencies.tipReduce')}
        </h4>

        <div className="mt-3 flex flex-wrap gap-2">
          <LegendButton
            active={filter === 'good'}
            className="bg-status-good text-status-good-foreground"
            onClick={() => setFilter(filter === 'good' ? '*' : 'good')}
          >
            1 - 2
          </LegendButton>
          <LegendButton
            active={filter === 'warning'}
            className="bg-status-warning text-status-warning-foreground"
            onClick={() => setFilter(filter === 'warning' ? '*' : 'warning')}
          >
            3 - 5
          </LegendButton>
          <LegendButton
            active={filter === 'bad'}
            className="bg-status-bad text-status-bad-foreground"
            onClick={() => setFilter(filter === 'bad' ? '*' : 'bad')}
          >
            6 - 10
          </LegendButton>
          <LegendButton
            active={filter === 'terrible'}
            className="bg-status-terrible text-status-terrible-foreground"
            onClick={() => setFilter(filter === 'terrible' ? '*' : 'terrible')}
          >
            11+
          </LegendButton>
        </div>

        <div className="mt-4 grid max-h-96 grid-cols-1 gap-2 overflow-y-auto pr-1 sm:grid-cols-2 xl:grid-cols-3">
          {Object.entries(dependenciesByVersions).map(([name, dependency]) => (
            <DependencyRow key={name} name={name} dependency={dependency} />
          ))}
        </div>
      </div>
    </main>
  )
}

function SummaryStat({
  label,
  info,
  emoji,
  items,
}: {
  label: string
  info: string
  emoji: string
  items: Array<string>
}) {
  return (
    <div>
      <Tooltip>
        <TooltipTrigger render={<p className="m-0 text-xs text-muted-foreground" />}>
          {label}
        </TooltipTrigger>
        <TooltipContent>{info}</TooltipContent>
      </Tooltip>
      <Tooltip>
        <TooltipTrigger render={<p className="m-0 text-lg font-semibold" />}>
          {emoji} {items.length}
        </TooltipTrigger>
        <TooltipContent>{items.join(' / ') || '—'}</TooltipContent>
      </Tooltip>
    </div>
  )
}

function LegendButton({
  active,
  className,
  onClick,
  children,
}: {
  active: boolean
  className: string
  onClick: () => void
  children: React.ReactNode
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`rounded-full px-3 py-1 text-xs font-semibold transition ${className} ${
        active ? 'ring-2 ring-offset-2 ring-offset-background ring-foreground' : 'opacity-90 hover:opacity-100'
      }`}
    >
      {children}
    </button>
  )
}
