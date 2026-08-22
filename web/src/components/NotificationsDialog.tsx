import { useTranslation } from 'react-i18next'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import type { GithubNotification } from '@/lib/types'

interface NotificationsDialogProps {
  notifications: Array<GithubNotification>
  repositoryUrl: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export default function NotificationsDialog({
  notifications,
  repositoryUrl,
  open,
  onOpenChange,
}: NotificationsDialogProps) {
  const { t } = useTranslation()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>🔓 {t('dialog.githubNotifications')}</DialogTitle>
        </DialogHeader>

        <ul className="flex flex-col gap-3">
          {notifications.map((notification) => (
            <li key={notification.url}>
              <a
                href={`${repositoryUrl}/network/alerts`}
                target="_blank"
                rel="noreferrer"
                className="text-sm font-medium text-primary underline-offset-4 hover:underline"
              >
                {notification.title}
              </a>
              <p className="mt-0.5 text-xs text-muted-foreground">
                {new Date(notification.updatedAt).toLocaleString()}
              </p>
            </li>
          ))}
        </ul>

        <DialogFooter>
          <Button onClick={() => onOpenChange(false)}>
            {t('dialog.close')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
