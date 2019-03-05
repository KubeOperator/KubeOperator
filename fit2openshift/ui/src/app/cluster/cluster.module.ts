import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ClusterComponent} from './cluster.component';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {CoreModule} from '../core/core.module';
import {ClusterService} from './cluster.service';
import {ClusterDetailComponent} from './cluster-detail/cluster-detail.component';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {TipModule} from '../tip/tip.module';
import {ClusterRoutingResolverService} from './cluster-routing-resolver.service';
import { HostsFilterPipe } from './hosts-filter.pipe';

@NgModule({
  declarations: [ClusterComponent, ClusterListComponent, ClusterDetailComponent, ClusterCreateComponent, HostsFilterPipe],
  imports: [
    CommonModule,
    TipModule,
    CoreModule
  ],
  providers: [ClusterService, ClusterRoutingResolverService]
})
export class ClusterModule {
}
