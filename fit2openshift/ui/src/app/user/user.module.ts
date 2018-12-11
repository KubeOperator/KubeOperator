import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {UserComponent} from './user.component';
import {UserListComponent} from './user-list/user-list.component';
import {CoreModule} from '../core/core.module';
import {UserService} from './user.service';

@NgModule({
  declarations: [UserComponent, UserListComponent],
  imports: [
    CommonModule,
    CoreModule
  ],
  providers: [UserService]
})
export class UserModule {
}
