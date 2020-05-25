import {NgModule} from '@angular/core';
import {ClusterComponent} from './cluster.component';
import {CoreModule} from '../../core/core.module';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import { ClusterCreateComponent } from './cluster-create/cluster-create.component';
import { ClusterDeleteComponent } from './cluster-delete/cluster-delete.component';
import { ClusterDetailComponent } from './cluster-detail/cluster-detail.component';
import { OverviewComponent } from './cluster-detail/overview/overview.component';
import {RouterModule} from '@angular/router';


@NgModule({
    declarations: [ClusterComponent, ClusterListComponent, ClusterCreateComponent, ClusterDeleteComponent, ClusterDetailComponent, OverviewComponent],
    imports: [
        CoreModule,
        RouterModule,
    ]
})
export class ClusterModule {
}
