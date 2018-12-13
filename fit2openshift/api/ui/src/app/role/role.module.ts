import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from '../shared/shared.module';
import { CoreModule } from '../core/core.module';
import { RoleListComponent } from './role-list/role-list.component';

@NgModule({
  imports: [
    CommonModule,
    CoreModule,
    SharedModule,
  ],
  declarations: [RoleListComponent]
})
export class RoleModule { }
