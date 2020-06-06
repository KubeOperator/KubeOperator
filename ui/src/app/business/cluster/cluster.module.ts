import {NgModule} from '@angular/core';
import {ClusterComponent} from './cluster.component';
import {CoreModule} from '../../core/core.module';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {ClusterDeleteComponent} from './cluster-delete/cluster-delete.component';
import {ClusterDetailComponent} from './cluster-detail/cluster-detail.component';
import {OverviewComponent} from './cluster-detail/overview/overview.component';
import {RouterModule} from '@angular/router';
import {ClusterConditionComponent} from './cluster-condition/cluster-condition.component';
import {NodeComponent} from './cluster-detail/node/node.component';
import { NamespaceComponent } from './cluster-detail/namespace/namespace.component';
import { NamespaceListComponent } from './cluster-detail/namespace/namespace-list/namespace-list.component';
import {SharedModule} from "../../shared/shared.module";


@NgModule({
    declarations: [ClusterComponent, ClusterListComponent, ClusterCreateComponent, ClusterDeleteComponent, ClusterDetailComponent, OverviewComponent, ClusterConditionComponent, NodeComponent, NamespaceComponent, NamespaceListComponent],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
    ]
})
export class ClusterModule {
}
