import {NgModule} from '@angular/core';
import {ClusterComponent} from './cluster.component';
import {CoreModule} from '../../core/core.module';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import { ClusterCreateComponent } from './cluster-create/cluster-create.component';
import { ClusterDeleteComponent } from './cluster-delete/cluster-delete.component';


@NgModule({
    declarations: [ClusterComponent, ClusterListComponent, ClusterCreateComponent, ClusterDeleteComponent],
    imports: [
        CoreModule,
    ]
})
export class ClusterModule {
}
