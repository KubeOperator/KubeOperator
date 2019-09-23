import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BackupStorageCreateComponent } from './backup-storage-create/backup-storage-create.component';
import { BackupStorageListComponent } from './backup-storage-list/backup-storage-list.component';
import {ClrDatagridModule, ClrIconModule} from '@clr/angular';
import {SharedModule} from '../../shared/shared.module';
import {TipModule} from '../../tip/tip.module';
import { BackupStorageStatusPipe } from './backup-storage-status.pipe';



@NgModule({
  declarations: [BackupStorageCreateComponent, BackupStorageListComponent, BackupStorageStatusPipe],
  exports: [
    BackupStorageListComponent,
    BackupStorageCreateComponent
  ],
  imports: [
    CommonModule,
    ClrDatagridModule,
    ClrIconModule,
    SharedModule,
    TipModule
  ]
})
export class BackupStorageSettingModule { }
