import { Component, inject } from '@angular/core';
import { Update } from '../../interface/repository.interface';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialogTitle, MatDialogContent, MatDialogActions } from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import { MatButton } from '@angular/material/button';

export interface DataDialog {
    updates: Array<Update>;
}

@Component({
    selector: 'app-updates-dialog',
    templateUrl: './updates-dialog.component.html',
    styleUrls: ['./updates-dialog.component.scss'],
    standalone: true,
    imports: [MatDialogTitle, MatDialogContent, MatDialogActions, MatButton, TranslateModule]
})
export class UpdatesDialogComponent {
  dialogRef = inject<MatDialogRef<UpdatesDialogComponent>>(MatDialogRef);
  data = inject(MAT_DIALOG_DATA);


  close(): void {
    this.dialogRef.close();
  }
}
