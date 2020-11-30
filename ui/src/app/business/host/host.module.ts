import { NgModule } from '@angular/core';
import { HostComponent } from './host.component';
import { HostListComponent } from './host-list/host-list.component';
import { HostCreateComponent } from './host-create/host-create.component';
import { HostDeleteComponent } from './host-delete/host-delete.component';
import {CoreModule} from '../../core/core.module';
import { HostDetailComponent } from './host-detail/host-detail.component';
import { HostGrantComponent } from './host-grant/host-grant.component';
import {SharedModule} from '../../shared/shared.module';
import { HostStatusDetailComponent } from './host-status-detail/host-status-detail.component';
import { HostImportComponent } from './host-import/host-import.component';



@NgModule({
  declarations: [HostComponent, HostListComponent, HostCreateComponent, HostDeleteComponent, HostDetailComponent,HostGrantComponent, HostStatusDetailComponent, HostImportComponent],
    imports: [
        CoreModule,
        SharedModule,
    ]
})
export class HostModule { }
