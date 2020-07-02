import {Component, OnInit, ViewChild} from '@angular/core';
import {StorageClassCreateComponent} from './storage-class-create/storage-class-create.component';
import {StorageClassListComponent} from './storage-class-list/storage-class-list.component';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';

@Component({
    selector: 'app-storage-class',
    templateUrl: './storage-class.component.html',
    styleUrls: ['./storage-class.component.css']
})
export class StorageClassComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    @ViewChild(StorageClassCreateComponent, {static: true})
    create: StorageClassCreateComponent;

    @ViewChild(StorageClassListComponent, {static: true})
    list: StorageClassListComponent;

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openCreate() {
        this.create.open();
    }

    refresh() {
        this.list.refresh();
    }
}
