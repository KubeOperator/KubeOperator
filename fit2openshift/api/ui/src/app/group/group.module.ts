import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CoreModule } from '../core/core.module';
import { SharedModule } from '../shared/shared.module';
import { GroupComponent } from './group.component';
import { GroupListComponent } from './group-list/group-list.component';
import { GroupCreateComponent } from './group-create/group-create.component';
import { GroupEditComponent } from './group-edit/group-edit.component';

@NgModule({
  declarations: [GroupComponent, GroupListComponent, GroupCreateComponent, GroupEditComponent, ],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule
  ]
})
export class GroupModule { }
