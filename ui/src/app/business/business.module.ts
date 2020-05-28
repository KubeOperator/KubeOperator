import {NgModule} from '@angular/core';
import {CoreModule} from '../core/core.module';
import {ClusterModule} from './cluster/cluster.module';
import {SettingModule} from './setting/setting.module';
import {HostModule} from './host/host.module';


@NgModule({
    declarations: [],
    imports: [
        CoreModule,
        ClusterModule,
        SettingModule,
        HostModule,
    ],
    exports: [
        ClusterModule,
        SettingModule,
        HostModule,
    ]
})
export class BusinessModule {
}
