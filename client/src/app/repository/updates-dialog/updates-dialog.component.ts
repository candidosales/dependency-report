import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Update } from '../../interface/repository.interface';

export interface DataDialog {
    updates: Array<Update>;
}

@Component({
  selector: 'app-updates-dialog',
  templateUrl: './updates-dialog.component.html',
  styleUrls: ['./updates-dialog.component.scss']
})
export class UpdatesDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<UpdatesDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}

  close(): void {
    this.dialogRef.close();
  }
}
