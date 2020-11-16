import {NgModule} from '@angular/core';
import {CoreModule} from '../core/core.module';
import {ClusterModule} from './cluster/cluster.module';
import {SettingModule} from './setting/setting.module';
import {HostModule} from './host/host.module';
import {SharedModule} from '../shared/shared.module';
import {RouterModule} from '@angular/router';
import {DeployPlanModule} from './deploy-plan/deploy-plan.module';
import {ProjectModule} from './project/project.module';
import {MessageCenterModule} from './message-center/message-center.module';
import {UserModule} from "./user/user.module";
import {ManifestModule} from "./manifest/manifest.module";
import {MultiClusterModule} from "./multi-cluster/multi-cluster.module";


@NgModule({
    declarations: [],
    imports: [
        CoreModule,
        ClusterModule,
        SettingModule,
        UserModule,
        HostModule,
        SharedModule,
        RouterModule,
        DeployPlanModule,
        ProjectModule,
        MessageCenterModule,
        MultiClusterModule,
        ManifestModule
    ],
    exports: [
        ClusterModule,
        UserModule,
        ManifestModule,
        SettingModule,
        HostModule,
        DeployPlanModule,
        MessageCenterModule,
    ]
})
export class BusinessModule {
}
