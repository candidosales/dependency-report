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
    <div className="flex flex-col overflow-hidden rounded-lg bg-muted/60 transition hover:-translate-y-px hover:shadow-md">
      <div className="p-2">
        <p className="truncate pr-2 text-sm font-medium">
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
        <div className="mt-1 flex items-start justify-between gap-2">
          <span className="flex items-center gap-1 text-xs text-muted-foreground">
            {getIconLibraryByFilter(repository.filter) && (
              <img
                src={getIconLibraryByFilter(repository.filter)}
                alt=""
                className="size-3"
              />
            )}
            {getVersionByFilter(repository.filter)}
          </span>
          <span className="text-xs text-muted-foreground">
            {repository.version}
          </span>
        </div>
      </div>

      <div className="mt-auto flex h-9 items-center justify-around bg-accent/60">
          {repository.documentation.frontend && (
            <Tooltip>
              <TooltipTrigger
                render={
                  <a
                    href={repository.documentation.frontend}
                    target="_blank"
                    rel="noreferrer"
                    className="flex size-7 items-center justify-center rounded-md text-muted-foreground hover:text-foreground"
                  />
                }
              >
                <FileTextIcon className="size-4" />
              </TooltipTrigger>
              <TooltipContent>Frontend Doc</TooltipContent>
            </Tooltip>
          )}

          {repository.documentation.design &&
            getIconDesignByUrl(repository.documentation.design) && (
            <Tooltip>
              <TooltipTrigger
                render={
                  <a
                    href={repository.documentation.design}
                    target="_blank"
                    rel="noreferrer"
                    className="flex size-7 items-center justify-center rounded-md text-muted-foreground hover:text-foreground"
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
                    className="flex size-7 items-center justify-center rounded-md"
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
                    className="flex size-7 items-center justify-center rounded-md"
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
        repositoryUrl={repository.url}
        open={notificationsOpen}
        onOpenChange={setNotificationsOpen}
      />
    </div>
  )
}
