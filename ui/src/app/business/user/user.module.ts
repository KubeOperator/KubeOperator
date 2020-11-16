import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {UserComponent} from "./user.component";
import {UserCreateComponent} from "./user-create/user-create.component";
import {UserListComponent} from "./user-list/user-list.component";
import {UserUpdateComponent} from "./user-update/user-update.component";
import {UserDeleteComponent} from "./user-delete/user-delete.component";


@NgModule({
    declarations: [
        UserComponent, UserCreateComponent, UserListComponent, UserUpdateComponent, UserDeleteComponent
    ],
    imports: [
        CommonModule
    ]
})
export class UserModule {
}
