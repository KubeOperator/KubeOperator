import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ClusterComponent} from './cluster.component';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {CoreModule} from '../core/core.module';
import {ClusterService} from './cluster.service';
import {ClusterDetailComponent} from './cluster-detail/cluster-detail.component';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';

import {ClusterRoutingResolverService} from './cluster-routing-resolver.service';
import {HostsFilterPipe} from './hosts-filter.pipe';
import {DeviceCheckService} from './device-check.service';
import {SharedModule} from '../shared/shared.module';

@NgModule({
  declarations: [ClusterComponent, ClusterListComponent, ClusterDetailComponent, ClusterCreateComponent, HostsFilterPipe],
  imports: [
    CommonModule,

    CoreModule,
    SharedModule
  ],
  providers: [ClusterService, ClusterRoutingResolverService, DeviceCheckService]
})
export class ClusterModule {
}
