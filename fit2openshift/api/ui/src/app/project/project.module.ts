import { NgModule } from '@angular/core';

import { CoreModule } from '../core/core.module';
import { SharedModule } from '../shared/shared.module';
import { BaseModule } from '../base/base.module';

import { ProjectComponent } from './project.component';
import { ProjectListComponent } from './project-list/project-list.component';
import { ProjectDetailComponent } from './project-detail/project-detail.component';
import { ProjectOverviewComponent } from './project-overview/project-overview.component';
import { ProjectInventoryComponent } from './project-inventory/project-inventory.component';

import { ProjectRoutingResolver } from './project-routing-resolver.service';
import { ProjectCreateComponent } from './project-create/project-create.component';

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    BaseModule,
  ],
  declarations: [
    ProjectListComponent, ProjectDetailComponent, ProjectOverviewComponent,
    ProjectInventoryComponent, ProjectCreateComponent, ProjectComponent
  ],
  providers: [ProjectRoutingResolver]
})
export class ProjectModule { }
