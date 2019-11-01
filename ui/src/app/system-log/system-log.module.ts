import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {SystemLogComponent} from './system-log.component';
import {SystemLogListComponent} from './system-log-list/system-log-list.component';
import {CoreModule} from '../core/core.module';
import {SharedModule} from '../shared/shared.module';


@NgModule({
  declarations: [SystemLogComponent, SystemLogListComponent],
  imports: [
    CommonModule,
    CoreModule,
  ]
})
export class SystemLogModule {
}
