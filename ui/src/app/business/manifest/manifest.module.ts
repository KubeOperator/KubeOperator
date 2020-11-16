import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ManifestComponent} from "./manifest.component";
import {ManifestListComponent} from "./manifest-list/manifest-list.component";
import {ManifestDetailComponent} from "./manifest-detail/manifest-detail.component";


@NgModule({
    declarations: [
        ManifestComponent, ManifestListComponent, ManifestDetailComponent
    ],
    imports: [
        CommonModule
    ]
})
export class ManifestModule {
}
