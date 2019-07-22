import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {SettingComponent} from './setting.component';
import {TipModule} from '../tip/tip.module';
import {CoreModule} from '../core/core.module';
import {SystemSettingComponent} from './system-setting/system-setting.component';

@NgModule({
  declarations: [SettingComponent, SystemSettingComponent],
  imports: [
    CommonModule,
    CoreModule,
    TipModule,

  ]
})
export class SettingModule {
}
