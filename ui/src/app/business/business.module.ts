import {NgModule} from '@angular/core';
import {CoreModule} from '../core/core.module';
import {ClusterModule} from './cluster/cluster.module';


@NgModule({
    declarations: [],
    imports: [
        CoreModule,
        ClusterModule
    ],
    exports: [
        ClusterModule
    ]
})
export class BusinessModule {
}
