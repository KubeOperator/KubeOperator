import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {HostComponent} from './host.component';
import {TipModule} from '../tip/tip.module';
import {CoreModule} from '../core/core.module';
import {HostListComponent} from './host-list/host-list.component';
import {HostCreateComponent} from './host-create/host-create.component';
import {HostFilterPipe} from './host-filter.pipe';

@NgModule({
  declarations: [HostComponent, HostListComponent, HostCreateComponent],
  imports: [
    CommonModule,
    TipModule,
    CoreModule
  ]
})
export class HostModule {
}
