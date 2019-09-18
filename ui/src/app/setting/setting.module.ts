import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {SettingComponent} from './setting.component';
import {TipModule} from '../tip/tip.module';
import {CoreModule} from '../core/core.module';
import {SystemSettingComponent} from './system-setting/system-setting.component';
import { BackupStorageSettingComponent } from './backup-storage-setting/backup-storage-setting.component';

@NgModule({
  declarations: [SettingComponent, SystemSettingComponent, BackupStorageSettingComponent],
  imports: [
    CommonModule,
    CoreModule,
    TipModule,

  ]
})
export class SettingModule {
}
