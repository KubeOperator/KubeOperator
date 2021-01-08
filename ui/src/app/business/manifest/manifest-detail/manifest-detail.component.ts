import {Component, OnInit} from '@angular/core';
import {Item, Manifest} from '../manifest';

@Component({
    selector: 'app-manifest-detail',
    templateUrl: './manifest-detail.component.html',
    styleUrls: ['./manifest-detail.component.css']
})
export class ManifestDetailComponent implements OnInit {

    constructor() {
    }

    opened = false;
    item: Manifest = new Manifest();

    ngOnInit(): void {
    }

    open(item: Manifest) {
        this.item = item;
        this.opened = true;
    }

}
