import {NgModule} from '@angular/core';
import {ClusterComponent} from './cluster.component';
import {CoreModule} from '../../core/core.module';
import {ClusterListComponent} from './cluster-list/cluster-list.component';


@NgModule({
    declarations: [ClusterComponent, ClusterListComponent],
    imports: [
        CoreModule,
    ]
})
export class ClusterModule {
}
