import { useTranslation } from 'react-i18next'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import type { Update } from '@/lib/types'

interface UpdatesDialogProps {
  updates: Array<Update>
  open: boolean
  onOpenChange: (open: boolean) => void
}

export default function UpdatesDialog({
  updates,
  open,
  onOpenChange,
}: UpdatesDialogProps) {
  const { t } = useTranslation()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>🔄 {t('dialog.dependenciesToUpdate')}</DialogTitle>
        </DialogHeader>

        <ul className="flex flex-col gap-2">
          {updates.map((update) => (
            <li
              key={update.name}
              className="flex items-center justify-between gap-3 text-sm"
            >
              <span className="font-medium">{update.name}</span>
              <span className="text-muted-foreground">
                {update.current} → {update.update}
              </span>
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
