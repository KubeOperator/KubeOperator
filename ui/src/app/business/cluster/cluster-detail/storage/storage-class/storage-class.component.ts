import {Component, OnInit, ViewChild} from '@angular/core';
import {StorageClassCreateComponent} from './storage-class-create/storage-class-create.component';
import {StorageClassListComponent} from './storage-class-list/storage-class-list.component';
import {CreateStorageClassRequest} from '../storage';
import {StorageClassCreateNfsComponent} from './storage-class-create/storage-class-create-nfs/storage-class-create-nfs.component';

@Component({
    selector: 'app-storage-class',
    templateUrl: './storage-class.component.html',
    styleUrls: ['./storage-class.component.css']
})
export class StorageClassComponent implements OnInit {

    constructor() {
    }

    @ViewChild(StorageClassCreateComponent, {static: true})
    create: StorageClassCreateComponent;

    @ViewChild(StorageClassCreateNfsComponent, {static: true})
    createNfs: StorageClassCreateNfsComponent;

    @ViewChild(StorageClassListComponent, {static: true})
    list: StorageClassListComponent;

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    refresh() {
        this.list.refresh();
    }


    openSelected(item: CreateStorageClassRequest) {
        if (item.provisioner === 'nfs') {
            this.createNfs.open();
        }
    }

}
