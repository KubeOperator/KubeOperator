import {NgModule} from '@angular/core';
import {CoreModule} from '../core/core.module';
import {ClusterModule} from './cluster/cluster.module';
import {CredentialModule} from './setting/credential/credential.module';
import {SettingModule} from './setting/setting.module';


@NgModule({
    declarations: [],
    imports: [
        CoreModule,
        ClusterModule,
        CredentialModule,
        SettingModule,
    ],
    exports: [
        ClusterModule,
        SettingModule,
        CredentialModule,
    ]
})
export class BusinessModule {
}
