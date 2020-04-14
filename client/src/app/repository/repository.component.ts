import { Component, ChangeDetectionStrategy, Input } from '@angular/core';
import { Repository, GithubNotification } from '../interface/repository.interface';
import { MatDialog } from '@angular/material/dialog';
import { NotificationsDialogComponent } from './notifications-dialog/notifications-dialog.component';

@Component({
  selector: 'app-repository',
  templateUrl: './repository.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RepositoryComponent {
    @Input() repository: Repository;

    constructor(public dialog: MatDialog) {}

    showNotifications(notifications: Array<GithubNotification>) {
      console.log('notifications', notifications);
      if (notifications.length > 0) {
        const dialogRef = this.dialog.open(NotificationsDialogComponent, {
          width: '500px',
          data: { notifications }
        });
      }
    }
}