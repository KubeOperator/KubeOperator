import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {HostComponent} from './host.component';

import {CoreModule} from '../core/core.module';
import {HostListComponent} from './host-list/host-list.component';
import {HostCreateComponent} from './host-create/host-create.component';
import {HostInfoComponent} from './host-info/host-info.component';
import {SharedModule} from '../shared/shared.module';
import {HostImportComponent} from './host-import/host-import.component';

@NgModule({
  declarations: [HostComponent, HostListComponent, HostCreateComponent, HostInfoComponent, HostImportComponent],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule
  ]
})
export class HostModule {
}
