import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ProjectCreateComponent} from './project-create/project-create.component';
import {ProjectListComponent} from './project-list/project-list.component';
import {ProjectDeleteComponent} from './project-delete/project-delete.component';
import {ProjectUpdateComponent} from './project-update/project-update.component';
import {CoreModule} from '../../core/core.module';
import {SharedModule} from '../../shared/shared.module';
import {ProjectDetailComponent} from './project-detail/project-detail.component';
import {RouterModule} from '@angular/router';
import {ProjectResourceComponent} from './project-resource/project-resource.component';
import {ProjectResourceModule} from './project-resource/project-resource.module';
import {ProjectMemberComponent} from './project-member/project-member.component';
import {ProjectMemberModule} from './project-member/project-member.module';
import {ProjectComponent} from "./project.component";


@NgModule({
    declarations: [ProjectComponent, ProjectCreateComponent, ProjectListComponent, ProjectDeleteComponent,
        ProjectUpdateComponent, ProjectDetailComponent, ProjectResourceComponent, ProjectMemberComponent],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule,
        RouterModule,
        ProjectResourceModule,
        ProjectMemberModule
    ]
})
export class ProjectModule {
}
