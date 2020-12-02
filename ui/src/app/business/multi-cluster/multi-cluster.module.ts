import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MultiClusterSettingComponent} from './multi-cluster-repository-detail/multi-cluster-setting/multi-cluster-setting.component';
import {MultiClusterBrowserComponent} from './multi-cluster-repository-detail/multi-cluster-browser/multi-cluster-browser.component';
import {MultiClusterComponent} from './multi-cluster.component';
import {CoreModule} from "../../core/core.module";
import {RouterModule} from "@angular/router";
import {SharedModule} from "../../shared/shared.module";
import {MultiClusterRepositoryListComponent} from './multi-cluster-repository-list/multi-cluster-repository-list.component';
import {MultiClusterRepositoryCreateComponent} from './multi-cluster-repository-create/multi-cluster-repository-create.component';
import {MultiClusterRepositoryDetailComponent} from './multi-cluster-repository-detail/multi-cluster-repository-detail.component';
import {CodemirrorModule} from "ng2-codemirror";
import { FileCreateComponent } from './multi-cluster-repository-detail/multi-cluster-browser/file-create/file-create.component';
import { FileDeleteComponent } from './multi-cluster-repository-detail/multi-cluster-browser/file-delete/file-delete.component';
import { MultiClusterRelationComponent } from './multi-cluster-repository-detail/multi-cluster-relation/multi-cluster-relation.component';
import { MultiClusterRelationListComponent } from './multi-cluster-repository-detail/multi-cluster-relation/multi-cluster-relation-list/multi-cluster-relation-list.component';
import { MultiClusterRelationCreateComponent } from './multi-cluster-repository-detail/multi-cluster-relation/multi-cluster-relation-create/multi-cluster-relation-create.component';
import { MultiClusterRelationDeleteComponent } from './multi-cluster-repository-detail/multi-cluster-relation/multi-cluster-relation-delete/multi-cluster-relation-delete.component';
import { MultiClusterLogComponent } from './multi-cluster-repository-detail/multi-cluster-log/multi-cluster-log.component';
import { MultiClusterRepositoryDeleteComponent } from './multi-cluster-repository-delete/multi-cluster-repository-delete.component';
import { MultiClusterLogListComponent } from './multi-cluster-repository-detail/multi-cluster-log/multi-cluster-log-list/multi-cluster-log-list.component';
import { MultiClusterLogDetailComponent } from './multi-cluster-repository-detail/multi-cluster-log/multi-cluster-log-detail/multi-cluster-log-detail.component';


@NgModule({
    declarations: [MultiClusterSettingComponent, MultiClusterBrowserComponent,
        MultiClusterComponent, MultiClusterRepositoryListComponent,
        MultiClusterRepositoryCreateComponent, MultiClusterRepositoryDetailComponent, FileCreateComponent, FileDeleteComponent, MultiClusterRelationComponent, MultiClusterRelationListComponent, MultiClusterRelationCreateComponent, MultiClusterRelationDeleteComponent, MultiClusterLogComponent, MultiClusterRepositoryDeleteComponent, MultiClusterLogListComponent, MultiClusterLogDetailComponent],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
        CodemirrorModule
    ]
})
export class MultiClusterModule {
}
