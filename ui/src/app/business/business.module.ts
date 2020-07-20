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
import {RouterModule} from '@angular/router';
import {DeployPlanModule} from './deploy-plan/deploy-plan.module';
import { ProjectComponent } from './project/project.component';
import {ProjectModule} from './project/project.module';


@NgModule({
    declarations: [UserComponent, UserCreateComponent, UserListComponent, UserUpdateComponent, UserDeleteComponent, ProjectComponent],
    imports: [
        CoreModule,
        ClusterModule,
        SettingModule,
        HostModule,
        SharedModule,
        RouterModule,
        DeployPlanModule,
        ProjectModule,
    ],
    exports: [
        ClusterModule,
        SettingModule,
        HostModule,
        DeployPlanModule,
    ]
})
export class BusinessModule {
}
