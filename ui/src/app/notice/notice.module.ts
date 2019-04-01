import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NoticeComponent} from './notice.component';
import {CoreModule} from '../core/core.module';

@NgModule({
  declarations: [NoticeComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class NoticeModule {
}
