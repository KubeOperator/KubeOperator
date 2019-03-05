import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {SettingComponent} from './setting.component';
import {TipModule} from '../tip/tip.module';
import {CoreModule} from '../core/core.module';

@NgModule({
  declarations: [SettingComponent],
  imports: [
    CommonModule,
    CoreModule,
    TipModule,

  ]
})
export class SettingModule {
}
