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


@NgModule({
    declarations: [UserComponent, UserCreateComponent, UserListComponent, UserUpdateComponent, UserDeleteComponent],
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
