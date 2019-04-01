import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {HostComponent} from './host.component';
import {TipModule} from '../tip/tip.module';
import {CoreModule} from '../core/core.module';
import {HostListComponent} from './host-list/host-list.component';
import {HostCreateComponent} from './host-create/host-create.component';
import { HostInfoComponent } from './host-info/host-info.component';

@NgModule({
  declarations: [HostComponent, HostListComponent, HostCreateComponent, HostInfoComponent],
  imports: [
    CommonModule,
    TipModule,
    CoreModule
  ]
})
export class HostModule {
}
