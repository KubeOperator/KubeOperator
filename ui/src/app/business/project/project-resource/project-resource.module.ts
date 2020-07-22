import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ProjectResourceListComponent } from './project-resource-list/project-resource-list.component';
import { ProjectResourceCreateComponent } from './project-resource-create/project-resource-create.component';
import { ProjectResourceDeleteComponent } from './project-resource-delete/project-resource-delete.component';
import {CoreModule} from '../../../core/core.module';
import {SharedModule} from '../../../shared/shared.module';



@NgModule({
    declarations: [ProjectResourceListComponent, ProjectResourceCreateComponent, ProjectResourceDeleteComponent],
    exports: [
        ProjectResourceListComponent,
        ProjectResourceCreateComponent,
        ProjectResourceDeleteComponent
    ],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule
    ]
})
export class ProjectResourceModule { }
