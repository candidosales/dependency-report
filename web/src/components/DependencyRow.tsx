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
    <div className="flex flex-col gap-2 rounded-lg bg-muted/60 p-3">
      <p className="m-0 flex items-center gap-2 text-sm font-medium">
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
        <span className="truncate">{name}</span>
      </p>

      <ul className="flex flex-wrap gap-3">
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
