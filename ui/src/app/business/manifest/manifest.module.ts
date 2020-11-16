import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ManifestComponent} from "./manifest.component";
import {ManifestListComponent} from "./manifest-list/manifest-list.component";
import {ManifestDetailComponent} from "./manifest-detail/manifest-detail.component";
import {CoreModule} from "../../core/core.module";
import {SharedModule} from "../../shared/shared.module";


@NgModule({
    declarations: [
        ManifestComponent, ManifestListComponent, ManifestDetailComponent
    ],
    imports: [
        CoreModule,
        SharedModule
    ]
})
export class ManifestModule {
}
