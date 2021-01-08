import {Component, OnInit, ViewChild} from '@angular/core';
import {ManifestListComponent} from './manifest-list/manifest-list.component';
import {ManifestDetailComponent} from './manifest-detail/manifest-detail.component';
import {Manifest} from './manifest';
import {ManifestAlertComponent} from './manifest-alert/manifest-alert.component';

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

    @ViewChild(ManifestAlertComponent)
    alert: ManifestAlertComponent;

    ngOnInit(): void {
    }

    openDetail(item: Manifest) {
        this.detail.open(item);
    }

    openAlert() {
        this.alert.open();
    }
}
