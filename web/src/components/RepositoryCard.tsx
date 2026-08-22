import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { FileTextIcon } from 'lucide-react'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import NotificationsDialog from './NotificationsDialog'
import UpdatesDialog from './UpdatesDialog'
import {
  getIconDesignByUrl,
  getIconLibraryByFilter,
  getVersionByFilter,
} from '@/lib/dependency-display'
import type { Repository } from '@/lib/types'

interface RepositoryCardProps {
  repository: Repository
}

export default function RepositoryCard({ repository }: RepositoryCardProps) {
  const { t } = useTranslation()
  const [notificationsOpen, setNotificationsOpen] = useState(false)
  const [updatesOpen, setUpdatesOpen] = useState(false)

  return (
    <div className="flex items-center justify-between gap-3 border-b border-border py-2.5 last:border-b-0">
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm font-medium">
          <a
            href={repository.url}
            title={repository.name}
            target="_blank"
            rel="noreferrer"
            className="text-foreground no-underline hover:underline"
          >
            {repository.name}
          </a>
        </p>
        <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
          <img
            src={getIconLibraryByFilter(repository.filter)}
            alt=""
            className="size-3.5"
          />
          <span>{getVersionByFilter(repository.filter)}</span>
          <span>&middot;</span>
          <span>{repository.version}</span>
        </div>
      </div>

      <div className="flex shrink-0 items-center gap-1">
        {repository.documentation.frontend && (
          <Tooltip>
            <TooltipTrigger
              render={
                <a
                  href={repository.documentation.frontend}
                  target="_blank"
                  rel="noreferrer"
                  className="flex size-7 items-center justify-center rounded-md text-muted-foreground hover:bg-muted hover:text-foreground"
                />
              }
            >
              <FileTextIcon className="size-4" />
            </TooltipTrigger>
            <TooltipContent>Frontend Doc</TooltipContent>
          </Tooltip>
        )}

        {repository.documentation.design && (
          <Tooltip>
            <TooltipTrigger
              render={
                <a
                  href={repository.documentation.design}
                  target="_blank"
                  rel="noreferrer"
                  className="flex size-7 items-center justify-center rounded-md text-muted-foreground hover:bg-muted hover:text-foreground"
                />
              }
            >
              <img
                src={getIconDesignByUrl(repository.documentation.design)}
                alt=""
                className="size-4"
              />
            </TooltipTrigger>
            <TooltipContent>Design Doc</TooltipContent>
          </Tooltip>
        )}

        {repository.updates.length > 0 && (
          <Tooltip>
            <TooltipTrigger
              render={
                <button
                  type="button"
                  onClick={() => setUpdatesOpen(true)}
                  className="flex size-7 items-center justify-center rounded-md hover:bg-muted"
                />
              }
            >
              🚧
            </TooltipTrigger>
            <TooltipContent>{t('dialog.dependenciesToUpdate')}</TooltipContent>
          </Tooltip>
        )}

        {repository.notifications.length > 0 && (
          <Tooltip>
            <TooltipTrigger
              render={
                <button
                  type="button"
                  onClick={() => setNotificationsOpen(true)}
                  className="flex size-7 items-center justify-center rounded-md hover:bg-muted"
                />
              }
            >
              🔓
            </TooltipTrigger>
            <TooltipContent>Vulnerability alert</TooltipContent>
          </Tooltip>
        )}
      </div>

      <UpdatesDialog
        updates={repository.updates}
        open={updatesOpen}
        onOpenChange={setUpdatesOpen}
      />
      <NotificationsDialog
        notifications={repository.notifications}
        open={notificationsOpen}
        onOpenChange={setNotificationsOpen}
      />
    </div>
  )
}
