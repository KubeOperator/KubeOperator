import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {SystemLogComponent} from './system-log.component';
import {SystemLogListComponent} from './system-log-list/system-log-list.component';
import {CoreModule} from '../core/core.module';
import {SharedModule} from '../shared/shared.module';
import { SystemLogDetailComponent } from './system-log-detail/system-log-detail.component';


@NgModule({
  declarations: [SystemLogComponent, SystemLogListComponent, SystemLogDetailComponent],
  imports: [
    CommonModule,
    CoreModule,
  ]
})
export class SystemLogModule {
}
