import { Component, ChangeDetectionStrategy, inject, input } from '@angular/core';
import { Repository, GithubNotification, Update } from '../interface/repository.interface';
import { MatDialog } from '@angular/material/dialog';
import { NotificationsDialogComponent } from './notifications-dialog/notifications-dialog.component';
import { UpdatesDialogComponent } from './updates-dialog/updates-dialog.component';
import { VersionPipe } from '../pipes/version.pipe';
import { IconLibraryPipe } from '../pipes/icon-library.pipe';
import { IconDesignPipe } from '../pipes/icon-design.pipe';
import { TranslateModule } from '@ngx-translate/core';
import { MatIcon } from '@angular/material/icon';
import { MatTooltip } from '@angular/material/tooltip';

@Component({
    selector: 'app-repository',
    templateUrl: './repository.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: true,
    imports: [MatTooltip, MatIcon, TranslateModule, IconDesignPipe, IconLibraryPipe, VersionPipe]
})
export class RepositoryComponent {
    dialog = inject(MatDialog);

    repository = input<Repository>();

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