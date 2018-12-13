import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CoreModule } from '../core/core.module';
import { SharedModule } from '../shared/shared.module';
import { HostComponent } from './host.component';
import { HostListComponent } from './host-list/host-list.component';
import { HostCreateComponent } from './host-create/host-create.component';

@NgModule({
  declarations: [HostListComponent, HostComponent, HostCreateComponent],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule
  ]
})
export class HostModule { }
