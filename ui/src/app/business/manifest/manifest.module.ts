import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ManifestComponent} from './manifest.component';
import {ManifestListComponent} from './manifest-list/manifest-list.component';
import {ManifestDetailComponent} from './manifest-detail/manifest-detail.component';
import {CoreModule} from '../../core/core.module';
import {SharedModule} from '../../shared/shared.module';
import { ManifestAlertComponent } from './manifest-alert/manifest-alert.component';


@NgModule({
    declarations: [
        ManifestComponent, ManifestListComponent, ManifestDetailComponent, ManifestAlertComponent
    ],
    imports: [
        CoreModule,
        SharedModule
    ]
})
export class ManifestModule {
}
