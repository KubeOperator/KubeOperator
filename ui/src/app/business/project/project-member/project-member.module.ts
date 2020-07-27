import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ProjectMemberListComponent } from './project-member-list/project-member-list.component';
import { ProjectMemberCreateComponent } from './project-member-create/project-member-create.component';
import { ProjectMemberDeleteComponent } from './project-member-delete/project-member-delete.component';
import {CoreModule} from '../../../core/core.module';
import {SharedModule} from '../../../shared/shared.module';



@NgModule({
    declarations: [ProjectMemberListComponent, ProjectMemberCreateComponent, ProjectMemberDeleteComponent],
    exports: [
        ProjectMemberListComponent,
        ProjectMemberCreateComponent,
        ProjectMemberDeleteComponent
    ],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule
    ]
})
export class ProjectMemberModule { }
