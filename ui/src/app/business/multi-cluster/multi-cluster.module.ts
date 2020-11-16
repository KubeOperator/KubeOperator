import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MultiClusterSettingComponent} from './multi-cluster-setting/multi-cluster-setting.component';
import {MultiClusterBrowserComponent} from './multi-cluster-browser/multi-cluster-browser.component';
import {MultiClusterComponent} from './multi-cluster.component';
import {CoreModule} from "../../core/core.module";
import {RouterModule} from "@angular/router";
import {SharedModule} from "../../shared/shared.module";


@NgModule({
    declarations: [MultiClusterSettingComponent, MultiClusterBrowserComponent, MultiClusterComponent],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
    ]
})
export class MultiClusterModule {
}
