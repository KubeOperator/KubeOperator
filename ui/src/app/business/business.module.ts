import {NgModule} from '@angular/core';
import {CoreModule} from '../core/core.module';
import {ClusterModule} from './cluster/cluster.module';
import {SettingModule} from './setting/setting.module';
import {HostModule} from './host/host.module';
import { UserComponent } from './user/user.component';
import { UserCreateComponent } from './user/user-create/user-create.component';
import { UserListComponent } from './user/user-list/user-list.component';
import { UserUpdateComponent } from './user/user-update/user-update.component';
import { UserDeleteComponent } from './user/user-delete/user-delete.component';
import {SharedModule} from '../shared/shared.module';
import { DeployPlanComponent } from './deploy-plan/deploy-plan.component';
import { RegionComponent } from './deploy-plan/region/region.component';
import {RouterModule} from "@angular/router";
import { RegionListComponent } from './deploy-plan/region/region-list/region-list.component';
import { RegionCreateComponent } from './deploy-plan/region/region-create/region-create.component';
import { RegionDeleteComponent } from './deploy-plan/region/region-delete/region-delete.component';


@NgModule({
    declarations: [UserComponent, UserCreateComponent, UserListComponent, UserUpdateComponent, UserDeleteComponent, DeployPlanComponent, RegionComponent, RegionListComponent, RegionCreateComponent, RegionDeleteComponent],
    imports: [
        CoreModule,
        ClusterModule,
        SettingModule,
        HostModule,
        SharedModule,
        RouterModule,
    ],
    exports: [
        ClusterModule,
        SettingModule,
        HostModule,
    ]
})
export class BusinessModule {
}
