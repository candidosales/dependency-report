import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import type { DependencyVersions } from '@/lib/types'

const STATUS_CLASSES: Record<DependencyVersions['type'], string> = {
  good: 'bg-status-good text-status-good-foreground',
  warning: 'bg-status-warning text-status-warning-foreground',
  bad: 'bg-status-bad text-status-bad-foreground',
  terrible: 'bg-status-terrible text-status-terrible-foreground',
}

interface DependencyRowProps {
  name: string
  dependency: DependencyVersions
}

export default function DependencyRow({ name, dependency }: DependencyRowProps) {
  const { t } = useTranslation()
  const versionEntries = Object.entries(dependency.versions)
  const versionCount = versionEntries.length

  return (
    <div className="flex flex-wrap items-center gap-3 border-b border-border py-2.5 last:border-b-0">
      <p className="m-0 flex min-w-48 items-center gap-2 text-sm font-medium">
        <Tooltip>
          <TooltipTrigger
            render={
              <Badge className={STATUS_CLASSES[dependency.type]} />
            }
          >
            {versionCount}
          </TooltipTrigger>
          <TooltipContent>
            {versionCount} {t('home.dependencies.info')}
          </TooltipContent>
        </Tooltip>
        {name}
      </p>

      <ul className="flex flex-1 flex-wrap gap-4">
        {versionEntries.map(([version, info]) => (
          <Tooltip key={version}>
            <TooltipTrigger
              render={<li className="text-center text-xs" />}
            >
              <p className="m-0 font-semibold">{info.quantity}</p>
              <p className="m-0 text-muted-foreground">{version}</p>
            </TooltipTrigger>
            <TooltipContent>{info.projects.join(' / ')}</TooltipContent>
          </Tooltip>
        ))}
      </ul>
    </div>
  )
}
