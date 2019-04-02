import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {TipComponent} from './tip.component';
import {CoreModule} from '../core/core.module';

@NgModule({
  declarations: [TipComponent],
  imports: [
    CommonModule,
    CoreModule
  ], exports: [
    TipComponent
  ]
})
export class TipModule {
}
