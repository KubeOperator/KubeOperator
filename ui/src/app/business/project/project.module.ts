import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ProjectCreateComponent } from './project-create/project-create.component';
import { ProjectListComponent } from './project-list/project-list.component';
import { ProjectDeleteComponent } from './project-delete/project-delete.component';
import { ProjectUpdateComponent } from './project-update/project-update.component';
import {CoreModule} from '../../core/core.module';
import {SharedModule} from '../../shared/shared.module';



@NgModule({
  declarations: [ProjectCreateComponent, ProjectListComponent, ProjectDeleteComponent, ProjectUpdateComponent],
  exports: [
    ProjectListComponent,
    ProjectCreateComponent,
    ProjectDeleteComponent,
    ProjectUpdateComponent
  ],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule
  ]
})
export class ProjectModule { }
