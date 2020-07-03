import {Component, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {StorageProvisionerListComponent} from './storage-provisioner-list/storage-provisioner-list.component';
import {StorageProvisionerCreateComponent} from './storage-provisioner-create/storage-provisioner-create.component';
import {StorageProvisionerCreateNfsComponent} from './storage-provisioner-create/storage-provisioner-create-nfs/storage-provisioner-create-nfs.component';
import {CreateStorageProvisionerRequest, StorageProvisioner} from './storage-provisioner';
import {StorageProvisionerDeleteComponent} from "./storage-provisioner-delete/storage-provisioner-delete.component";

@Component({
    selector: 'app-storage-provisioner',
    templateUrl: './storage-provisioner.component.html',
    styleUrls: ['./storage-provisioner.component.css']
})
export class StorageProvisionerComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    @ViewChild(StorageProvisionerListComponent, {static: true})
    list: StorageProvisionerListComponent;

    @ViewChild(StorageProvisionerCreateComponent, {static: true})
    create: StorageProvisionerCreateComponent;

    @ViewChild(StorageProvisionerDeleteComponent, {static: true})
    delete: StorageProvisionerDeleteComponent;

    @ViewChild(StorageProvisionerCreateNfsComponent, {static: true})
    nfs: StorageProvisionerCreateNfsComponent;
    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items: StorageProvisioner[]) {
        this.delete.open(items);
    }


    openSelected(item: CreateStorageProvisionerRequest) {
        console.log(item.name);
        switch (item.type) {
            case 'nfs':
                this.nfs.open(item);
                break;
        }
    }

    refresh() {
        this.list.refresh();
    }

}
