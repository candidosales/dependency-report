import { Component, ChangeDetectionStrategy, Input } from '@angular/core';
import { Repository, GithubNotification, Update } from '../interface/repository.interface';
import { MatDialog } from '@angular/material/dialog';
import { NotificationsDialogComponent } from './notifications-dialog/notifications-dialog.component';
import { UpdatesDialogComponent } from './updates-dialog/updates-dialog.component';

@Component({
  selector: 'app-repository',
  templateUrl: './repository.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RepositoryComponent {
    @Input() repository: Repository;

    constructor(public dialog: MatDialog) {}

    showNotifications(notifications: Array<GithubNotification>) {
      if (notifications.length > 0) {
        this.dialog.open(NotificationsDialogComponent, {
          width: '500px',
          data: { notifications }
        });
      }
    }

    showUpdates(updates: Array<Update>) {
      if (updates.length > 0) {
        this.dialog.open(UpdatesDialogComponent, {
          width: '500px',
          data: { updates }
        });
      }
    }
}