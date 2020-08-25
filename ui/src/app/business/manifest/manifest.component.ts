import {Component, OnInit, ViewChild} from '@angular/core';
import {ManifestListComponent} from "./manifest-list/manifest-list.component";
import {ManifestDetailComponent} from "./manifest-detail/manifest-detail.component";
import {Manifest} from "./manifest";

@Component({
    selector: 'app-manifest',
    templateUrl: './manifest.component.html',
    styleUrls: ['./manifest.component.css']
})
export class ManifestComponent implements OnInit {

    constructor() {
    }

    @ViewChild(ManifestListComponent)
    list: ManifestListComponent;

    @ViewChild(ManifestDetailComponent)
    detail: ManifestDetailComponent;

    ngOnInit(): void {
    }

    openDetail(item: Manifest) {
        this.detail.open(item);
    }
}
