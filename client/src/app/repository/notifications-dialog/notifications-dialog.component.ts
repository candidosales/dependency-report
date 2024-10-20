import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialogTitle, MatDialogContent, MatDialogActions } from '@angular/material/dialog';

import { GithubNotification } from '../../interface/repository.interface';
import { TranslateModule } from '@ngx-translate/core';
import { DatePipe } from '@angular/common';
import { MatButton } from '@angular/material/button';

export interface DataDialog {
    notifications: Array<GithubNotification>;
}

@Component({
    selector: 'app-notifications-dialog',
    templateUrl: './notifications-dialog.component.html',
    styleUrls: ['./notifications-dialog.component.scss'],
    standalone: true,
    imports: [MatDialogTitle, MatDialogContent, MatDialogActions, MatButton, DatePipe, TranslateModule]
})
export class NotificationsDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<NotificationsDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}

  close(): void {
    this.dialogRef.close();
  }
}
